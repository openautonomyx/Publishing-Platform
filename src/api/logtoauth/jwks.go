package logtoauth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

// jwksCache fetches the issuer's JSON Web Key Set and caches the public keys by kid.
// An unknown kid triggers a refetch (key rotation), at most once per minRefresh.
type jwksCache struct {
	url        string
	client     *http.Client
	ttl        time.Duration
	minRefresh time.Duration

	mu        sync.Mutex
	keys      map[string]any
	fetchedAt time.Time
}

func newJWKSCache(url string, client *http.Client) *jwksCache {
	return &jwksCache{url: url, client: client, ttl: time.Hour, minRefresh: 30 * time.Second}
}

func (c *jwksCache) key(ctx context.Context, kid string) (any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	stale := time.Since(c.fetchedAt) > c.ttl
	if k, ok := c.keys[kid]; ok && !stale {
		return k, nil
	}
	if stale || time.Since(c.fetchedAt) > c.minRefresh {
		if err := c.refresh(ctx); err != nil {
			// Keep serving cached keys if the issuer is briefly unreachable.
			if k, ok := c.keys[kid]; ok {
				return k, nil
			}
			return nil, err
		}
	}
	if k, ok := c.keys[kid]; ok {
		return k, nil
	}
	return nil, fmt.Errorf("no signing key with kid %q", kid)
}

func (c *jwksCache) refresh(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch jwks: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch jwks: status %d", resp.StatusCode)
	}

	var set struct {
		Keys []jwk `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return fmt.Errorf("decode jwks: %w", err)
	}

	keys := make(map[string]any, len(set.Keys))
	for _, k := range set.Keys {
		if k.Use != "" && k.Use != "sig" {
			continue
		}
		pub, err := k.publicKey()
		if err != nil {
			continue // skip key types we do not verify with
		}
		keys[k.Kid] = pub
	}
	if len(keys) == 0 {
		return errors.New("jwks contains no usable signing keys")
	}
	c.keys = keys
	c.fetchedAt = time.Now()
	return nil
}

type jwk struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func (k jwk) publicKey() (any, error) {
	switch k.Kty {
	case "EC":
		var curve elliptic.Curve
		switch k.Crv {
		case "P-256":
			curve = elliptic.P256()
		case "P-384":
			curve = elliptic.P384()
		case "P-521":
			curve = elliptic.P521()
		default:
			return nil, fmt.Errorf("unsupported curve %q", k.Crv)
		}
		x, err := b64Int(k.X)
		if err != nil {
			return nil, err
		}
		y, err := b64Int(k.Y)
		if err != nil {
			return nil, err
		}
		if !curve.IsOnCurve(x, y) {
			return nil, errors.New("ec point is not on curve")
		}
		return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
	case "RSA":
		n, err := b64Int(k.N)
		if err != nil {
			return nil, err
		}
		e, err := b64Int(k.E)
		if err != nil {
			return nil, err
		}
		if !e.IsInt64() || e.Int64() < 2 || e.Int64() > 1<<31-1 {
			return nil, errors.New("invalid rsa exponent")
		}
		return &rsa.PublicKey{N: n, E: int(e.Int64())}, nil
	default:
		return nil, fmt.Errorf("unsupported key type %q", k.Kty)
	}
}

func b64Int(s string) (*big.Int, error) {
	if s == "" {
		return nil, errors.New("missing key parameter")
	}
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(b), nil
}
