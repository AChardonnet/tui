package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	boxer "github.com/treilik/bubbleboxer"
)

const (
	mainAddr     = "main"
	logsAddr     = "logs"
	playlistAddr = "playlist"
)

func stripErr(n boxer.Node, _ error) boxer.Node {
	return n
}

func tui(logChannel chan logMsg) {
	vp := viewport.New(0, 0)
	main := &mainModel{
		viewport: vp,
		cacheBar: *newCacheBar(),
	}
	playlist := newPlaylistModel()
	logs := newLogModel()

	m := model{
		tui:     boxer.Boxer{},
		logChan: logChannel,
	}

	topPane := boxer.Node{

		SizeFunc: func(node boxer.Node, widthOrHeight int) []int {
			mainWidth := widthOrHeight * 72 / 100
			playlistWidth := widthOrHeight - mainWidth
			if playlistWidth <= 0 {
				playlistWidth = 1
			}
			if mainWidth <= 0 {
				mainWidth = 1
			}
			return []int{mainWidth, playlistWidth}
		},
		Children: []boxer.Node{
			stripErr(m.tui.CreateLeaf(mainAddr, main)),
			stripErr(m.tui.CreateLeaf(playlistAddr, playlist)),
		},
	}

	m.tui.LayoutTree = boxer.Node{
		VerticalStacked: true,
		SizeFunc: func(node boxer.Node, widthOrHeight int) []int {
			topHeight := widthOrHeight * 75 / 100
			logHeight := widthOrHeight - topHeight
			if logHeight <= 0 {
				logHeight = 1
			}
			if topHeight <= 0 {
				topHeight = 1
			}
			return []int{topHeight, logHeight}
		},
		Children: []boxer.Node{
			topPane,
			stripErr(m.tui.CreateLeaf(logsAddr, logs)),
		},
	}
	p := tea.NewProgram(m)
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Println(err)
	}
	p.ExitAltScreen()
}
