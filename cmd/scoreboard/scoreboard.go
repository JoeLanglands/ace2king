package main

import (
	"fmt"
	"strconv"

	"github.com/JoeLanglands/ace2king/pkg/scoring"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	uploadStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	infoStyle = lipgloss.NewStyle().
			Align(lipgloss.Left, lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder()).
			Foreground(lipgloss.Color("#58c7e0"))
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#b141f1"))
)

type gameState uint8

const (
	inProgress gameState = iota
	complete
)

type ScoreboardModel struct {
	state     gameState         `json:"-"`
	table     table.Model       `json:"-"`
	Players   *[]scoring.Player `json:"players"`
	cursor    int               `json:"-"`
	RoundCard string            `json:"roundCard"`
	Round     int               `json:"round"`
	textInput textinput.Model   `json:"-"`
}

func NewScoreboardModel(players *[]scoring.Player) ScoreboardModel {
	columns := []table.Column{
		{Title: "Card", Width: 5},
	}

	for _, player := range *players {
		columns = append(columns, table.Column{

			Title: player.Name,
			Width: maxInt(len(player.Name), 4),
		})
	}

	var rows []table.Row
	for _, card := range scoring.Cards {
		rows = append(rows, table.Row{card})
	}

	rows = append(rows, table.Row{"Total"})
	rows = append(rows, table.Row{"Wins"})

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#f0f0f0")).
		BorderBottom(true).
		Align(lipgloss.Center)

	t.SetStyles(s)

	ti := textinput.New()
	ti.Placeholder = "Enter a score"
	ti.CharLimit = 4
	ti.Width = 8
	ti.Focus()

	return ScoreboardModel{
		state:     inProgress,
		table:     t,
		Players:   players,
		textInput: ti,
		RoundCard: scoring.Cards[0],
	}
}

func (m ScoreboardModel) Init() tea.Cmd {
	return nil
}

func (m ScoreboardModel) View() string {
	var v string

	uploadView := "Update scores:\n\n"

	for i, player := range *m.Players {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		uploadView += fmt.Sprintf("%s %s - %d\n", cursor, player.Name, player.Scores[m.RoundCard])
	}

	infoView := m.makeInfoView()

	uploadView = lipgloss.JoinVertical(lipgloss.Left, infoView, uploadStyle.Render(uploadView), m.textInput.View())

	v += lipgloss.JoinHorizontal(lipgloss.Left, baseStyle.Render(m.table.View()), "    ", uploadView)

	help := helpStyle.Render("Use the arrow keys to navigate, ctrl+c to quit,\nEnter to add score, n to move to the next round")

	return lipgloss.JoinVertical(lipgloss.Left, v, help)
}

func (m ScoreboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			score, err := strconv.Atoi(m.textInput.Value())
			if err != nil {
				m.textInput.Err = err
				return m, nil
			}
			if m.state == inProgress {
				(*m.Players)[m.cursor].AddScore(score, m.RoundCard)
			}
			m.textInput.Reset()
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(*m.Players)-1 {
				m.cursor++
			}
		case "n":
			if m.Round == len(scoring.Cards)-1 && m.state == inProgress {
				// get stuck in this state when the game is over to view table
				rows := m.refreshTableRows()
				m.table.SetRows(rows)
				m.table.MoveDown(1)
				m.state = complete
				saveGameState(m)
				return m, nil
			}
			if m.state == inProgress {
				rows := m.refreshTableRows()
				m.table.SetRows(rows)
				m.Round++
				m.RoundCard = scoring.Cards[m.Round]
				m.table.MoveDown(1)
				m.cursor = 0
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m ScoreboardModel) refreshTableRows() []table.Row {
	var rows []table.Row

	for _, card := range scoring.Cards {
		row := table.Row{card}
		for _, player := range *m.Players {
			row = append(row, strconv.Itoa(player.Scores[card]))
		}
		rows = append(rows, row)
	}

	totalRow := table.Row{"Total"}
	winRow := table.Row{"Wins"}
	for _, player := range *m.Players {
		totalRow = append(totalRow, strconv.Itoa(player.Score))
		// not handling error right now
		wins, _ := player.CountWins(m.Round)
		winRow = append(winRow, strconv.Itoa(wins))
	}

	rows = append(rows, totalRow, winRow)

	return rows
}

func (m ScoreboardModel) makeInfoView() string {
	var iv string

	iv += fmt.Sprintf("Current card: %s\n", m.RoundCard)
	iv += fmt.Sprintf("Soufl??: %s\n", (*m.Players)[m.Round%len(*m.Players)].Name)
	iv += fmt.Sprintf("Round: %d\n", m.Round)

	return infoStyle.Render(iv)
}
