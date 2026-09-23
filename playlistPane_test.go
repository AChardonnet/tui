package main

import "testing"

func TestPlaylistModelIncludesHeaderAndItems(t *testing.T) {
	m := newPlaylistModel()
	m.items = []playlistItem{{Title: "Intro"}, {Title: "Main topic"}, {Title: "Outro"}}
	m.selected = 1

	view := m.View()
	if len(view) == 0 {
		t.Fatal("playlist view should not be empty")
	}
	if !contains(view, "Playlist") {
		t.Fatal("playlist view should render a header")
	}
	if !contains(view, "Main topic") {
		t.Fatal("playlist view should render current items")
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})())
}
