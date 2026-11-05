package main

import (
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateMenu state = iota
	stateQuestion
	stateShuffle
	stateReading
)

const (
	shuffleFrames = 18
	shuffleTick   = 70 * time.Millisecond
)

// tickMsg advances the shuffle animation. The id ties it to one shuffle so a
// leftover tick from a skipped animation can't speed up the next one.
type tickMsg struct{ id int }

func tick(id int) tea.Cmd {
	return tea.Tick(shuffleTick, func(time.Time) tea.Msg { return tickMsg{id} })
}

type model struct {
	state         state
	width, height int

	rng       *rand.Rand
	reversals bool

	menuIdx  int
	input    textinput.Model
	question string

	spread   Spread
	cards    []Drawn
	revealed []bool
	selected int

	frame     int
	shuffleID int
}

func newModel(rng *rand.Rand, reversals bool) model {
	ti := textinput.New()
	ti.Placeholder = "What should I focus on right now?"
	ti.CharLimit = 120
	ti.Width = 50
	ti.Prompt = "› "
	return model{rng: rng, reversals: reversals, input: ti, width: 80, height: 24}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tickMsg:
		if m.state != stateShuffle || msg.id != m.shuffleID {
			return m, nil
		}
		m.frame++
		if m.frame >= shuffleFrames {
			m.deal()
			return m, nil
		}
		return m, tick(m.shuffleID)

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch m.state {
		case stateMenu:
			return m.updateMenu(msg)
		case stateQuestion:
			return m.updateQuestion(msg)
		case stateShuffle:
			m.deal() // any key skips the animation
			return m, nil
		case stateReading:
			return m.updateReading(msg)
		}
	}

	// Forward everything else (e.g. cursor blink) to the text input.
	if m.state == stateQuestion {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k := msg.String(); k {
	case "q", "esc":
		return m, tea.Quit
	case "up", "k":
		if m.menuIdx > 0 {
			m.menuIdx--
		}
	case "down", "j":
		if m.menuIdx < len(spreads)-1 {
			m.menuIdx++
		}
	case "r":
		m.reversals = !m.reversals
	case "enter", " ":
		return m.chooseSpread()
	default:
		if n, err := strconv.Atoi(k); err == nil && n >= 1 && n <= len(spreads) {
			m.menuIdx = n - 1
			return m.chooseSpread()
		}
	}
	return m, nil
}

func (m model) chooseSpread() (tea.Model, tea.Cmd) {
	m.spread = spreads[m.menuIdx]
	m.state = stateQuestion
	m.input.SetValue("")
	return m, m.input.Focus()
}

func (m model) updateQuestion(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.input.Blur()
		m.state = stateMenu
		return m, nil
	case tea.KeyEnter:
		m.question = strings.TrimSpace(m.input.Value())
		m.input.Blur()
		m.state = stateShuffle
		m.frame = 0
		m.shuffleID++
		return m, tick(m.shuffleID)
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *model) deal() {
	m.cards = Draw(m.rng, len(m.spread.Positions), m.reversals)
	m.revealed = make([]bool, len(m.cards))
	m.selected = 0
	m.state = stateReading
}

func (m model) updateReading(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	n := len(m.cards)
	per := m.perRow()
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "esc", "n":
		m.state = stateMenu
	case "left", "h", "shift+tab":
		if m.selected > 0 {
			m.selected--
		}
	case "right", "l", "tab":
		if m.selected < n-1 {
			m.selected++
		}
	case "up", "k":
		if m.selected-per >= 0 {
			m.selected -= per
		}
	case "down", "j":
		if m.selected+per < n {
			m.selected += per
		}
	case " ", "enter":
		// Turn over the selected card, or if it's already up, move to and
		// turn over the next face-down card.
		if !m.revealed[m.selected] {
			m.revealed[m.selected] = true
		} else if i := m.nextHidden(); i >= 0 {
			m.selected = i
			m.revealed[i] = true
		}
	case "a":
		for i := range m.revealed {
			m.revealed[i] = true
		}
	}
	return m, nil
}

// nextHidden returns the first face-down card after the selection, wrapping
// around, or -1 if every card is face up.
func (m model) nextHidden() int {
	n := len(m.revealed)
	for off := 1; off <= n; off++ {
		if i := (m.selected + off) % n; !m.revealed[i] {
			return i
		}
	}
	return -1
}

func (m model) revealedCount() int {
	c := 0
	for _, r := range m.revealed {
		if r {
			c++
		}
	}
	return c
}
