package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Club})
	fmt.Println(Card{Rank: Three, Suit: Diamond})
	fmt.Println(Card{Rank: Four, Suit: Heart})
	fmt.Println(Card{Rank: King, Suit: Spade})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Hearts
	// Two of Clubs
	// Three of Diamonds
	// Four of Hearts
	// King of Spades
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	// 13 ranks and 4 suits
	if len(cards) != 13*4 {
		t.Error("Wrong number of cards in the new deck!")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	exp := Card{Rank: Ace, Suit: Spade}
	if cards[0] != exp {
		t.Error("First card is not the expected Ace of Spades. Received:", cards[0])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	exp := Card{Rank: Ace, Suit: Spade}
	if cards[0] != exp {
		t.Error("First card is not the expected Ace of Spades. Received:", cards[0])
	}
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(4))
	count := 0
	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}
	if count != 4 {
		t.Error("Expected 4 Jokers, received:", count)
	}
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}
	cards := New(Filter(filter))
	for _, c := range cards {
		if c.Rank == Two || c.Rank == Three {
			t.Error("Expected to filter cards with Rank Two and Three")
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))

	if len(cards) != 13*4*3 {
		t.Errorf("Expected %d cards. Received %d cards.", 13*4*3, len(cards))
	}
}
