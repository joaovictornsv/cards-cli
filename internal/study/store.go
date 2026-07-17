package study

import (
	"context"
	"errors"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/models"
)

var ErrDeckNotFound = errors.New("deck not found")

type Store interface {
	GetDeck(ctx context.Context, name string) (models.Deck, error)
	ListQueueCardIDs(ctx context.Context, deckName string) ([]int64, error)
	GetCard(ctx context.Context, deckName string, cardID int64) (models.Card, error)
	ReplaceQueue(ctx context.Context, deckID int64, cardIDs []int64) error
	SetReplaceEligible(ctx context.Context, deckName string, cardID int64) error
}

type DBStore struct {
	repo *db.Repository
}

func NewDBStore(repo *db.Repository) DBStore {
	return DBStore{repo: repo}
}

func (s DBStore) GetDeck(ctx context.Context, name string) (models.Deck, error) {
	deck, err := s.repo.GetDeckByName(ctx, name)
	if errors.Is(err, db.ErrDeckNotFound) {
		return models.Deck{}, ErrDeckNotFound
	}
	return deck, err
}

func (s DBStore) ListQueueCardIDs(ctx context.Context, deckName string) ([]int64, error) {
	ids, err := s.repo.ListQueueCardIDsByDeck(ctx, deckName)
	if errors.Is(err, db.ErrDeckNotFound) {
		return nil, ErrDeckNotFound
	}
	return ids, err
}

func (s DBStore) GetCard(ctx context.Context, deckName string, cardID int64) (models.Card, error) {
	return s.repo.GetCardByDeckAndID(ctx, deckName, cardID)
}

func (s DBStore) ReplaceQueue(ctx context.Context, deckID int64, cardIDs []int64) error {
	return s.repo.ReplaceDeckQueue(ctx, deckID, cardIDs)
}

func (s DBStore) SetReplaceEligible(ctx context.Context, deckName string, cardID int64) error {
	return s.repo.SetReplaceEligible(ctx, deckName, cardID, true)
}
