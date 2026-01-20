package inmemory

import (
	"testing"

	"github.com/6ermvH/trash-bot/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestGetCurrent(t *testing.T) {
	t.Parallel()

	t.Run("One user", func(t *testing.T) {
		t.Parallel()

		chats := []repository.Chat{
			{
				ID:      1,
				Users:   []string{"German"},
				Current: 0,
			},
		}
		repo := newTestRepo(t, chats)

		ctx := t.Context()

		username, err := repo.GetCurrent(ctx, chats[0].ID)
		require.NoError(t, err)

		require.Equal(t, username, chats[0].Users[0])
	})

	t.Run("More Users", func(t *testing.T) {
		t.Parallel()

		chats := []repository.Chat{
			{
				ID:      1,
				Users:   []string{"German", "Anthon", "Vitaly"},
				Current: 0,
			},
		}
		repo := newTestRepo(t, chats)

		ctx := t.Context()

		for ind := range 3 {
			repo.chats[1].Current = ind

			username, err := repo.GetCurrent(ctx, chats[0].ID)
			require.NoError(t, err)

			require.Equal(t, username, chats[0].Users[ind])
		}
	})
}

func TestSimpleWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("Next", func(t *testing.T) {
		t.Parallel()

		chats := []repository.Chat{
			{
				ID:      1,
				Users:   []string{"German", "Anthon", "Vitaly"},
				Current: 0,
			},
		}
		repo := newTestRepo(t, chats)

		ctx := t.Context()

		for ind := range 3 {
			username, err := repo.GetCurrent(ctx, chats[0].ID)
			require.NoError(t, err)

			require.Equal(t, username, chats[0].Users[ind%len(chats[0].Users)])

			require.NoError(t, repo.SetNext(ctx, chats[0].ID))
		}
	})

	t.Run("Prev", func(t *testing.T) {
		t.Parallel()

		chats := []repository.Chat{
			{
				ID:      1,
				Users:   []string{"German", "Anthon", "Vitaly"},
				Current: 0,
			},
		}
		repo := newTestRepo(t, chats)

		ctx := t.Context()

		for ind := range -3 {
			username, err := repo.GetCurrent(ctx, chats[0].ID)
			require.NoError(t, err)

			require.Equal(t, username, chats[0].
				Users[(ind+len(chats[0].Users))%len(chats[0].Users)])

			require.NoError(t, repo.SetPrev(ctx, chats[0].ID))
		}
	})
}

func TestSetEstablish(t *testing.T) {
	t.Parallel()

	t.Run("Simple", func(t *testing.T) {
		t.Parallel()

		chats := []repository.Chat{
			{
				ID:      1,
				Users:   []string{},
				Current: 0,
			},
		}
		repo := newTestRepo(t, chats)

		ctx := t.Context()

		users := []string{"German", "Anthon", "Vitaly"}
		require.NoError(t, repo.SetEstablish(ctx, chats[0].ID, users))

		require.Equal(t, repo.chats[1].Users, users)
	})
}

func TestSubscribe(t *testing.T) {
	t.Parallel()

	t.Run("Subscribe to existing chat", func(t *testing.T) {
		t.Parallel()

		chats := []repository.Chat{
			{
				ID:      1,
				Users:   []string{"German", "Anthon"},
				Current: 0,
			},
		}
		repo := newTestRepo(t, chats)
		ctx := t.Context()

		notifyTime := "09:00"
		err := repo.Subscribe(ctx, chats[0].ID, notifyTime)
		require.NoError(t, err)

		require.NotNil(t, repo.chats[1].NotifyTime)
		require.Equal(t, notifyTime, *repo.chats[1].NotifyTime)
	})

	t.Run("Subscribe to non-existing chat", func(t *testing.T) {
		t.Parallel()

		repo := New()
		ctx := t.Context()

		err := repo.Subscribe(ctx, 999, "09:00")
		require.ErrorIs(t, err, repository.ErrChatIsNotInitialize)
	})

	t.Run("Update subscription time", func(t *testing.T) {
		t.Parallel()

		oldTime := "08:00"
		chats := []repository.Chat{
			{
				ID:         1,
				Users:      []string{"German"},
				Current:    0,
				NotifyTime: &oldTime,
			},
		}
		repo := newTestRepo(t, chats)
		ctx := t.Context()

		newTime := "10:00"
		err := repo.Subscribe(ctx, chats[0].ID, newTime)
		require.NoError(t, err)

		require.Equal(t, newTime, *repo.chats[1].NotifyTime)
	})
}

