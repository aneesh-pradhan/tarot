package main

import (
	"fmt"
	"math/rand/v2"
)

// Suit identifies which part of the deck a card belongs to.
type Suit int

const (
	Major Suit = iota
	Wands
	Cups
	Swords
	Pentacles
)

func (s Suit) String() string {
	return [...]string{"Major Arcana", "Wands", "Cups", "Swords", "Pentacles"}[s]
}

// Symbol is the glyph drawn in the middle of a face-up card. The minor suits
// use their playing-card descendants (Wands→clubs, Cups→hearts, and so on).
func (s Suit) Symbol() string {
	return [...]string{"✦", "♣", "♥", "♠", "♦"}[s]
}

// Theme is a one-line description of what the suit governs.
func (s Suit) Theme() string {
	return [...]string{
		"Major Arcana cards mark significant life lessons, turning points and deep currents.",
		"Wands · Fire — passion, ambition, creativity and the drive to act.",
		"Cups · Water — emotions, relationships, intuition and the heart.",
		"Swords · Air — thought, truth, conflict and communication.",
		"Pentacles · Earth — work, money, health and the material world.",
	}[s]
}

// Card is a single tarot card with keyword meanings for each orientation.
type Card struct {
	Name     string
	Suit     Suit
	Rank     int // 0–21 for the Major Arcana, 1–14 (Ace…King) for minors
	Upright  string
	Reversed string
}

var romanNumerals = []string{
	"0", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X",
	"XI", "XII", "XIII", "XIV", "XV", "XVI", "XVII", "XVIII", "XIX", "XX", "XXI",
}

var rankNames = []string{
	"", "Ace", "Two", "Three", "Four", "Five", "Six", "Seven",
	"Eight", "Nine", "Ten", "Page", "Knight", "Queen", "King",
}

var rankLabels = []string{
	"", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Pg", "Kn", "Q", "K",
}

// Label is the short corner index printed on the card.
func (c Card) Label() string {
	if c.Suit == Major {
		return romanNumerals[c.Rank]
	}
	return rankLabels[c.Rank] + " " + c.Suit.Symbol()
}

// Drawn is a card as it lands on the table.
type Drawn struct {
	Card     Card
	Reversed bool
}

// Meaning returns the keywords for the orientation the card was drawn in.
func (d Drawn) Meaning() string {
	if d.Reversed {
		return d.Card.Reversed
	}
	return d.Card.Upright
}

// OtherMeaning returns the keywords for the opposite orientation.
func (d Drawn) OtherMeaning() string {
	if d.Reversed {
		return d.Card.Upright
	}
	return d.Card.Reversed
}

type meaning struct{ name, up, rev string }

var majorArcana = []meaning{
	{"The Fool", "beginnings, innocence, spontaneity, a leap of faith", "recklessness, hesitation, naivety, needless risk"},
	{"The Magician", "manifestation, willpower, skill, resourcefulness", "manipulation, untapped talent, trickery"},
	{"The High Priestess", "intuition, mystery, the subconscious, your inner voice", "secrets, ignored intuition, withdrawal"},
	{"The Empress", "abundance, nurturing, fertility, creativity", "dependence, creative block, smothering"},
	{"The Emperor", "authority, structure, stability, leadership", "rigidity, domination, lack of discipline"},
	{"The Hierophant", "tradition, conformity, spiritual guidance, institutions", "rebellion, unconventional paths, personal beliefs"},
	{"The Lovers", "love, harmony, partnership, meaningful choices", "imbalance, disharmony, misaligned values"},
	{"The Chariot", "determination, control, victory, momentum", "lack of direction, aggression, self-doubt"},
	{"Strength", "courage, compassion, inner strength, patience", "self-doubt, insecurity, raw emotion"},
	{"The Hermit", "introspection, solitude, inner guidance, soul-searching", "isolation, loneliness, withdrawal"},
	{"Wheel of Fortune", "cycles, fate, turning points, luck", "bad luck, resistance to change, broken cycles"},
	{"Justice", "fairness, truth, cause and effect, law", "unfairness, dishonesty, avoiding accountability"},
	{"The Hanged Man", "surrender, pause, new perspective, letting go", "stalling, resistance, indecision"},
	{"Death", "endings, transformation, transition, release", "resisting change, stagnation, fear of endings"},
	{"Temperance", "balance, moderation, patience, purpose", "imbalance, excess, lack of harmony"},
	{"The Devil", "attachment, temptation, materialism, the shadow self", "release, breaking free, reclaiming power"},
	{"The Tower", "sudden upheaval, revelation, chaos, awakening", "averted disaster, fear of change, delayed collapse"},
	{"The Star", "hope, renewal, inspiration, serenity", "despair, lost faith, disconnection"},
	{"The Moon", "illusion, fear, intuition, the unknown", "clarity returning, released fear, confusion lifting"},
	{"The Sun", "joy, success, vitality, positivity", "temporary gloom, overconfidence, dimmed enthusiasm"},
	{"Judgement", "rebirth, reckoning, an inner calling, absolution", "self-doubt, refusing the call, harsh self-judgement"},
	{"The World", "completion, fulfilment, integration, travel", "incompletion, shortcuts, lack of closure"},
}

