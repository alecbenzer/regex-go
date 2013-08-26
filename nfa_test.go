package regex

import (
	"testing"
)

func TestSingleCharacter(t *testing.T) {
	for c := rune(0); c < 256; c++ {
		if !characterNFA(c).process(string(c)) {
			t.Error("Failed to match rune", c)
		}
		if characterNFA(c).process(string(rune(c + 1))) {
			t.Error("Matched rune", c+1)
		}
	}

	// try some unicode characters
	for c := rune(10000); c < 10000 + 256; c++ {
		if !characterNFA(c).process(string(c)) {
			t.Error("Failed to match rune", c)
		}
		if characterNFA(c).process(string(rune(c + 1))) {
			t.Error("Matched rune", c+1)
		}
	}
}

func TestConcat(t *testing.T) {
	a := characterNFA('a')
	b := characterNFA('b')
	a.concat(b)

	bad := []string{
		"a",
		"b",
		"",
	}

	if !a.process("ab") {
		t.Errorf("Failed to match %#v\n", "ab")
	}

	for _, s := range bad {
		if a.process(s) {
			t.Errorf("Matched %#v, but should not have\n", s)
		}
	}
}

func TestMultipleConcat(t *testing.T) {
	a := characterNFA('a')
	b := characterNFA('b')
	c := characterNFA('c')

	a.concat(b)
	a.concat(c)

	bad := []string{
		"ab",
		"bc",
		"",
		"a",
		"b",
		"c",
	}

	if !a.process("abc") {
		t.Errorf("Failed to match %#v\n", "abc")
	}

	for _, s := range bad {
		if a.process(s) {
			t.Errorf("Matched %#v, but should not have\n", s)
		}
	}
}

func TestOr(t *testing.T) {
	a := characterNFA('a')
	b := characterNFA('b')
	a.or(b)

	good := []string{
		"a",
		"b",
	}

	bad := []string{
		"ab",
		"",
		"abc",
	}

	for _, s := range good {
		if !a.process(s) {
			t.Errorf("Failed to match %#v\n", s)
		}
	}

	for _, s := range bad {
		if a.process(s) {
			t.Errorf("Matched %#v, but should not have\n", s)
		}
	}
}

func TestOptional(t *testing.T) {
	a := characterNFA('a')
	a.makeOptional()

	good := []string{
		"a",
		"",
	}

	bad := []string{
		"b",
		"c",
		"ab",
		"bc",
	}

	for _, s := range good {
		if !a.process(s) {
			t.Errorf("Failed to match %#v\n", s)
		}
	}

	for _, s := range bad {
		if a.process(s) {
			t.Errorf("Matched %#v, but should not have\n", s)
		}
	}
}

func TestLoop(t *testing.T) {
	a := characterNFA('a')
	a.loop()

	good := []string{
		"",
		"a",
		"aa",
		"aaa",
		"aaaa",
	}

	bad := []string{
		"ab",
		"b",
		"bab",
		"c",
		"aaccb",
	}

	for _, s := range good {
		if !a.process(s) {
			t.Errorf("Failed to match %#v\n", s)
		}
	}

	for _, s := range bad {
		if a.process(s) {
			t.Errorf("Matched %#v, but should not have\n", s)
		}
	}
}
