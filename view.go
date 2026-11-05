package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	cardInner = 11            // interior width of a card
	cardWidth = cardInner + 2 // including border
	cellWidth = cardWidth + 2 // card plus gap in the grid
)

var (
	gold   = lipgloss.AdaptiveColor{Light: "#9A6B00", Dark: "#E6C35C"}
	purple = lipgloss.AdaptiveColor{Light: "#5B3E96", Dark: "#B39DDB"}
	dim    = lipgloss.AdaptiveColor{Light: "#6B6B6B", Dark: "#8A8A8A"}
	red    = lipgloss.AdaptiveColor{Light: "#B00020", Dark: "#EF5350"}
	green  = lipgloss.AdaptiveColor{Light: "#2E7D32", Dark: "#81C784"}

	suitColors = map[Suit]lipgloss.AdaptiveColor{
		Major:     purple,
		Wands:     {Light: "#B5471B", Dark: "#FF8A50"},
		Cups:      {Light: "#1D5FA8", Dark: "#64B5F6"},
		Swords:    {Light: "#4A6572", Dark: "#B0BEC5"},
		Pentacles: green,
	}

	titleStyle  = lipgloss.NewStyle().Foreground(gold).Bold(true)
	dimStyle    = lipgloss.NewStyle().Foreground(dim)
	italicStyle = lipgloss.NewStyle().Foreground(dim).Italic(true)
	revStyle    = lipgloss.NewStyle().Foreground(red)
	upStyle     = lipgloss.NewStyle().Foreground(green)
	helpStyle   = lipgloss.NewStyle().Foreground(dim).MarginTop(1)
	panelStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(dim).Padding(0, 1)
)

func (m model) View() string {
	var body string
	switch m.state {
	case stateMenu:
		body = m.viewMenu()
	case stateQuestion:
		body = m.viewQuestion()
	case stateShuffle:
		body = m.viewShuffle()
	case stateReading:
		body = m.viewReading()
	}
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, body)
}

func banner() string {
	return titleStyle.Render("✦  T A R O T  ✦") + "\n" + italicStyle.Render("a terminal oracle")
}

func (m model) viewMenu() string {
	var b strings.Builder
	b.WriteString(banner() + "\n\n")
	b.WriteString("Choose a spread:\n\n")
	for i, s := range spreads {
		count := fmt.Sprintf("%2d card", len(s.Positions))
		if len(s.Positions) != 1 {
			count += "s"
		}
		line := fmt.Sprintf("%d  %-30s %s", i+1, s.Name, dimStyle.Render(count))
		if i == m.menuIdx {
			b.WriteString(titleStyle.Render("› ") + lipgloss.NewStyle().Bold(true).Render(line) + "\n")
		} else {
			b.WriteString("  " + line + "\n")
		}
	}
	b.WriteString("\n" + italicStyle.Render(spreads[m.menuIdx].Desc) + "\n\n")

	rev := revStyle.Render("off")
	if m.reversals {
		rev = upStyle.Render("on")
	}
	b.WriteString("Reversed cards: " + rev)
	b.WriteString(helpStyle.Render("\n↑/↓ choose · enter select · r toggle reversals · q quit"))
	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}

func (m model) viewQuestion() string {
	var b strings.Builder
	b.WriteString(banner() + "\n\n")
	b.WriteString(titleStyle.Render(m.spread.Name) + "\n\n")
	b.WriteString("Hold a question in mind, or type it below (optional).\n\n")
	b.WriteString(m.input.View() + "\n")
	b.WriteString(helpStyle.Render("enter shuffle the deck · esc back"))
	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}

func (m model) viewShuffle() string {
	bob := []int{0, 1, 2, 1}
	backs := make([]string, 5)
	for i := range backs {
		off := bob[(m.frame+i)%len(bob)]
		backs[i] = strings.Repeat("\n", off) + renderCardBack(false)
	}
	stack := lipgloss.JoinHorizontal(lipgloss.Top, backs...)
	dots := strings.Repeat("·", m.frame%4+1)
	text := titleStyle.Render("Shuffling the deck "+dots) + "\n" + dimStyle.Render("press any key to skip")
	return lipgloss.NewStyle().Padding(2, 2).Render(
		lipgloss.JoinVertical(lipgloss.Center, stack, "", text))
}

