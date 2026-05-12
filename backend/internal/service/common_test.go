package service

import "testing"

func TestEncodeParamsHeaderUsesXTokenOnly(t *testing.T) {
	headers := encodeParamsHeader(map[string]string{"code": "hello world"})
	if headers[paramHeaderName] == "" {
		t.Fatalf("%s header is empty", paramHeaderName)
	}
	if _, ok := headers["X-CD-Params"]; ok {
		t.Fatalf("legacy X-CD-Params header should not be emitted")
	}
}
