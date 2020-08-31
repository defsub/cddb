// Copyright (C) 2020 The cdmbz Authors.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package server

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"os"
)

type session struct {
	discid     string
	trackCount int
	offsets    []int
	seconds    int
	mbzResult
}

func handle(conn net.Conn) {
	defer conn.Close()
	sess := &session{}
	conn.SetDeadline(time.Now().Add(20 * time.Second))

	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	version := "0.1"
	date := time.Now().Local().String()
	banner := fmt.Sprintf("201 %s CDDBP %s at %s\n", host, version, date)

	conn.Write([]byte(banner))
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
