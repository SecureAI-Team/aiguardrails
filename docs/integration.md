# Integration Guide (SDK + Proxy)

## SDK Mode
- Use appId/appSecret headers: `X-App-Id`, `X-App-Secret`.
- Endpoints:
  - `POST /v1/guardrails/prompt-check`
  - `POST /v1/guardrails/output-filter`
  - `POST /v1/agent/plan`
  - `GET /v1/mcp/capabilities`
- SDKs: Go (pkg/sdk), extend similarly for Node/Python; include retries, timeouts, and error handling for 429/403.

## Proxy Mode (Sidecar/Edge)
- Deploy a reverse-proxy that:
  - Adds `X-App-Id`/`X-App-Secret`
  - Routes LLM traffic through guardrail endpoints before upstream
  - Captures logs/audit
- Example (docker-compose sidecar):
  - Run `api` service from compose; sidecar proxies to upstream LLM with pre/post hooks calling `prompt-check` and `output-filter`.
- Config knobs:
  - Upstream base URL, allowed headers, timeout/retry
  - Block/mark mode for output, rate limit behavior on 429
  - Optional: inject tenant/app IDs from platform config

## Admin/OIDC
- Platform admin: `X-Admin-Token` for tenant/app/policy/capability/rule mgmt.
- OIDC (issuer/audience/JWKS) with role mapping (admin/user) for tenant-level operations.

