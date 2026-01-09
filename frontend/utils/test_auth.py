import requests
import sys

# Usage: python test_auth.py <url> <app_id> <app_secret>

def test_auth(url, app_id, app_secret):
    headers = {
        "Content-Type": "application/json",
        "X-App-Id": app_id,
        "X-App-Secret": app_secret
    }
    data = {
        "prompt": "test"
    }
    try:
        print(f"Sending request to {url}...")
        print(f"Headers: {headers}")
        resp = requests.post(url, json=data, headers=headers)
        print(f"Status: {resp.status_code}")
        print(f"Body: {resp.text}")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print("Usage: python test_auth.py <url> <app_id> <app_secret>")
        sys.exit(1)
    
    url = sys.argv[1]
    app_id = sys.argv[2]
    app_secret = sys.argv[3]
    test_auth(url, app_id, app_secret)
