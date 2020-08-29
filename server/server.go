package server

import (
	"bufio"
	"fmt"
	"net"
)

type session struct {
	discid  string
	trackCount  int
	offsets []int
	seconds int
	mbzResult
}

func handle(conn net.Conn) {
	defer conn.Close()
	sess := &session{}
	conn.Write([]byte("100 welcome\n"))
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		sess.dispatch(scanner.Text(), conn)
	}
}

func Serve() {
	listen := ":8880"
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		fmt.Printf("listen error %s\n", err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("accept error %s\n", err)
			break
		}
		go handle(conn)
	}
}
