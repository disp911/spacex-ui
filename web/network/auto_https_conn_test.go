package network

import (
	"bufio"
	"net"
	"net/http"
	"testing"
)

func TestPlainHTTPGetsAProperRedirect(t *testing.T) {
	server, client := net.Pipe()
	go func() {
		NewAutoHttpsConn(server).Read(make([]byte, 64))
	}()
	go client.Write([]byte("GET /docs/?a=1 HTTP/1.1\r\nHost: example.com:8443\r\n\r\n"))
	resp, err := http.ReadResponse(bufio.NewReader(client), nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Proto != "HTTP/1.1" || resp.StatusCode != http.StatusTemporaryRedirect || !resp.Close {
		t.Fatalf("%s %d close=%v", resp.Proto, resp.StatusCode, resp.Close)
	}
	if loc := resp.Header.Get("Location"); loc != "https://example.com:8443/docs/?a=1" {
		t.Fatalf("Location %q", loc)
	}
}
