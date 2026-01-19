package inmemory

import (
	"context"
	"sync"

	"github.com/6ermvH/trash-bot/internal/repository"
)

type RepoInMem struct {
	chats map[int64]*repository.Chat
	mu    sync.Mutex
}

func New() *RepoInMem {
	return &RepoInMem{chats: make(map[int64]*repository.Chat)}
}

func (r *RepoInMem) GetChats(ctx context.Context) ([]repository.Chat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]repository.Chat, 0)
	for _, chat := range r.chats {
		result = append(result, *chat)
	}

	return result, nil
}

func (r *RepoInMem) GetChat(ctx context.Context, chatID int64) (*repository.Chat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	chat, ok := r.chats[chatID]
	if !ok {
		return nil, repository.ErrChatIsNotInitialize
	}

	chatCopy := *chat
	return &chatCopy, nil
}

func (r *RepoInMem) GetCurrent(ctx context.Context, chatID int64) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

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
	r.mu.Lock()
	defer r.mu.Unlock()

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
	r.mu.Lock()
	defer r.mu.Unlock()

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
	r.mu.Lock()
	defer r.mu.Unlock()

	chat, ok := r.chats[chatID]
	if !ok {
		chat = &repository.Chat{
			ID:      chatID,
			Current: 0,
		}
		r.chats[chatID] = chat
	}

	chat.Users = users
	chat.Current = 0

	return nil
}

func (r *RepoInMem) Subscribe(ctx context.Context, chatID int64) error {
	return nil
}

func (r *RepoInMem) Unsubscribe(ctx context.Context, chatID int64) error {
	return nil
}
