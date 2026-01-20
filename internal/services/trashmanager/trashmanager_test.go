package trashmanager

import (
	"context"
	"errors"
	"testing"

	"github.com/6ermvH/trash-bot/internal/repository"
	"github.com/stretchr/testify/require"
)

var errDatabaseConnection = errors.New("database connection failed")

type mockRepo struct {
	chats map[int64]*repository.Chat
}

func newMockRepo() *mockRepo {
	return &mockRepo{chats: make(map[int64]*repository.Chat)}
}

func (m *mockRepo) GetChats(ctx context.Context) ([]repository.Chat, error) {
	result := make([]repository.Chat, 0, len(m.chats))
	for _, chat := range m.chats {
		result = append(result, *chat)
	}

	return result, nil
}

func (m *mockRepo) GetChat(ctx context.Context, chatID int64) (*repository.Chat, error) {
	chat, ok := m.chats[chatID]
	if !ok {
		return nil, repository.ErrChatIsNotInitialize
	}

	return chat, nil
}

func (m *mockRepo) GetSubscribedChats(ctx context.Context) ([]repository.Chat, error) {
	result := make([]repository.Chat, 0)

	for _, chat := range m.chats {
		if chat.NotifyTime != nil {
			result = append(result, *chat)
		}
	}

	return result, nil
}

func (m *mockRepo) GetCurrent(ctx context.Context, chatID int64) (string, error) {
	chat, ok := m.chats[chatID]
	if !ok {
		return "", repository.ErrChatIsNotInitialize
	}

	if len(chat.Users) == 0 {
		return "", repository.ErrChatIsEmpty
	}

	return chat.Users[chat.Current], nil
}

func (m *mockRepo) SetNext(ctx context.Context, chatID int64) error {
	chat, ok := m.chats[chatID]
	if !ok {
		return repository.ErrChatIsNotInitialize
	}

	if len(chat.Users) == 0 {
		return repository.ErrChatIsEmpty
	}

	chat.Current = (chat.Current + 1) % len(chat.Users)

	return nil
}

func (m *mockRepo) SetPrev(ctx context.Context, chatID int64) error {
	chat, ok := m.chats[chatID]
	if !ok {
		return repository.ErrChatIsNotInitialize
	}

	if len(chat.Users) == 0 {
		return repository.ErrChatIsEmpty
	}

	chat.Current = (len(chat.Users) + chat.Current - 1) % len(chat.Users)

	return nil
}

func (m *mockRepo) SetEstablish(ctx context.Context, chatID int64, users []string) error {
	chat, ok := m.chats[chatID]
	if !ok {
		chat = &repository.Chat{
			ID:      chatID,
			Current: 0,
		}
		m.chats[chatID] = chat
	}

	chat.Users = users
	chat.Current = 0

	return nil
}

func (m *mockRepo) Subscribe(ctx context.Context, chatID int64, notifyTime string) error {
	chat, ok := m.chats[chatID]
	if !ok {
		return repository.ErrChatIsNotInitialize
	}

	chat.NotifyTime = &notifyTime

	return nil
}

func (m *mockRepo) Unsubscribe(ctx context.Context, chatID int64) error {
	chat, ok := m.chats[chatID]
	if !ok {
		return nil
	}

	chat.NotifyTime = nil

	return nil
}

func TestService_Subscribe(t *testing.T) {
	t.Parallel()

	t.Run("Subscribe to existing chat", func(t *testing.T) {
		t.Parallel()

		repo := newMockRepo()
		repo.chats[1] = &repository.Chat{
			ID:      1,
			Users:   []string{"German", "Anthon"},
			Current: 0,
		}

		service := New(repo)
		ctx := t.Context()

		err := service.Subscribe(ctx, 1, "09:00")
		require.NoError(t, err)

		require.NotNil(t, repo.chats[1].NotifyTime)
		require.Equal(t, "09:00", *repo.chats[1].NotifyTime)
	})

	t.Run("Subscribe to non-existing chat returns error", func(t *testing.T) {
		t.Parallel()

		repo := newMockRepo()
		service := New(repo)
		ctx := t.Context()

		err := service.Subscribe(ctx, 999, "09:00")
		require.ErrorIs(t, err, ErrTryToInitialize)
	})

	t.Run("Update subscription time", func(t *testing.T) {
		t.Parallel()

		oldTime := "08:00"
		repo := newMockRepo()
		repo.chats[1] = &repository.Chat{
			ID:         1,
			Users:      []string{"German"},
			Current:    0,
			NotifyTime: &oldTime,
		}

		service := New(repo)
		ctx := t.Context()

		err := service.Subscribe(ctx, 1, "10:00")
		require.NoError(t, err)

		require.Equal(t, "10:00", *repo.chats[1].NotifyTime)
	})
}

