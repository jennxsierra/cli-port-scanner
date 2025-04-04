package banner

import (
	"bytes"
	"net"
	"testing"
	"time"
)

// fakeConn is a simple implementation of net.Conn for testing purposes.
type fakeConn struct {
	input  *bytes.Buffer
	output *bytes.Buffer
}

func (f *fakeConn) Read(b []byte) (int, error) {
	return f.input.Read(b)
}
func (f *fakeConn) Write(b []byte) (int, error) {
	return f.output.Write(b)
}
func (f *fakeConn) Close() error {
	return nil
}
func (f *fakeConn) LocalAddr() net.Addr                { return nil }
func (f *fakeConn) RemoteAddr() net.Addr               { return nil }
func (f *fakeConn) SetDeadline(t time.Time) error      { return nil }
func (f *fakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (f *fakeConn) SetWriteDeadline(t time.Time) error { return nil }

func TestGrabBannerNonHTTP(t *testing.T) {
	// Prepare a fake connection that returns a fixed banner.
	bannerData := "FTP Service Ready"
	input := bytes.NewBufferString(bannerData)
	output := &bytes.Buffer{}
	conn := &fakeConn{input: input, output: output}

	result := GrabBanner(conn, 21)
	if result != bannerData {
		t.Errorf("Expected banner %q, got %q", bannerData, result)
	}
}

func TestGrabBannerHTTP(t *testing.T) {
	// Prepare fake HTTP response data.
	httpResponse := "HTTP/1.0 200 OK\r\nServer: TestServer\r\n\r\n"
	input := bytes.NewBufferString(httpResponse)
	output := &bytes.Buffer{}
	conn := &fakeConn{input: input, output: output}

	result := GrabBanner(conn, 80)
	// Expect to extract only the TestServer portion.
	if result != "TestServer" {
		t.Errorf("Expected banner ‘TestServer’, got %q", result)
	}
}
