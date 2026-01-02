package trashmanager

import (
	"context"
	"errors"
	"fmt"
)

var ErrChatIsNotInitialize = errors.New("")

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
	username, err := s.repo.GetCurrent(ctx, chatID)
	switch err {
	case nil:
		return username, nil
	default:
		return "", fmt.Errorf("get Who From chat: %w", err)
	}
}

func (s *Service) Next(ctx context.Context, chatID int64) (string, error) {
	err := s.repo.SetNext(ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("get next: %w", err)
	}

	return s.Who(ctx, chatID)
}

func (s *Service) Prev(ctx context.Context, chatID int64) (string, error) {
	err := s.repo.SetPrev(ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("get prev: %w", err)
	}

	return s.Who(ctx, chatID)
}

func (s *Service) SetEstablish(ctx context.Context, chatID int64, users []string) error {
	if err := s.repo.SetEstablish(ctx, chatID, users); err != nil {
		return fmt.Errorf("set establish: %w", err)
	}

	return nil
}

func (s *Service) Subscribe(ctx context.Context, chatID int64) error {
	return nil
}

func (s *Service) Unsubscribe(ctx context.Context, chatID int64) error {
	return nil
}
