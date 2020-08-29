package server

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"strings"
)

func (sess *session) dispatch(msg string, w io.Writer) error {
	fmt.Printf(">> %s\n", msg)
	w2 := io.MultiWriter(w, os.Stdout)

	cmd := strings.Split(msg, " ")
	if len(cmd) > 2 {
		if cmd[0] == "cddb" && cmd[1] == "query" {
			// cddb query discid ntrks off1 off2 ... nsecs
			sess.query(cmd[2:], w2)
		} else if cmd[0] == "cddb" && cmd[1] == "read" {
			// cddb read categ discid
			sess.read(cmd[2:], w2)
		} else if cmd[0] == "cddb" && cmd[1] == "hello" {
			// cddb hello username hostname clientname version
			sess.hello(cmd[2:], w2)
		}
	} else if len(cmd) == 2 {
		if cmd[0] == "proto" {
			sess.proto(cmd[1:], w2)
		}
	}
	return nil
}

func (sess *session) query(cmd []string, w io.Writer) {
	// cddb query ee0fbc11 17 150 15918 33885 51210 72633 91825 107498
	// 123166 142810 157443 171760 190726 208521 230252 252341 269891
	// 286394 4030
	sess.discid = cmd[0]
	sess.trackCount = atoi(cmd[1])
	for i := 0; i < sess.trackCount; i++ {
		sess.offsets = append(sess.offsets, atoi(cmd[2+i]))
	}
	sess.seconds = atoi(cmd[len(cmd)-1])

	sess.mbzLookup()
	release := sess.pick()
	if release == nil {
		w.Write([]byte("202 not found\n"))
		return
	}

	wr := bufio.NewWriter(w)
	wr.WriteString(fmt.Sprintf("200 rock %s %s / %s\n", sess.discid,
		release.Artist[0].Name, release.Title))

	wr.Flush()
}

func (sess *session) read(cmd []string, w io.Writer) {
	wr := bufio.NewWriter(w)
	wr.WriteString(fmt.Sprintf("210 %s %s\n", "rock", sess.discid))
	wr.WriteString("# xmcd\n")
	wr.WriteString("#\n")
	wr.WriteString("# Track frame offsets:\n")
	for _, v := range sess.offsets {
		wr.WriteString(fmt.Sprintf("#\t%d\n", v))
	}
	wr.WriteString("#\n")
	wr.WriteString(fmt.Sprintf("# Disc length: %d seconds\n", sess.seconds))
	wr.WriteString("#\n")
	wr.WriteString(fmt.Sprintf("# Revision: 1\n"))
	wr.WriteString(fmt.Sprintf("# Submitted via: MusicBrainz 1.0 tbd\n"))
	wr.WriteString("#\n")

	release := sess.pick()
	if release != nil {
		wr.WriteString(fmt.Sprintf("DISCID=%s\n", sess.discid))
		wr.WriteString(fmt.Sprintf("DTITLE=%s / %s\n", release.Artist[0].Name, release.Title))
		date := strings.Split(release.Date, "-")
		wr.WriteString(fmt.Sprintf("DYEAR=%s\n", date[0]))

		media := release.media(sess.trackCount)
		if media != nil {
			for _, t := range media.Tracks  {
				wr.WriteString(fmt.Sprintf("TTITLE%d=%s\n", t.Position-1, t.Title))
			}
		}
	}
	wr.WriteString(".\n")
	wr.Flush()
}

func (sess *session) hello(cmd []string, w io.Writer) {
	w.Write([]byte("200 hello and welcome.\n"))
}

func (sess *session) proto(cmd []string, w io.Writer) {
	w.Write([]byte("200 ok\n"))
}
