package game

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
)

/*
- API to get the words : https://www.wordgamedb.com/api/v1/words/random
- Interface -> CheckWord, UpdateUI
*/

/*

 Hangman
 Word: c r _ _ d
  _______
 |       0
 |      /|\
 |       |
 |      /
 |
---

Keyboard:
| c | r | a | i | d |   |   |   |
|   |   |   |   |   |   |   |   |
|   |   |   |   |   |   |   |   |

*/
// creed : : maxGuesses : 6, guesses : 4
// keyboard : {c, r, d, a} --> coloured
// c : 0
// r : 1
// e : [2, 3]
// d : 4

// identifiy the keyboard location for exmaple(a: 0, selected: true) as selected or not.
type selected int

type Game struct {
	word         WordDescription
	won          bool
	session      ssh.Session
	guesses      int
	currentGuess []string
	maxGuesses   int
	keyboard     map[string]selected
	model        model
}

type model struct {
	Tabs       []string
	TabContent []string // change this to include a view
	ActiveTab  int      // index of the currect active tab
	TextInput  textinput.Model
}

func Play(width int, height int, session ssh.Session) Game {
	// create a new game for that player
	var game Game

	// get a random word
	word, err := GetRandomWord()
	if err != nil {
		log.Errorf("Couldn't fetch a word : %s", err)
	}

	game.word = word
	game.currentGuess = make([]string, word.Length)
	game.session = session
	game.model.Tabs = []string{"Game", "Results"}
	game.model.TabContent = []string{"Game Content go here", "Results go here!"}

	// text input
	ti := textinput.New()
	ti.Placeholder = "Char"
	ti.Focus()
	ti.CharLimit = 1
	ti.Width = 10

	game.model.TextInput = ti

	return game
}

// bubbletea
func (g Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "right", "tab":
			g.model.ActiveTab = min(g.model.ActiveTab+1, len(g.model.Tabs)-1)
			return g, nil
		case "left", "shift+tab":
			g.model.ActiveTab = max(g.model.ActiveTab-1, 0)
			return g, nil
		case "ctrl+c", "esc":
			return g, tea.Quit
		case "enter":
			// update the word
			char := g.model.TextInput.Value()
			g.model.TextInput.Reset()
			g.tryChar(char)
			return g, nil
		}
	}
	g.model.TextInput, cmd = g.model.TextInput.Update(msg)

	return g, cmd
}

func (g Game) View() string {
	doc := strings.Builder{}
	m := g.model

	pty, _, _ := g.session.Pty()
	windowWidth, _ := pty.Window.Width, pty.Window.Height

	var tabs []string
	for i, t := range m.Tabs {

		// activeTab := tab.Copy().Border(activeTabBorder, true)
		// check with one is active now
		var tabText string
		_, _, isActive := i == 0, i == len(m.Tabs)-1, i == m.ActiveTab
		if isActive {
			tabText = lipgloss.JoinHorizontal(
				lipgloss.Top,
				activeTab.Render(t),
			)
		} else {
			tabText = lipgloss.JoinHorizontal(
				lipgloss.Top,
				tab.Render(t),
			)
		}

		tabs = append(tabs, tabText)
	}

	rowTab := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabs...,
	)
	gap := tabGap.Render(strings.Repeat(" ", max(0, windowWidth-lipgloss.Width(tabs[0])-2)))
	row := lipgloss.JoinHorizontal(lipgloss.Bottom, rowTab, gap)
	doc.WriteString(row + "\n\n")

	// display the word
	// pass the updated word as string to the render function
	word := lipgloss.NewStyle().Width(10).Align(lipgloss.Center).Render(g.showWord())
	ui := lipgloss.JoinVertical(lipgloss.Center, word)

	wordBox := lipgloss.Place(windowWidth, 1,
		lipgloss.Left, lipgloss.Top,
		wordBoxStyle.Render(ui),
		lipgloss.WithWhitespaceForeground(subtle),
	)

	// doc.WriteString(lipgloss.JoinHorizontal(
	// 	lipgloss.Top,
	// 	hangmanShape.Width(windowWidth/3).Align(lipgloss.Right).Render("Hangman"),
	// ))

	doc.WriteString(wordBox + "\n\n")
	doc.WriteString(fmt.Sprintf("Word: %s\n%s\n", g.model.TextInput.View(), "(esc to exit)"))

	return doc.String()

}

func (g Game) Init() tea.Cmd {
	return textinput.Blink
}
