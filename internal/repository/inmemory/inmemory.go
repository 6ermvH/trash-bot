package inmemory

import (
	"context"

	"github.com/6ermvH/trash-bot/internal/model"
	"github.com/6ermvH/trash-bot/internal/repository"
)

type RepoInMem struct {
	chats map[int64]*model.Chat
}

func New() *RepoInMem {
	return &RepoInMem{chats: make(map[int64]*model.Chat)}
}

func (r *RepoInMem) GetCurrent(ctx context.Context, chatID int64) (string, error) {
	chat, ok := r.chats[chatID]
	if !ok {
		return "", repository.ErrChatIsNotInitialize
	}

	if len(chat.Users) == 0 {
		return "", repository.ErrChatIsEmpty
	}

	return chat.Users[chat.Current], nil
}

func (r *RepoInMem) SetNext(ctx context.Context, chatID int64) error {
	chat, ok := r.chats[chatID]
	if !ok {
		return repository.ErrChatIsNotInitialize
	}

	if len(chat.Users) == 0 {
		return repository.ErrChatIsEmpty
	}

	chat.Current = (chat.Current + 1) % len(chat.Users)

	return nil
}

func (r *RepoInMem) SetPrev(ctx context.Context, chatID int64) error {
	chat, ok := r.chats[chatID]
	if !ok {
		return repository.ErrChatIsNotInitialize
	}

	if len(chat.Users) == 0 {
		return repository.ErrChatIsEmpty
	}

	chat.Current = (len(chat.Users) + chat.Current - 1) % len(chat.Users)

	return nil
}

func (r *RepoInMem) SetEstablish(ctx context.Context, chatID int64, users []string) error {
	chat, ok := r.chats[chatID]
	if !ok {
		chat = &model.Chat{
			ID:      chatID,
			Current: 0,
			Users:   make([]string, 0, len(users)),
		}
		r.chats[chatID] = chat
	}

	chat.Users = append(chat.Users, users...)

	return nil

}

func (r *RepoInMem) Subscribe(ctx context.Context, chatID int64) error {
	return nil
}

func (r *RepoInMem) Unsubscribe(ctx context.Context, chatID int64) error {
	return nil
}
