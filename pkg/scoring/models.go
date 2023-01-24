package scoring

type Player struct {
	Name   string
	Score  int
	Wins   int
	Scores map[string]int
}

var Cards = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

func NewPlayer(name string) Player {
	return Player{
		Name:   name,
		Score:  0,
		Wins:   0,
		Scores: make(map[string]int, len(Cards)),
	}
}

func (p *Player) AddScore(score int, card string) {
	p.Scores[card] = score
	var scoreSum int
	var winSum int
	for _, s := range p.Scores {
		scoreSum += s
		if s <= 0 {
			winSum++
		}
	}
	p.Score = scoreSum
	p.Wins = winSum
}
