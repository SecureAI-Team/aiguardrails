import json
import requests


class AiguardrailsClient:
    def __init__(self, base_url: str, app_id: str, secret: str, timeout: float = 5.0):
        self.base_url = base_url.rstrip("/")
        self.app_id = app_id
        self.secret = secret
        self.timeout = timeout

    def prompt_check(self, prompt: str):
        return self._post("/v1/guardrails/prompt-check", {"prompt": prompt})

    def output_filter(self, output: str):
        return self._post("/v1/guardrails/output-filter", {"output": output})

    def agent_plan(self, prompt: str, tools=None):
        if tools is None:
            tools = []
        return self._post("/v1/agent/plan", {"prompt": prompt, "tools": tools})

    def _post(self, path: str, body: dict):
        url = f"{self.base_url}{path}"
        headers = {
            "Content-Type": "application/json",
            "X-App-Id": self.app_id,
            "X-App-Secret": self.secret,
        }
        resp = requests.post(url, headers=headers, data=json.dumps(body), timeout=self.timeout)
        if not resp.ok:
            raise RuntimeError(f"Request failed {resp.status_code}: {resp.text}")
        return resp.json()


# Example:
# client = AiguardrailsClient("http://localhost:8080", "appId", "secret")
# print(client.prompt_check("hello"))