func (m model) perRow() int {
	return max(1, (m.width-2)/cellWidth)
}

func (m model) viewReading() string {
	if v := m.readingLayout(false); lipgloss.Height(v) <= m.height {
		return v
	}
	return m.readingLayout(true)
}

// readingLayout renders the table. In compact mode it drops spacing and the
// suit blurb, and if the card grid still doesn't fit it shows only the rows
// around the selected card.
func (m model) readingLayout(compact bool) string {
	header := titleStyle.Render(m.spread.Name)
	q := "an open reading"
	if m.question != "" {
		q = "“" + m.question + "”"
	}
	if compact {
		header += dimStyle.Render("  ·  ") + italicStyle.Render(q)
	} else {
		header += "\n" + italicStyle.Render(q)
	}

	panelWidth := max(30, min(m.width-4, 88))
	panel := panelStyle.Width(panelWidth).Render(m.detail(m.selected, compact))

	done := m.revealedCount()
	help := fmt.Sprintf("%d/%d revealed · ←/→ select · space turn card · a reveal all · n new reading · q quit",
		done, len(m.cards))
	if done == len(m.cards) {
		help = "The reading is complete. ←/→ revisit cards · n new reading · q quit"
	}

	// Lay the cards out in rows that fit the terminal width.
	per := m.perRow()
	totalRows := (len(m.cards) + per - 1) / per
	firstRow, lastRow := 0, totalRows
	if compact {
		const rowHeight = 8 // label + 7-line card
		fixed := lipgloss.Height(header) + lipgloss.Height(panel) + 2 + 1
		fit := max(1, (m.height-fixed)/rowHeight)
		if fit < totalRows {
			sel := m.selected / per
			firstRow = min(max(0, sel-fit+1), totalRows-fit)
			lastRow = firstRow + fit
			help = fmt.Sprintf("rows %d–%d of %d (↑/↓) · ", firstRow+1, lastRow, totalRows) + help
		}
	}

	var rows []string
	for r := firstRow; r < lastRow; r++ {
		start := r * per
		end := min(start+per, len(m.cards))
		cells := make([]string, 0, end-start)
		for i := start; i < end; i++ {
			label := truncate(fmt.Sprintf("%d. %s", i+1, m.spread.Positions[i].Name), cellWidth-1)
			if i == m.selected {
				label = titleStyle.Render(label)
			} else {
				label = dimStyle.Render(label)
			}
			card := renderCard(m.cards[i], m.revealed[i], i == m.selected)
			cells = append(cells, lipgloss.NewStyle().Width(cellWidth).Render(label+"\n"+card))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}
	grid := lipgloss.JoinVertical(lipgloss.Left, rows...)

	if compact {
		help = lipgloss.NewStyle().Foreground(dim).Width(m.width - 2).Render(help)
		return lipgloss.NewStyle().Padding(0, 1).Render(
			lipgloss.JoinVertical(lipgloss.Left, header, grid, panel, help))
	}
	return lipgloss.NewStyle().Padding(1, 1).Render(
		lipgloss.JoinVertical(lipgloss.Left, header, "", grid, "", panel, helpStyle.Render(help)))
}

func (m model) detail(i int, compact bool) string {
	gap := "\n\n"
	if compact {
		gap = "\n"
	}
	pos := m.spread.Positions[i]
	head := titleStyle.Render(fmt.Sprintf("%d · %s", i+1, pos.Name)) + dimStyle.Render(" — "+pos.Desc)
	if !m.revealed[i] {
		return head + gap + italicStyle.Render("This card is still face down. Press space to turn it over.")
	}

	d := m.cards[i]
	nameStyle := lipgloss.NewStyle().Foreground(suitColors[d.Card.Suit]).Bold(true)
	orient, otherLabel, reversed := upStyle.Render("Upright"), "Reversed", ""
	if d.Reversed {
		orient, otherLabel, reversed = revStyle.Render("Reversed"), "Upright", " reversed"
	}

	article := ""
	if d.Card.Suit != Major {
		article = "the "
	}

	var b strings.Builder
	b.WriteString(head + gap)
	b.WriteString(nameStyle.Render(d.Card.Name) + "  " + orient + dimStyle.Render("  ·  "+d.Card.Label()) + gap)
	b.WriteString(fmt.Sprintf("In the %s position, %s%s%s speaks of %s.",
		pos.Name, article, d.Card.Name, reversed, d.Meaning()) + gap)
	b.WriteString(dimStyle.Render(otherLabel + ": " + d.OtherMeaning()))
	if !compact {
		b.WriteString("\n" + italicStyle.Render(d.Card.Suit.Theme()))
	}
	return b.String()
}

// renderCard draws a card face up or face down. Every card renders to exactly
// cardWidth × 7 cells so the grid stays aligned.
func renderCard(d Drawn, faceUp, selected bool) string {
	if !faceUp {
		return renderCardBack(selected)
	}
	color := suitColors[d.Card.Suit]
	border := lipgloss.RoundedBorder()
	if selected {
		border = lipgloss.DoubleBorder()
	}
	symStyle := lipgloss.NewStyle().Foreground(color).Bold(true)
	nameStyle := lipgloss.NewStyle().Bold(true)

	n1, n2 := wrap2(d.Card.Name, cardInner)
	sym := center(symStyle.Render(d.Card.Suit.Symbol()))
	var lines []string
	if d.Reversed {
		lines = []string{
			center(revStyle.Render("reversed")),
			sym,
			center(nameStyle.Render(n1)),
			center(nameStyle.Render(n2)),
			lipgloss.PlaceHorizontal(cardInner, lipgloss.Right, dimStyle.Render(d.Card.Label())),
		}
	} else {
		lines = []string{
			lipgloss.PlaceHorizontal(cardInner, lipgloss.Left, dimStyle.Render(d.Card.Label())),
			"",
			sym,
			center(nameStyle.Render(n1)),
			center(nameStyle.Render(n2)),
		}
	}
	for i, l := range lines {
		lines[i] = lipgloss.PlaceHorizontal(cardInner, lipgloss.Left, l)
	}
	bc := color
	if selected {
		bc = gold
	}
	return lipgloss.NewStyle().Border(border).BorderForeground(bc).Render(strings.Join(lines, "\n"))
}

func renderCardBack(selected bool) string {
	a := "╱╲╱╲╱╲╱╲╱╲╱"
	b := "╲╱╲╱╲╱╲╱╲╱╲"
	pattern := lipgloss.NewStyle().Foreground(purple).Faint(true)
	body := strings.Join([]string{
		pattern.Render(a), pattern.Render(b),
		pattern.Render("╱╲╱╲ ") + titleStyle.Render("✦") + pattern.Render(" ╲╱╲╱"),
		pattern.Render(b), pattern.Render(a),
	}, "\n")
	border, bc := lipgloss.RoundedBorder(), lipgloss.TerminalColor(purple)
	if selected {
		border, bc = lipgloss.DoubleBorder(), gold
	}
	return lipgloss.NewStyle().Border(border).BorderForeground(bc).Render(body)
}

func center(s string) string {
	return lipgloss.PlaceHorizontal(cardInner, lipgloss.Center, s)
}

// wrap2 word-wraps s onto at most two lines of the given width.
func wrap2(s string, width int) (string, string) {
	words := strings.Fields(s)
	first := ""
	i := 0
	for ; i < len(words); i++ {
		next := words[i]
		if first != "" {
			next = first + " " + words[i]
		}
		if lipgloss.Width(next) > width && first != "" {
			break
		}
		first = next
	}
	return truncate(first, width), truncate(strings.Join(words[i:], " "), width)
}

func truncate(s string, width int) string {
	r := []rune(s)
	if lipgloss.Width(s) <= width {
		return s
	}
	for len(r) > 0 && lipgloss.Width(string(r)) > width-1 {
		r = r[:len(r)-1]
	}
	return string(r) + "…"
}
