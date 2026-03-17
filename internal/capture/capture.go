package capture

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type VxConnection struct {
	Timestamp  time.Time
	RemoteAddr string
	Service    string
	Payload    string
}

// fake banners to bait a first-line response before we close
var _banners = map[string]string{
	"SSH":  "SSH-2.0-OpenSSH_8.9p1 Ubuntu-3ubuntu0.6\r\n",
	"FTP":  "220 (vsFTPd 3.0.5)\r\n",
	"HTTP": "",
}

func VxStartListener(service string, port int, ch chan<- VxConnection) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go _handle(conn, service, ch)
	}
}

func _handle(conn net.Conn, service string, ch chan<- VxConnection) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(8 * time.Second))

	if b := _banners[service]; b != "" {
		_, _ = conn.Write([]byte(b))
	}

	line, _ := bufio.NewReader(io.LimitReader(conn, 512)).ReadString('\n')

	ch <- VxConnection{
		Timestamp:  time.Now(),
		RemoteAddr: conn.RemoteAddr().String(),
		Service:    service,
		Payload:    strings.TrimSpace(line),
	}
}