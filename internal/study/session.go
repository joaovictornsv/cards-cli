package study

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/joaovictornsv/cards-cli/internal/queue"
)

type Options struct {
	BatchSize int
	QueueOpts queue.Options
}

type Session struct {
	DeckName string
	Out      io.Writer
	Store    Store
	Input    Input
	Opts     Options
}

func (s *Session) Run(ctx context.Context) error {
	deck, err := s.Store.GetDeck(ctx, s.DeckName)
	if err != nil {
		return err
	}

	ids, err := s.Store.ListQueueCardIDs(ctx, s.DeckName)
	if err != nil {
		return err
	}

	limit := s.Opts.BatchSize
	if limit > len(ids) {
		limit = len(ids)
	}
	if limit == 0 {
		return fmt.Errorf("no cards in deck %q", s.DeckName)
	}

	batch, tail := queue.Pull(ids, limit)
	pending := append([]int64(nil), batch...)

	fmt.Fprintf(s.Out, "\nSession: %s (batch %d/%d, %d cards in deck)\n\n",
		s.DeckName, limit, limit, len(ids))

	for i := 0; i < len(batch); i++ {
		cardID := pending[0]
		card, err := s.Store.GetCard(ctx, s.DeckName, cardID)
		if err != nil {
			return err
		}

		fmt.Fprintf(s.Out, "[%d/%d] %s\n", i+1, limit, card.Front)
		if i == 0 {
			fmt.Fprintln(s.Out, "      (space/enter to reveal, 1/2/3 or arrows to grade, q to quit)")
		}
		fmt.Fprintln(s.Out)

		if err := s.Input.WaitReveal(); err != nil {
			if errors.Is(err, ErrQuit) {
				return s.persist(ctx, deck.ID, pending, tail)
			}
			return err
		}

		fmt.Fprintf(s.Out, "      %s\n\n", card.Back)
		fmt.Fprintln(s.Out, "      [1] again   [2] hard   [3] easy")

		grade, err := s.Input.ReadGrade()
		if err != nil {
			if errors.Is(err, ErrQuit) {
				return s.persist(ctx, deck.ID, pending, tail)
			}
			return err
		}

		pending = pending[1:]
		tail, err = queue.ReinsertAfterGrade(tail, cardID, grade, s.Opts.QueueOpts)
		if err != nil {
			return err
		}
		if err := s.persist(ctx, deck.ID, pending, tail); err != nil {
			return err
		}
	}

	fmt.Fprintf(s.Out, "\nSession complete.\n\n")
	return nil
}

func (s *Session) persist(ctx context.Context, deckID int64, pending, tail []int64) error {
	full := append(append([]int64(nil), pending...), tail...)
	return s.Store.ReplaceQueue(ctx, deckID, full)
}
