package main

import (
	"math/rand/v2"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestDeckIsComplete(t *testing.T) {
	deck := NewDeck()
	if len(deck) != 78 {
		t.Fatalf("deck has %d cards, want 78", len(deck))
	}
	seen := map[string]bool{}
	for _, c := range deck {
		if seen[c.Name] {
			t.Errorf("duplicate card %q", c.Name)
		}
		seen[c.Name] = true
		if c.Upright == "" || c.Reversed == "" {
			t.Errorf("%s is missing a meaning", c.Name)
		}
	}
}

func TestDrawIsUniqueAndSeeded(t *testing.T) {
	a := Draw(rand.New(rand.NewPCG(7, 7)), 10, true)
	b := Draw(rand.New(rand.NewPCG(7, 7)), 10, true)
	seen := map[string]bool{}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("same seed gave different draws at %d", i)
		}
		if seen[a[i].Card.Name] {
			t.Fatalf("card %q drawn twice", a[i].Card.Name)
		}
		seen[a[i].Card.Name] = true
	}
	for _, d := range Draw(rand.New(rand.NewPCG(1, 2)), 78, false) {
		if d.Reversed {
			t.Fatal("got a reversed card with reversals disabled")
		}
	}
}

// Every card must render to the same box so the grid lines up.
func TestCardsRenderToFixedSize(t *testing.T) {
	check := func(name, s string) {
		t.Helper()
		if w, h := lipgloss.Width(s), lipgloss.Height(s); w != cardWidth || h != 7 {
			t.Errorf("%s renders %dx%d, want %dx7", name, w, h, cardWidth)
		}
		for _, line := range strings.Split(s, "\n") {
			if lipgloss.Width(line) != cardWidth {
				t.Errorf("%s has a ragged line %q", name, line)
			}
		}
	}
	check("card back", renderCardBack(false))
	for _, c := range NewDeck() {
		for _, rev := range []bool{false, true} {
			for _, sel := range []bool{false, true} {
				check(c.Name, renderCard(Drawn{c, rev}, true, sel))
			}
		}
	}
}

func key(s string) tea.KeyMsg {
	switch s {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case " ":
		return tea.KeyMsg{Type: tea.KeySpace, Runes: []rune(" ")}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func TestFullReadingFlow(t *testing.T) {
	var tm tea.Model = newModel(rand.New(rand.NewPCG(3, 3)), true)
	tm, _ = tm.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	for _, k := range []string{"5", "enter", "x"} { // Celtic Cross, no question, skip shuffle
		tm, _ = tm.Update(key(k))
	}
	m := tm.(model)
	if m.state != stateReading || len(m.cards) != 10 {
		t.Fatalf("state=%v cards=%d, want reading with 10 cards", m.state, len(m.cards))
	}
	for range 10 {
		tm, _ = tm.Update(key(" "))
	}
	m = tm.(model)
	if m.revealedCount() != 10 {
		t.Fatalf("revealed %d cards after 10 presses, want 10", m.revealedCount())
	}
	if v := m.View(); !strings.Contains(v, "reading is complete") || !strings.Contains(v, m.cards[9].Card.Name) {
		t.Errorf("final view missing completion text or last card:\n%s", v)
	}
}

// Readings must fit on screen, including a classic 80×24 terminal.
func TestReadingFitsTerminal(t *testing.T) {
	for _, size := range [][2]int{{80, 24}, {100, 30}, {160, 50}} {
		for si := range spreads {
			var tm tea.Model = newModel(rand.New(rand.NewPCG(5, 5)), true)
			tm, _ = tm.Update(tea.WindowSizeMsg{Width: size[0], Height: size[1]})
			m := tm.(model)
			m.menuIdx = si
			tm, _ = m.chooseSpread()
			tm, _ = tm.Update(key("enter"))
			tm, _ = tm.Update(key("x"))
			for i := 0; i < len(spreads[si].Positions); i++ {
				tm, _ = tm.Update(key(" "))
				v := tm.View()
				if h := lipgloss.Height(v); h > size[1] {
					t.Errorf("%s at %dx%d, card %d: view is %d lines tall", spreads[si].Name, size[0], size[1], i+1, h)
				}
				if !strings.Contains(v, spreads[si].Name) {
					t.Errorf("%s at %dx%d: header missing", spreads[si].Name, size[0], size[1])
				}
			}
		}
	}
}
