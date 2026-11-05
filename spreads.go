package main

// Position is one slot in a spread.
type Position struct {
	Name string
	Desc string
}

// Spread is a named layout of positions.
type Spread struct {
	Name      string
	Desc      string
	Positions []Position
}

var spreads = []Spread{
	{
		Name: "Single Card",
		Desc: "One card for quick guidance or a daily draw.",
		Positions: []Position{
			{"Guidance", "the energy or advice most relevant to you right now"},
		},
	},
	{
		Name: "Past · Present · Future",
		Desc: "A classic three-card look at how you got here and where you're going.",
		Positions: []Position{
			{"Past", "influences and events that led to this moment"},
			{"Present", "where you stand and what surrounds you now"},
			{"Future", "where things are heading if the current path continues"},
		},
	},
	{
		Name: "Situation · Action · Outcome",
		Desc: "Three cards for a practical problem: what's happening, what to do, what follows.",
		Positions: []Position{
			{"Situation", "the heart of the matter as it stands"},
			{"Action", "the approach the cards recommend"},
			{"Outcome", "the likely result of taking that action"},
		},
	},
	{
		Name: "Mind · Body · Spirit",
		Desc: "A three-card check-in on your whole self.",
		Positions: []Position{
			{"Mind", "your thoughts, focus and mental state"},
			{"Body", "your physical wellbeing and material circumstances"},
			{"Spirit", "your inner life, purpose and sense of meaning"},
		},
	},
	{
		Name: "Celtic Cross",
		Desc: "The traditional ten-card spread for a deep, detailed reading.",
		Positions: []Position{
			{"Present", "the heart of the situation"},
			{"Challenge", "what crosses you; the obstacle to work with"},
			{"Foundation", "the root cause beneath the situation"},
			{"Recent Past", "what is passing out of your life"},
			{"Crown", "your conscious goal or best possible outcome"},
			{"Near Future", "what is coming into your life soon"},
			{"Self", "your attitude and how you see yourself"},
			{"Environment", "the people and forces around you"},
			{"Hopes & Fears", "what you long for and what you dread"},
			{"Outcome", "where all of this is leading"},
		},
	},
}