func TestUnsubscribe(t *testing.T) {
	t.Parallel()

	t.Run("Unsubscribe from subscribed chat", func(t *testing.T) {
		t.Parallel()

		notifyTime := "09:00"
		chats := []repository.Chat{
			{
				ID:         1,
				Users:      []string{"German"},
				Current:    0,
				NotifyTime: &notifyTime,
			},
		}
		repo := newTestRepo(t, chats)
		ctx := t.Context()

		err := repo.Unsubscribe(ctx, chats[0].ID)
		require.NoError(t, err)

		require.Nil(t, repo.chats[1].NotifyTime)
	})

	t.Run("Unsubscribe from non-subscribed chat", func(t *testing.T) {
		t.Parallel()

		chats := []repository.Chat{
			{
				ID:      1,
				Users:   []string{"German"},
				Current: 0,
			},
		}
		repo := newTestRepo(t, chats)
		ctx := t.Context()

		err := repo.Unsubscribe(ctx, chats[0].ID)
		require.NoError(t, err)

		require.Nil(t, repo.chats[1].NotifyTime)
	})

	t.Run("Unsubscribe from non-existing chat", func(t *testing.T) {
		t.Parallel()

		repo := New()
		ctx := t.Context()

		// Не должно быть ошибки для несуществующего чата
		err := repo.Unsubscribe(ctx, 999)
		require.NoError(t, err)
	})
}

func TestGetSubscribedChats(t *testing.T) {
	t.Parallel()

	t.Run("Get subscribed chats when some are subscribed", func(t *testing.T) {
		t.Parallel()

		time1 := "09:00"
		time2 := "18:00"
		chats := []repository.Chat{
			{
				ID:         1,
				Users:      []string{"German"},
				Current:    0,
				NotifyTime: &time1,
			},
			{
				ID:      2,
				Users:   []string{"Anthon"},
				Current: 0,
			},
			{
				ID:         3,
				Users:      []string{"Vitaly"},
				Current:    0,
				NotifyTime: &time2,
			},
		}
		repo := newTestRepo(t, chats)
		ctx := t.Context()

		subscribed, err := repo.GetSubscribedChats(ctx)
		require.NoError(t, err)

		require.Len(t, subscribed, 2)

		// Проверяем, что вернулись только подписанные чаты
		ids := make(map[int64]bool)
		for _, chat := range subscribed {
			ids[chat.ID] = true

			require.NotNil(t, chat.NotifyTime)
		}

		require.True(t, ids[1])
		require.True(t, ids[3])
		require.False(t, ids[2])
	})

	t.Run("Get subscribed chats when none are subscribed", func(t *testing.T) {
		t.Parallel()

		chats := []repository.Chat{
			{
				ID:      1,
				Users:   []string{"German"},
				Current: 0,
			},
			{
				ID:      2,
				Users:   []string{"Anthon"},
				Current: 0,
			},
		}
		repo := newTestRepo(t, chats)
		ctx := t.Context()

		subscribed, err := repo.GetSubscribedChats(ctx)
		require.NoError(t, err)

		require.Empty(t, subscribed)
	})

	t.Run("Get subscribed chats from empty repo", func(t *testing.T) {
		t.Parallel()

		repo := New()
		ctx := t.Context()

		subscribed, err := repo.GetSubscribedChats(ctx)
		require.NoError(t, err)

		require.Empty(t, subscribed)
	})
}

func newTestRepo(t *testing.T, chats []repository.Chat) *RepoInMem {
	t.Helper()

	repo := New()

	for _, chat := range chats {
		chatCopy := chat
		repo.chats[chat.ID] = &chatCopy
	}

	return repo
}
