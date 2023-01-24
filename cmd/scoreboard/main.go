package main

import (
	"flag"

	"github.com/JoeLanglands/ace2king/pkg/scoring"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func parsePlayers(names []string) []scoring.Player {
	var players []scoring.Player
	caser := cases.Title(language.English)
	for _, n := range names {
		name := caser.String(n)
		players = append(players, scoring.NewPlayer(name))
	}
	return players
}

func main() {
	flag.Parse()
	names := flag.Args()

	players := parsePlayers(names)

	f, err := tea.LogToFile("debug.log", "debug > ")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := tea.NewProgram(NewScoreboardModel(&players), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
