// Filename: banner.go
// Description: Contains the GrabBanner function to read banner information from a connection.
// This function is used in the port scanning process to retrieve service banners.
package banner

import (
	"bufio"
	"net"
	"strings"
	"time"
)

// GrabBanner reads banner information from a connection. For HTTP (port 80),
// it sends a HEAD request and extracts the "Server:" header if present.
func GrabBanner(conn net.Conn, port int) string {
	// For HTTP port, send a HEAD request to get HTTP response headers.
	if port == 80 {
		_, err := conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
		if err != nil {
			return ""
		}
	}

	// Set a deadline for reading.
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	var banner []byte
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			banner = append(banner, buf[:n]...)
			// Extend the deadline a bit to capture any extra data.
			conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		}
		if err != nil {
			break
		}
	}

	result := strings.TrimSpace(string(banner))

	// For HTTP, extract only the "Server:" header.
	if port == 80 {
		scanner := bufio.NewScanner(strings.NewReader(result))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(strings.ToLower(line), "server:") {
				return strings.TrimSpace(line[7:]) // Remove the "Server:" prefix.
			}
		}
		// Return the entire result if no Server header found.
		return result
	}

	return result
}