// minorArcana holds [upright, reversed] keywords for Ace through King.
var minorArcana = map[Suit][14][2]string{
	Wands: {
		{"inspiration, a new venture, creative spark, potential", "delays, lack of motivation, a false start"},
		{"planning, future vision, decisions, discovery", "fear of the unknown, poor planning, playing it safe"},
		{"expansion, foresight, progress, opportunity abroad", "obstacles, frustrating delays, limited vision"},
		{"celebration, homecoming, harmony, milestones", "instability, transition, lack of support"},
		{"conflict, competition, tension, disagreement", "avoiding conflict, resolution, inner conflict"},
		{"victory, recognition, public success, confidence", "ego, a fall from grace, lack of recognition"},
		{"defending your position, perseverance, challenge", "overwhelm, giving up, exhaustion"},
		{"swift action, movement, rapid progress, news", "delays, frustration, slowing down"},
		{"resilience, persistence, a last stand, boundaries", "paranoia, fatigue, defensiveness"},
		{"burden, responsibility, hard work, overload", "delegating, release, burnout"},
		{"enthusiasm, exploration, a free spirit, discovery", "lack of direction, procrastination, setbacks"},
		{"energy, passion, adventure, impulsiveness", "haste, scattered energy, frustration"},
		{"confidence, determination, warmth, vibrancy", "jealousy, insecurity, selfishness"},
		{"visionary leadership, entrepreneurship, honour", "impulsiveness, being overbearing, high expectations"},
	},
	Cups: {
		{"new love, compassion, emotional awakening, creativity", "emotional blockage, repressed feelings, emptiness"},
		{"partnership, mutual attraction, unity, connection", "imbalance, broken communication, tension"},
		{"friendship, celebration, community, joy", "overindulgence, gossip, isolation"},
		{"apathy, contemplation, missed opportunity, discontent", "renewed interest, acceptance, a new perspective"},
		{"loss, grief, regret, disappointment", "acceptance, moving on, forgiveness"},
		{"nostalgia, childhood memories, innocence, reunion", "living in the past, rose-tinted memories, moving forward"},
		{"choices, illusion, fantasy, wishful thinking", "clarity, alignment, a sober decision"},
		{"walking away, disillusionment, seeking deeper meaning", "fear of change, aimless drifting, staying put"},
		{"contentment, satisfaction, wishes fulfilled", "smugness, dissatisfaction, materialism"},
		{"harmony, family, lasting happiness, alignment", "a broken home, misaligned values, disconnection"},
		{"creative opportunity, curiosity, intuitive messages", "emotional immaturity, creative block, insecurity"},
		{"romance, charm, following the heart, idealism", "moodiness, unrealistic expectations, disappointment"},
		{"compassion, calm, emotional security, intuition", "codependency, insecurity, emotional overwhelm"},
		{"emotional balance, diplomacy, generosity, wisdom", "manipulation, moodiness, coldness"},
	},
	Swords: {
		{"clarity, breakthrough, truth, a sharp mind", "confusion, miscommunication, clouded judgement"},
		{"a difficult choice, stalemate, indecision, avoidance", "information overload, confusion, the lesser of two evils"},
		{"heartbreak, sorrow, grief, a painful truth", "recovery, forgiveness, releasing pain"},
		{"rest, recuperation, contemplation, retreat", "restlessness, burnout, stagnation"},
		{"conflict, hollow victory, tension, winning at all costs", "reconciliation, making amends, lingering resentment"},
		{"transition, moving on, leaving behind, calmer waters", "resistance to change, unfinished business, baggage"},
		{"deception, strategy, stealth, getting away with it", "confession, conscience, coming clean"},
		{"feeling trapped, restriction, self-imposed limits", "release, a new perspective, self-acceptance"},
		{"anxiety, worry, nightmares, despair", "hope, reaching out, recovery from anxiety"},
		{"a painful ending, rock bottom, betrayal, finality", "recovery, regeneration, resisting the inevitable"},
		{"curiosity, new ideas, vigilance, a thirst for knowledge", "gossip, all talk, scattered thoughts"},
		{"ambition, action, fast thinking, drive", "recklessness, impatience, no direction"},
		{"independence, clear boundaries, direct communication", "coldness, cruelty, bitterness"},
		{"intellectual power, authority, truth, ethics", "manipulation, tyranny, misuse of power"},
	},
	Pentacles: {
		{"a new opportunity, prosperity, manifestation, abundance", "a lost opportunity, poor planning, scarcity"},
		{"balance, adaptability, juggling priorities", "overwhelm, disorganisation, overcommitment"},
		{"teamwork, collaboration, skill, craftsmanship", "disharmony, poor teamwork, misalignment"},
		{"security, control, saving, conservatism", "greed, possessiveness, letting go of control"},
		{"hardship, loss, isolation, financial worry", "recovery, help arriving, the end of hard times"},
		{"generosity, charity, giving and receiving", "strings attached, debt, one-sided giving"},
		{"patience, a long-term view, investment, perseverance", "impatience, limited reward, wasted effort"},
		{"diligence, mastery, skill development, dedication", "perfectionism, lack of focus, uninspired work"},
		{"independence, luxury, self-sufficiency, reward", "overwork, superficiality, financial setbacks"},
		{"wealth, legacy, family, long-term success", "financial failure, a lost legacy, family disputes"},
		{"ambition, diligence, a new skill, manifestation", "lack of progress, procrastination, missed chances"},
		{"efficiency, routine, reliability, hard work", "boredom, stagnation, laziness"},
		{"nurturing, practicality, comfort, financial security", "self-neglect, work-home imbalance, smothering"},
		{"abundance, security, discipline, business acumen", "greed, stubbornness, indulgence"},
	},
}

// NewDeck returns the full 78-card Rider–Waite–Smith style deck in order.
func NewDeck() []Card {
	deck := make([]Card, 0, 78)
	for i, m := range majorArcana {
		deck = append(deck, Card{Name: m.name, Suit: Major, Rank: i, Upright: m.up, Reversed: m.rev})
	}
	for _, s := range []Suit{Wands, Cups, Swords, Pentacles} {
		for i, kw := range minorArcana[s] {
			rank := i + 1
			deck = append(deck, Card{
				Name:     fmt.Sprintf("%s of %s", rankNames[rank], s),
				Suit:     s,
				Rank:     rank,
				Upright:  kw[0],
				Reversed: kw[1],
			})
		}
	}
	return deck
}

// Draw shuffles a fresh deck and deals n cards from the top. When reversals
// is true each card has an even chance of landing upside down.
func Draw(r *rand.Rand, n int, reversals bool) []Drawn {
	deck := NewDeck()
	r.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	n = min(n, len(deck))
	out := make([]Drawn, n)
	for i := range out {
		out[i] = Drawn{Card: deck[i], Reversed: reversals && r.IntN(2) == 1}
	}
	return out
}