func TestService_Unsubscribe(t *testing.T) {
	t.Parallel()

	t.Run("Unsubscribe from subscribed chat", func(t *testing.T) {
		t.Parallel()

		notifyTime := "09:00"
		repo := newMockRepo()
		repo.chats[1] = &repository.Chat{
			ID:         1,
			Users:      []string{"German"},
			Current:    0,
			NotifyTime: &notifyTime,
		}

		service := New(repo)
		ctx := t.Context()

		err := service.Unsubscribe(ctx, 1)
		require.NoError(t, err)

		require.Nil(t, repo.chats[1].NotifyTime)
	})

	t.Run("Unsubscribe from non-existing chat succeeds", func(t *testing.T) {
		t.Parallel()

		repo := newMockRepo()
		service := New(repo)
		ctx := t.Context()

		// Не должно быть ошибки
		err := service.Unsubscribe(ctx, 999)
		require.NoError(t, err)
	})
}

func TestService_GetSubscribedChats(t *testing.T) {
	t.Parallel()

	t.Run("Get subscribed chats", func(t *testing.T) {
		t.Parallel()

		time1 := "09:00"
		time2 := "18:00"

		repo := newMockRepo()
		repo.chats[1] = &repository.Chat{
			ID:         1,
			Users:      []string{"German"},
			Current:    0,
			NotifyTime: &time1,
		}
		repo.chats[2] = &repository.Chat{
			ID:      2,
			Users:   []string{"Anthon"},
			Current: 0,
		}
		repo.chats[3] = &repository.Chat{
			ID:         3,
			Users:      []string{"Vitaly"},
			Current:    0,
			NotifyTime: &time2,
		}

		service := New(repo)
		ctx := t.Context()

		subscribed, err := service.GetSubscribedChats(ctx)
		require.NoError(t, err)

		require.Len(t, subscribed, 2)
	})

	t.Run("Get subscribed chats when none subscribed", func(t *testing.T) {
		t.Parallel()

		repo := newMockRepo()
		repo.chats[1] = &repository.Chat{
			ID:      1,
			Users:   []string{"German"},
			Current: 0,
		}

		service := New(repo)
		ctx := t.Context()

		subscribed, err := service.GetSubscribedChats(ctx)
		require.NoError(t, err)

		require.Empty(t, subscribed)
	})
}

func TestService_SubscribeUnsubscribeWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("Full workflow: establish -> subscribe -> who -> unsubscribe", func(t *testing.T) {
		t.Parallel()

		repo := newMockRepo()
		service := New(repo)
		ctx := t.Context()

		// Устанавливаем пользователей
		err := service.SetEstablish(ctx, 1, []string{"German", "Anthon", "Vitaly"})
		require.NoError(t, err)

		// Проверяем кто выносит
		who, err := service.Who(ctx, 1)
		require.NoError(t, err)
		require.Equal(t, "German", who)

		// Подписываемся на уведомления
		err = service.Subscribe(ctx, 1, "09:00")
		require.NoError(t, err)

		// Проверяем что подписка активна
		subscribed, err := service.GetSubscribedChats(ctx)
		require.NoError(t, err)
		require.Len(t, subscribed, 1)
		require.Equal(t, int64(1), subscribed[0].ID)
		require.Equal(t, "09:00", *subscribed[0].NotifyTime)

		// Отписываемся
		err = service.Unsubscribe(ctx, 1)
		require.NoError(t, err)

		// Проверяем что подписок больше нет
		subscribed, err = service.GetSubscribedChats(ctx)
		require.NoError(t, err)
		require.Empty(t, subscribed)
	})
}

// Тест на проверку ошибки репозитория.
type errorRepo struct {
	mockRepo

	subscribeErr error
}

func (e *errorRepo) Subscribe(ctx context.Context, chatID int64, notifyTime string) error {
	if e.subscribeErr != nil {
		return e.subscribeErr
	}

	return e.mockRepo.Subscribe(ctx, chatID, notifyTime)
}

func TestService_Subscribe_RepoError(t *testing.T) {
	t.Parallel()

	t.Run("Repository error is wrapped", func(t *testing.T) {
		t.Parallel()

		repo := &errorRepo{
			mockRepo:     mockRepo{chats: make(map[int64]*repository.Chat)},
			subscribeErr: errDatabaseConnection,
		}
		repo.chats[1] = &repository.Chat{
			ID:      1,
			Users:   []string{"German"},
			Current: 0,
		}

		service := New(repo)
		ctx := t.Context()

		err := service.Subscribe(ctx, 1, "09:00")
		require.Error(t, err)
		require.ErrorIs(t, err, errDatabaseConnection)
	})
}
