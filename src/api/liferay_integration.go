package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

// LiferayClient provides integration with Liferay DXP
type LiferayClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewLiferayClient creates a new Liferay client
func NewLiferayClient() *LiferayClient {
	return &LiferayClient{
		baseURL: os.Getenv("LIFERAY_URL"),
		apiKey:  os.Getenv("LIFERAY_API_KEY"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// LiferayContent represents content in Liferay
type LiferayContent struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	Body         string                 `json:"body"`
	ContentType  string                 `json:"contentType"`
	Status       string                 `json:"status"`
	CreatorID    string                 `json:"creatorId"`
	CreatedDate  time.Time              `json:"createdDate"`
	ModifiedDate time.Time              `json:"modifiedDate"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// GetContent retrieves content from Liferay
func (lc *LiferayClient) GetContent(contentID string) (*LiferayContent, error) {
	logger.Info("fetching content from Liferay", zap.String("content_id", contentID))

	url := fmt.Sprintf("%s/o/content-service/v1.0/content/%s", lc.baseURL, contentID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.apiKey))
	req.Header.Set("Accept", "application/json")

	resp, err := lc.client.Do(req)
	if err != nil {
		logger.Error("failed to fetch content from Liferay", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("liferay returned %d: %s", resp.StatusCode, string(body))
	}

	var content LiferayContent
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		logger.Error("failed to decode Liferay response", zap.Error(err))
		return nil, err
	}

	return &content, nil
}

// CreateContent creates new content in Liferay
func (lc *LiferayClient) CreateContent(title, body, contentType string) (string, error) {
	logger.Info("creating content in Liferay", zap.String("title", title))

	url := fmt.Sprintf("%s/o/content-service/v1.0/content", lc.baseURL)

	payload := map[string]interface{}{
		"title":       title,
		"body":        body,
		"contentType": contentType,
		"status":      "DRAFT",
		"metadata": map[string]interface{}{
			"createdBy": "creative-platform",
		},
	}

	bodyBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := lc.client.Do(req)
	if err != nil {
		logger.Error("failed to create content in Liferay", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("liferay returned %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logger.Error("failed to decode Liferay response", zap.Error(err))
		return "", err
	}

	contentID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("no content ID in response")
	}

	return contentID, nil
}

// UpdateContent updates content in Liferay
func (lc *LiferayClient) UpdateContent(contentID, title, body string) error {
	logger.Info("updating content in Liferay", zap.String("content_id", contentID))

	url := fmt.Sprintf("%s/o/content-service/v1.0/content/%s", lc.baseURL, contentID)

	payload := map[string]interface{}{
		"title": title,
		"body":  body,
	}

	bodyBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := lc.client.Do(req)
	if err != nil {
		logger.Error("failed to update content in Liferay", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("liferay returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// PublishContent publishes content in Liferay (changes status from DRAFT to PUBLISHED)
func (lc *LiferayClient) PublishContent(contentID string) error {
	logger.Info("publishing content in Liferay", zap.String("content_id", contentID))

	url := fmt.Sprintf("%s/o/content-service/v1.0/content/%s/publish", lc.baseURL, contentID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := lc.client.Do(req)
	if err != nil {
		logger.Error("failed to publish content in Liferay", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("liferay returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetContentStatus gets approval status of content
func (lc *LiferayClient) GetContentStatus(contentID string) (string, error) {
	content, err := lc.GetContent(contentID)
	if err != nil {
		return "", err
	}

	return content.Status, nil
}

// SearchContent searches for content in Liferay
func (lc *LiferayClient) SearchContent(query, contentType string) ([]*LiferayContent, error) {
	logger.Info("searching content in Liferay", zap.String("query", query))

	url := fmt.Sprintf(
		"%s/o/content-service/v1.0/content?search=%s&contentType=%s",
		lc.baseURL, query, contentType,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.apiKey))
	req.Header.Set("Accept", "application/json")

	resp, err := lc.client.Do(req)
	if err != nil {
		logger.Error("failed to search content in Liferay", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("liferay returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Items []*LiferayContent `json:"items"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logger.Error("failed to decode Liferay response", zap.Error(err))
		return nil, err
	}

	return result.Items, nil
}

// UpdateContentStatus updates the status of content in Liferay
func (lc *LiferayClient) UpdateContentStatus(contentID, status string) error {
	logger.Info("updating content status in Liferay",
		zap.String("content_id", contentID),
		zap.String("status", status),
	)

	url := fmt.Sprintf("%s/o/content-service/v1.0/content/%s", lc.baseURL, contentID)

	payload := map[string]interface{}{
		"status": status,
	}

	bodyBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := lc.client.Do(req)
	if err != nil {
		logger.Error("failed to update content status in Liferay", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("liferay returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// BatchPublish publishes multiple content items
func (lc *LiferayClient) BatchPublish(contentIDs []string) error {
	logger.Info("batch publishing content", zap.Int("count", len(contentIDs)))

	for _, id := range contentIDs {
		err := lc.PublishContent(id)
		if err != nil {
			logger.Warn("failed to publish content", zap.String("id", id), zap.Error(err))
			// Continue with other items
		}
	}

	return nil
}

// GetContentVersions retrieves all versions of content
func (lc *LiferayClient) GetContentVersions(contentID string) ([]*LiferayContent, error) {
	logger.Info("fetching content versions", zap.String("content_id", contentID))

	url := fmt.Sprintf("%s/o/content-service/v1.0/content/%s/versions", lc.baseURL, contentID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.apiKey))
	req.Header.Set("Accept", "application/json")

	resp, err := lc.client.Do(req)
	if err != nil {
		logger.Error("failed to fetch content versions", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("liferay returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Items []*LiferayContent `json:"items"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logger.Error("failed to decode versions response", zap.Error(err))
		return nil, err
	}

	return result.Items, nil
}

// DeleteContent deletes content from Liferay
func (lc *LiferayClient) DeleteContent(contentID string) error {
	logger.Info("deleting content from Liferay", zap.String("content_id", contentID))

	url := fmt.Sprintf("%s/o/content-service/v1.0/content/%s", lc.baseURL, contentID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.apiKey))

	resp, err := lc.client.Do(req)
	if err != nil {
		logger.Error("failed to delete content", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("liferay returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
