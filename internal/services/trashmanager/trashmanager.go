package trashmanager

import (
	"context"
	"errors"
	"fmt"

	"github.com/6ermvH/trash-bot/internal/repository"
	"github.com/6ermvH/trash-bot/internal/repository/inmemory"
)

var (
	ErrTryToInitialize = errors.New("проведите инициализацию при помощи команды /set")
	ErrTryToAddUsers   = errors.New("добавьте пользователей в список через команду /set")
)

type Repository interface {
	GetChats(ctx context.Context) []inmemory.Chat

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

func (s *Service) Chats(ctx context.Context) []inmemory.Chat {
	return s.repo.GetChats(ctx)
}

func (s *Service) Who(ctx context.Context, chatID int64) (string, error) {
	username, err := s.repo.GetCurrent(ctx, chatID)

	switch {
	case err == nil:
		return username, nil
	case errors.Is(err, repository.ErrChatIsEmpty):
		return "", ErrTryToAddUsers
	case errors.Is(err, repository.ErrChatIsNotInitialize):
		return "", ErrTryToInitialize
	default:
		return "", fmt.Errorf("get who from repo: %w", err)
	}
}

func (s *Service) Next(ctx context.Context, chatID int64) (string, error) {
	err := s.repo.SetNext(ctx, chatID)

	switch {
	case err == nil:
		break
	case errors.Is(err, repository.ErrChatIsEmpty):
		return "", ErrTryToAddUsers
	case errors.Is(err, repository.ErrChatIsNotInitialize):
		return "", ErrTryToInitialize
	default:
		return "", fmt.Errorf("get next from repo: %w", err)
	}

	return s.Who(ctx, chatID)
}

func (s *Service) Prev(ctx context.Context, chatID int64) (string, error) {
	err := s.repo.SetPrev(ctx, chatID)

	switch {
	case err == nil:
		break
	case errors.Is(err, repository.ErrChatIsEmpty):
		return "", ErrTryToAddUsers
	case errors.Is(err, repository.ErrChatIsNotInitialize):
		return "", ErrTryToInitialize
	default:
		return "", fmt.Errorf("get prev from repo: %w", err)
	}

	return s.Who(ctx, chatID)
}

func (s *Service) SetEstablish(ctx context.Context, chatID int64, users []string) error {
	if err := s.repo.SetEstablish(ctx, chatID, users); err != nil {
		return fmt.Errorf("set establish from repo: %w", err)
	}

	return nil
}

func (s *Service) Subscribe(ctx context.Context, chatID int64) error {
	return nil
}

func (s *Service) Unsubscribe(ctx context.Context, chatID int64) error {
	return nil
}
