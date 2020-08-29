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

func (sess *session) pick() *mbzRelease {
	if len(sess.mbzResult.Releases) == 0 {
		return nil
	}
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
