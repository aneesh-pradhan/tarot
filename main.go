// Command tarot is a terminal tarot reading game.
package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	seed := flag.Uint64("seed", 0, "seed for a reproducible reading (0 = random)")
	noReversals := flag.Bool("no-reversals", false, "never draw reversed cards")
	flag.Parse()

	s1, s2 := rand.Uint64(), rand.Uint64()
	if *seed != 0 {
		s1, s2 = *seed, *seed
	}
	rng := rand.New(rand.NewPCG(s1, s2))

	p := tea.NewProgram(newModel(rng, !*noReversals), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "tarot:", err)
		os.Exit(1)
	}
}
