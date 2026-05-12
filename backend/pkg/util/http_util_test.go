package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostRequestWithOptionsJSONAndCookie(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Cookie"); got != "ASPSESSIONIDXYZ=abc" {
			t.Fatalf("cookie = %q", got)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["timezone"] != "seed" || body["sign"] != "payload" {
			t.Fatalf("unexpected body: %#v", body)
		}
		http.SetCookie(w, &http.Cookie{Name: "ASPSESSIONIDXYZ", Value: "next"})
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	res, err := PostRequestWithOptions(server.URL, "seed", "payload", ShellPostOptions{
		Session:   "ASPSESSIONIDXYZ=abc",
		ShellType: "asp",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Cookie != "ASPSESSIONIDXYZ=next" {
		t.Fatalf("cookie capture = %q", res.Cookie)
	}
	if res.Body != "ok" {
		t.Fatalf("body = %q", res.Body)
	}
}

func TestXOREncryptDecryptOffsets(t *testing.T) {
	seed := "1710000000"
	plain := []byte("CloudDrop")
	encoded := EncryptWithOffset(plain, seed, 5)
	decoded, err := DecryptWithOffset(encoded, seed, 5)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != string(plain) {
		t.Fatalf("decoded = %q", decoded)
	}
}
