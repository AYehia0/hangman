package game

import (
	"strings"

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
	currentGuess string
	maxGuesses   int
	keyboard     map[string]selected
	model        model
}

type model struct {
	Tabs       []string
	TabContent []string // change this to include a view
	ActiveTab  int      // index of the currect active tab
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func Play(width int, height int, session ssh.Session) Game {
	// create a new game for that player
	var game Game

	// get a random word
	word, err := GetRandomWord()
	if err != nil {
		log.Errorf("Couldn't fetch a word : %s", err)
	}

	game.word = word
	game.session = session
	game.model.Tabs = []string{"Game", "Results"}
	game.model.TabContent = []string{"Game Content go here", "Results go here!"}

	return game
}

// bubbletea
func (g Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return g, tea.Quit
		case "right", "tab":
			g.model.ActiveTab = min(g.model.ActiveTab+1, len(g.model.Tabs)-1)
			return g, nil
		case "left", "shift+tab":
			g.model.ActiveTab = max(g.model.ActiveTab-1, 0)
			return g, nil
		}
	}

	return g, nil
}

func (g Game) View() string {
	doc := strings.Builder{}
	m := g.model

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.ActiveTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.ActiveTab]))
	return docStyle.Render(doc.String())
}

func (g Game) Init() tea.Cmd {
	return nil
}
