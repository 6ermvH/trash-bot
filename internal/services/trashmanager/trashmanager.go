package trashmanager

import (
	"context"
	"errors"
)

var (
	ErrChatIsNotInitialize = errors.New("")
)

type Repository interface {
	GetCurrent(ctx context.Context, chatID int64) (string, error)
	SetNext(ctx context.Context, chatID int64) error
	SetPrev(ctx context.Context, chatID int64) error
	SetEstablish(ctx context.Context, chatID int64, users []string) error
	Subscribe(ctx context.Context, chatID int64) error
	Unsubscribe(ctx context.Context, chatID int64) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Who(ctx context.Context, chatID int64) (string, error) {
	return s.repo.GetCurrent(ctx, chatID)
}

func (s *Service) Next(ctx context.Context, chatID int64) (string, error) {
	if err := s.repo.SetNext(ctx, chatID); err != nil {
		return "", err
	}

	return s.Who(ctx, chatID)
}

func (s *Service) Prev(ctx context.Context, chatID int64) (string, error) {
	if err := s.repo.SetPrev(ctx, chatID); err != nil {
		return "", err
	}

	return s.Who(ctx, chatID)
}

func (s *Service) SetEstablish(ctx context.Context, chatID int64, users []string) error {
	return s.repo.SetEstablish(ctx, chatID, users)
}

func (s *Service) Subscribe(ctx context.Context, chatID int64) error {
	return nil
}

func (s *Service) Unsubscribe(ctx context.Context, chatID int64) error {
	return nil
}
