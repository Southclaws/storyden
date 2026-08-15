package searchscore

import "strings"

type Field struct {
	Text           string
	ExactMatch     int
	SubstringMatch int
	TermMatch      int
}

type Query struct {
	phrase string
	terms  []string
}

func New(query string) Query {
	phrase := strings.ToLower(strings.TrimSpace(query))
	return Query{
		phrase: phrase,
		terms:  strings.Fields(phrase),
	}
}

func (q Query) Empty() bool {
	return q.phrase == ""
}

func (q Query) Score(fields ...Field) int {
	if q.Empty() {
		return 0
	}

	score := 0
	for _, field := range fields {
		text := strings.ToLower(field.Text)
		if field.ExactMatch > 0 && text == q.phrase {
			score += field.ExactMatch
		} else if field.SubstringMatch > 0 && strings.Contains(text, q.phrase) {
			score += field.SubstringMatch
		}

		if field.TermMatch == 0 {
			continue
		}
		for _, term := range q.terms {
			if strings.Contains(text, term) {
				score += field.TermMatch
			}
		}
	}
	return score
}
