package inmemory

import (
	"context"
	"sync"

	"github.com/6ermvH/trash-bot/internal/repository"
)

type Chat struct {
	ID      int64    `json:"id"`
	Current int      `json:"currentUser"`
	Users   []string `json:"activeUsers"`
}

type RepoInMem struct {
	chats map[int64]*Chat
	mu    sync.Mutex
}

func New() *RepoInMem {
	return &RepoInMem{chats: make(map[int64]*Chat)}
}

func (r *RepoInMem) GetChats(ctx context.Context) []Chat {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]Chat, 0)
	for _, chat := range r.chats {
		result = append(result, *chat)
	}

	return result
}

func (r *RepoInMem) GetChat(ctx context.Context, chatID int64) (*Chat, error) {
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
		chat = &Chat{
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
