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
	"fmt"
	"strings"
)

type mbzResult struct {
	Releases     []mbzRelease `json:"releases"`
	ReleaseCount int          `json:"release-count"`
}

type mbzRelease struct {
	Title   string            `json:"title"`
	Status  string            `json:"status"`
	Country string            `json:"country"`
	Artist  []mbzArtistCredit `json:"artist-credit"`
	Media   []mbzMedia        `json:"media"`
	Date    string            `json:"date"`
}

type mbzArtistCredit struct {
	Name string `json:"name"`
}

type mbzMedia struct {
	Format     string     `json:"format"`
	Title      string     `json:"title"`
	TrackCount int        `json:"track-count"`
	Tracks     []mbzTrack `json:"tracks"`
	Position   int        `json:"position"`
}

type mbzTrack struct {
	Title    string `json:"title"`
	Position int    `json:"position"`
}

func (r *mbzRelease) media(trackCount int) *mbzMedia {
	for _, m := range r.Media {
		if m.TrackCount == trackCount {
			return &m
		}
	}
	return nil
}

func (r *mbzRelease) mediaCount(trackCount int) int {
	return len(r.Media)
}

func (sess *session) pick() *mbzRelease {
	if len(sess.mbzResult.Releases) == 0 {
		return nil
	}
	// TODO for now prefer official US releases and fallback to the first
	// one.
	for _, r := range sess.mbzResult.Releases {
		if r.Status == "Official" && r.Country == "US" {
			return &r
		}
	}
	return &sess.mbzResult.Releases[0]
}

func (sess *session) mbzLookup() error {
	var frames []string
	for _, v := range sess.offsets {
		frames = append(frames, itoa(v))
	}
	inc := "recordings+artist-credits"
	url := fmt.Sprintf("https://musicbrainz.org/ws/2/discid/-?toc=1+%d+%d+%s&fmt=json&inc=%s",
		sess.trackCount, sess.seconds*75, strings.Join(frames, "+"), inc)
	err := getJson(url, &sess.mbzResult)
	return err
}
