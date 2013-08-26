package regex

type node struct {
	edges map[rune][]*node
}

type nfa struct {
	finals []*node
	start  *node
}

// Copy b's transitions to a
func copyTransitions(a, b *node) {
	for label, nextStates := range b.edges {
		a.edges[label] = append(a.edges[label], nextStates...)
	}
}

// Concats an NFA in place with another.
//
// This is acheived by copying the out-arrows from b's start state, and adding them to a's final states.
func (a *nfa) concat(b *nfa) {
	for _, finalNode := range a.finals {
		copyTransitions(finalNode, b.start)
	}

	a.finals = b.finals
}

// "OR"s an NFA in place with another.
//
// This is acheived by creating a new start state and copying all of a's and b's start state transitions to this new start state.
func (a *nfa) or(b *nfa) {
	newStart := &node{map[rune][]*node{}}
	copyTransitions(newStart, a.start)
	copyTransitions(newStart, b.start)

	a.start = newStart
	a.finals = append(a.finals, b.finals...)
}

// Modifies the NFA to optionally skip processing.
//
// We achieve this by simply making the start state final.
func (a *nfa) makeOptional() {
	// check if the start state is already final (in which case there's nothing to do)
	for _, n := range a.finals {
		if n == a.start {
			return
		}
	}

	// otherwise just add it
	a.finals = append(a.finals, a.start)
}

// Modifies the NFA to loop on itself 0 or more times (ie, the empty string is valid as well)
func (a *nfa) loop() {
	// this ensures the start state is final
	a.makeOptional()

	// now we just hook up the final states to the start state
	// (except we avoid a self-loop for the start state itself)
	for i := range a.finals {
		if a.finals[i] != a.start {
			copyTransitions(a.finals[i], a.start)
		}
	}
}

// Creates an NFA for matching a single character
func characterNFA(r rune) *nfa {
	final := &node{map[rune][]*node{}}
	return &nfa{
		finals: []*node{final},
		start: &node{
			map[rune][]*node{r: []*node{final}},
		},
	}
}

func (m *nfa) process(s string) bool {
	activeStates := []*node{m.start}
	for _, c := range s {
		nextStates := []*node{}
		for _, currentState := range activeStates {
			if nextState, ok := currentState.edges[c]; ok {
				nextStates = append(nextStates, nextState...)
			}
		}
		activeStates = nextStates
	}

	// are any of the active states final? (ugly double for-loop atm)
	for _, final := range m.finals {
		for _, active := range activeStates {
			if final == active {
				return true
			}
		}
	}

	return false
}
