package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func requireE2E(t *testing.T) string {
	t.Helper()

	if os.Getenv("RUN_E2E") != "1" {
		t.Skip("set RUN_E2E=1 to run e2e tests against a live server")
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return baseURL
}

func requestJSON(t *testing.T, client *http.Client, method, url, token string, body any) *http.Response {
	t.Helper()

	var payload *bytes.Reader
	if body == nil {
		payload = bytes.NewReader(nil)
	} else {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		payload = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request %s %s: %v", method, url, err)
	}

	return resp
}

func decodeBody[T any](t *testing.T, resp *http.Response) T {
	t.Helper()

	defer resp.Body.Close()

	var value T
	if err := json.NewDecoder(resp.Body).Decode(&value); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	return value
}

func mustLogin(t *testing.T, client *http.Client, baseURL, role string) string {
	t.Helper()

	resp := requestJSON(t, client, http.MethodPost, baseURL+"/dummyLogin", "", map[string]any{
		"role": role,
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("dummyLogin expected 200, got %d", resp.StatusCode)
	}

	var payload struct {
		Token string `json:"token"`
	}
	payload = decodeBody[struct {
		Token string `json:"token"`
	}](t, resp)

	if payload.Token == "" {
		t.Fatalf("dummyLogin returned empty token")
	}

	return payload.Token
}

func tomorrowDateUTC() string {
	return time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")
}

func uniqueRoomName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
