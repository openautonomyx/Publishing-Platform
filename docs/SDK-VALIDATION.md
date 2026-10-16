# SDK build and validation

The root npm package is the published SDK. Its TypeScript compilation and ESLint
scope follow the public `src/index.ts` entrypoint and the modules it exports.
Standalone React pages, the agent prototype, local LLM prototype and deployment
utilities are not part of that library. They remain incomplete and need separate
application build targets; a passing SDK build does not validate them.

The build cleans `dist` first and refuses to emit on compiler errors. This avoids
shipping stale prototype output. The missing `./client` export is removed, the
existing factory/configuration types are exported, and validators import their
actual dependencies directly instead of creating an index-module cycle.

Run:

```sh
npm ci
npm run lint
npm test
npm pack --ignore-scripts
```

`npm test` first builds the library and uses Node's built-in test runner for the
SDK and evaluation API. It also discovers PR #13's language tests once merged.
Jest previously found no tests and collided with the Electron package name;
this command executes real tests without a success-on-empty override.

Verified on Node 24.19.0 / npm 10.9.4 on 2026-10-16:

- Clean `npm ci --ignore-scripts --no-audit --no-fund` succeeded.
- TypeScript SDK build and ESLint passed.
- Four tests passed: SDK public entrypoint/validation, HTTP request/auth/error
  contract, and the two existing workspace HTTP tests.
- The generated tarball contains 51 files (compiled SDK, package metadata,
  README and license). Installing it in a separate local consumer and importing
  both `OpenAutonomyX` and `createOpenAutonomyX` succeeded.
- `npm audit --omit=dev --json` reported zero known production dependency
  vulnerabilities in that snapshot. This does not cover development dependencies,
  container images, unpublished prototypes, or future advisories.
- Existing PR #13's two language-registry tests passed in an isolated checkout.

## Remaining integration and container limitations

The SDK's legacy `/api/v1/publications` methods are tested against a local mock
HTTP server. The Go API exposes `/api/v1/content`, with a different data model.
The workspace API exposes `/api/publications` and different methods. These are
not interchangeable. This change establishes a loadable/testable package; it
does not claim compatibility with the Go API or real channel delivery. A versioned
API/client contract must be designed and tested before connecting those layers.

The root Dockerfile is still not a validated application image: its optional
COPY fallbacks, ignored documentation inputs, and SDK-as-server startup assumptions
remain recorded in `IMPLEMENTATION-AUDIT.md`. Its build stage now copies the SDK
build helper. No container was built or deployed, and nothing was published to npm.

The CI workflow repairs and broader repository audit are in PR #14; write-integrity
regressions are in PR #15. Hosted checks remain blocked by Issue #6.
