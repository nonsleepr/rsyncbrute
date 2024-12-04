package rsync

import (
	"encoding/base64"
	"fmt"
	"net"
	"regexp"

	"golang.org/x/crypto/md4"
)

const (
	RsyncVersion = 29
)

func Connect(host string, port int) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write([]byte(fmt.Sprintf("@RSYNCD: %d\n", RsyncVersion)))
	return conn, err
}

func Auth(conn net.Conn, username, password, share string) bool {
	// Send share name
	_, err := conn.Write([]byte(share + "\n"))
	if err != nil {
		return false
	}

	// Read challenge
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}

	challengeStr := string(buffer[:n])
	re := regexp.MustCompile(`@RSYNCD: AUTHREQD (.+)\n`)
	matches := re.FindStringSubmatch(challengeStr)
	if len(matches) < 2 {
		return false
	}

	challenge := matches[1]

	// Calculate response
	hash := md4.New()
	hash.Write([]byte{0, 0, 0, 0})
	hash.Write([]byte(password))
	hash.Write([]byte(challenge))
	digest := hash.Sum(nil)

	resp := base64.StdEncoding.EncodeToString(digest)
	resp = resp[:len(resp)-2] // Remove padding

	// Send authentication
	authStr := fmt.Sprintf("%s %s\n", username, resp)
	_, err = conn.Write([]byte(authStr))
	if err != nil {
		return false
	}

	// Check response
	n, err = conn.Read(buffer)
	if err != nil {
		return false
	}

	return string(buffer[:n]) == "@RSYNCD: OK\n"
}
