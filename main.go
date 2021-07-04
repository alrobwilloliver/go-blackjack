package main

import (
	"fmt"
	"strconv"
	"strings"

	"blackjack/deck"
)

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func (h Hand) MinScore() int {
	score := 0

	for _, c := range h {

		score += min(int(c.Rank), 10)
	}

	return score
}

func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			// ace worth 1 and we add 10 to make it worth 11 if possible higher score (less than 22 total)
			return minScore + 10
		}
	}
	return minScore
}

func (h Hand) Blackjack() bool {
	return h.Score() == 21 && len(h) == 2
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Bet(gs GameState) GameState {
	ret := clone(gs)
	fmt.Println("====BETTING====")
	fmt.Println("Player Balance", ret.PlayerBalance)
	fmt.Println("Dealer Balance", ret.DealerBalance)
	fmt.Println("How much will you bet?")
	var betInput string
	fmt.Scanf("%s\n", &betInput)
	for ret.PlayerBet == 0 {
		i, err := strconv.Atoi(betInput)
		if err != nil {
			fmt.Printf("Can't convert %s to integer\n", i)
			fmt.Scanf("%s\n", &betInput)
		} else if i < 1 {
			fmt.Println("You can't bet less than 1")
			fmt.Scanf("%s\n", &betInput)
		} else if i > ret.PlayerBalance {
			fmt.Println("You can't bet more than your available balance!")
			fmt.Scanf("%s\n", &betInput)
		} else {
			ret.PlayerBalance = ret.PlayerBalance - i
			ret.DealerBalance = ret.DealerBalance - i
			ret.PlayerBet = i
			ret.Pot = ret.PlayerBet * 2
		}
	}
	fmt.Println("===REMAINING BALANCES===")
	fmt.Println("Dealer balance", ret.DealerBalance)
	fmt.Println("Player balance", ret.PlayerBalance)
	fmt.Println("Player bet:", ret.PlayerBet)
	fmt.Println("Pot:", ret.Pot)
	return ret
}

func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)
		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)
	}
	ret.State = StatePlayerTurn
	return ret
}

func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()
	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(ret)
	}
	return ret
}

func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++
	return ret
}

func EndHand(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := ret.Player.Score(), ret.Dealer.Score()
	fmt.Println("===FINAL SCORES===")
	fmt.Println("Player:", ret.Player)
	fmt.Println("Score:", pScore)
	fmt.Println("Dealer:", ret.Dealer)
	fmt.Println("Score:", dScore)
	fmt.Println("Balance", ret.PlayerBalance, ret.DealerBalance)
	switch {
	case ret.Dealer.Blackjack():
		fmt.Println("Dealer drew a blackjack!")
		ret.DealerBalance += ret.Pot
	case ret.Player.Blackjack():
		fmt.Println("Player drew a blackjack!")
		ret.PlayerBalance += ret.Pot
	case pScore > 21:
		fmt.Println("Player busted!")
		ret.DealerBalance = ret.DealerBalance + ret.Pot
	case dScore > 21:
		fmt.Println("Dealer busted!")
		ret.PlayerBalance = ret.PlayerBalance + ret.Pot
	case dScore > pScore:
		fmt.Println("You lose!")
		ret.DealerBalance = ret.DealerBalance + ret.Pot
	case pScore > dScore:
		fmt.Println("You win!")
		ret.PlayerBalance = ret.PlayerBalance + ret.Pot
	case pScore == dScore:
		fmt.Println("Draw")
		ret.PlayerBalance = ret.PlayerBalance + ret.Pot/2
		ret.DealerBalance = ret.PlayerBalance + ret.Pot/2
	}
	ret.Player = nil
	ret.Dealer = nil
	ret.State = StatePlayerTurn
	fmt.Println("===============")
	return ret
}

func main() {
	var gs GameState
	gs = Shuffle(gs)
	gs = Deal(gs)
	gs.DealerBalance = 100
	gs.PlayerBalance = 100

	for gs.DealerBalance != 0 && gs.PlayerBalance != 0 {
		gs = Bet(gs)
		if gs.Dealer.Blackjack() || gs.Player.Blackjack() {
			gs = EndHand(gs)
			if gs.DealerBalance == 0 || gs.PlayerBalance == 0 {
				EndGame(gs)
			}
		}
		gs = Deal(gs)
		var input string
		for gs.State == StatePlayerTurn {
			fmt.Println("Player: ", gs.Player)
			fmt.Println("Dealer: ", gs.Dealer.DealerString())
			fmt.Println("What will you do? (h)it, (s)tand")
			fmt.Scanf("%s\n", &input)
			switch input {
			case "h":
				gs = Hit(gs)
			case "s":
				gs = Stand(gs)
			default:
				fmt.Println("Invalid input:", input)
			}
		}

		for gs.State == StateDealerTurn {
			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
				gs = Hit(gs)
			} else {
				gs = Stand(gs)
			}
		}
		gs = EndHand(gs)
		EndGame(gs)
	}
}

func EndGame(gs GameState) {
	if gs.PlayerBalance == 0 {
		fmt.Println("===WINNER===")
		fmt.Println("The Player has no remaining balance! Dealer Wins!")
		fmt.Println("============")
		return
	} else if gs.DealerBalance == 0 {
		fmt.Println("===WINNER===")
		fmt.Println("The Dealer has no remaining balance! Player Wins!")
		fmt.Println("============")
		return
	}
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

type State int8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck          []deck.Card
	State         State
	Player        Hand
	PlayerBet     int
	PlayerBalance int
	Dealer        Hand
	DealerBet     int
	DealerBalance int
	Pot           int
}

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("It is currently not any player's turn")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:          make([]deck.Card, len(gs.Deck)),
		State:         gs.State,
		PlayerBalance: gs.PlayerBalance,
		DealerBalance: gs.DealerBalance,
		Pot:           gs.Pot,
		Player:        make([]deck.Card, len(gs.Player)),
		Dealer:        make([]deck.Card, len(gs.Dealer)),
	}
	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)
	return ret
}
