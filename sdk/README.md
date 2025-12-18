# SDK Usage (Go / Node / Python)

## Go (module under backend)
```go
import "aiguardrails/pkg/sdk"

client := sdk.NewClient("http://localhost:8080", "appId", "secret")
res, err := client.PromptCheck("hello")
```

## Node
```js
import { AiguardrailsClient } from "./node/client.js";
const c = new AiguardrailsClient("http://localhost:8080", "appId", "secret");
const res = await c.promptCheck("hello");
```

## Python
```py
from client import AiguardrailsClient
c = AiguardrailsClient("http://localhost:8080", "appId", "secret")
print(c.prompt_check("hello"))
```

Notes:
- Set `X-App-Id` / `X-App-Secret`.
- Handle 429/403 for rate/quota and guardrail blocks; implement retries with backoff.
- Configure timeouts per environment.

