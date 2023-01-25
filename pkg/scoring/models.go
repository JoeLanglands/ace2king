package scoring

import "errors"

type Player struct {
	Name   string         `json:"name"`
	Score  int            `json:"score"`
	Scores map[string]int `json:"scores"`
}

var Cards = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

func NewPlayer(name string) Player {
	return Player{
		Name:   name,
		Score:  0,
		Scores: make(map[string]int, len(Cards)),
	}
}
func (p *Player) AddScore(score int, card string) {
	p.Scores[card] = score
	var scoreSum int
	for _, s := range p.Scores {
		scoreSum += s
	}
	p.Score = scoreSum
}

// CountWins returns the number of rounds won by the player up to cardIdx.
func (p *Player) CountWins(cardIdx int) (int, error) {
	var winSum int
	if cardIdx >= len(Cards) || cardIdx < 0 {
		return 0, errors.New("cardIdx out of range")
	}
	for _, card := range Cards[:cardIdx+1] {
		if p.Scores[card] == 0 {
			winSum++
		}
	}
	return winSum, nil
}
