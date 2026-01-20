package trashmanager

import (
	"context"
	"errors"
	"fmt"

	"github.com/6ermvH/trash-bot/internal/repository"
)

var (
	ErrTryToInitialize = errors.New("проведите инициализацию при помощи команды /set")
	ErrTryToAddUsers   = errors.New("добавьте пользователей в список через команду /set")
)

type Repository interface {
	GetChats(ctx context.Context) ([]repository.Chat, error)
	GetChat(ctx context.Context, chatID int64) (*repository.Chat, error)
	GetSubscribedChats(ctx context.Context) ([]repository.Chat, error)

	GetCurrent(ctx context.Context, chatID int64) (string, error)
	SetNext(ctx context.Context, chatID int64) error
	SetPrev(ctx context.Context, chatID int64) error
	SetEstablish(ctx context.Context, chatID int64, users []string) error
	Subscribe(ctx context.Context, chatID int64, notifyTime string) error
	Unsubscribe(ctx context.Context, chatID int64) error
}

type Stats struct {
	TotalChats      int     `json:"totalChats"`
	TotalUsers      int     `json:"totalUsers"`
	AvgUsersPerChat float64 `json:"avgUsersPerChat"`
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Chats(ctx context.Context) ([]repository.Chat, error) {
	chats, err := s.repo.GetChats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chats from repo: %w", err)
	}

	return chats, nil
}

func (s *Service) Chat(ctx context.Context, chatID int64) (*repository.Chat, error) {
	chat, err := s.repo.GetChat(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat from repo: %w", err)
	}

	return chat, nil
}

func (s *Service) Stats(ctx context.Context) (Stats, error) {
	chats, err := s.repo.GetChats(ctx)
	if err != nil {
		return Stats{}, fmt.Errorf("get chats for stats: %w", err)
	}

	totalChats := len(chats)

	totalUsers := 0
	for _, chat := range chats {
		totalUsers += len(chat.Users)
	}

	var avgUsers float64
	if totalChats > 0 {
		avgUsers = float64(totalUsers) / float64(totalChats)
	}

	return Stats{
		TotalChats:      totalChats,
		TotalUsers:      totalUsers,
		AvgUsersPerChat: avgUsers,
	}, nil
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

func (s *Service) Subscribe(ctx context.Context, chatID int64, notifyTime string) error {
	if _, err := s.repo.GetChat(ctx, chatID); err != nil {
		if errors.Is(err, repository.ErrChatIsNotInitialize) {
			return ErrTryToInitialize
		}

		return fmt.Errorf("get chat for subscribe: %w", err)
	}

	if err := s.repo.Subscribe(ctx, chatID, notifyTime); err != nil {
		return fmt.Errorf("subscribe in repo: %w", err)
	}

	return nil
}

func (s *Service) Unsubscribe(ctx context.Context, chatID int64) error {
	if err := s.repo.Unsubscribe(ctx, chatID); err != nil {
		return fmt.Errorf("unsubscribe in repo: %w", err)
	}

	return nil
}

func (s *Service) GetSubscribedChats(ctx context.Context) ([]repository.Chat, error) {
	chats, err := s.repo.GetSubscribedChats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get subscribed chats from repo: %w", err)
	}

	return chats, nil
}
