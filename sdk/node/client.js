import fetch from "node-fetch";

export class AiguardrailsClient {
  constructor(baseUrl, appId, secret) {
    this.baseUrl = baseUrl;
    this.appId = appId;
    this.secret = secret;
  }

  async promptCheck(prompt) {
    return this._post("/v1/guardrails/prompt-check", { prompt });
  }

  async outputFilter(output) {
    return this._post("/v1/guardrails/output-filter", { output });
  }

  async agentPlan(prompt, tools = []) {
    return this._post("/v1/agent/plan", { prompt, tools });
  }

  async _post(path, body) {
    const res = await fetch(this.baseUrl + path, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-App-Id": this.appId,
        "X-App-Secret": this.secret,
      },
      body: JSON.stringify(body),
    });
    if (!res.ok) {
      const text = await res.text();
      throw new Error(`Request failed: ${res.status} ${text}`);
    }
    return res.json();
  }
}

// Example usage:
// const client = new AiguardrailsClient("http://localhost:8080", "appId", "secret");
// client.promptCheck("hello").then(console.log);

