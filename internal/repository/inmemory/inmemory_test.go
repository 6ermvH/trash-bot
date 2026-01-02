package inmemory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCurrent(t *testing.T) {
	t.Parallel()

	t.Run("One user", func(t *testing.T) {
		t.Parallel()

		chats := []Chat{
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

		chats := []Chat{
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

		chats := []Chat{
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

		chats := []Chat{
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

		chats := []Chat{
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

func newTestRepo(t *testing.T, chats []Chat) *RepoInMem {
	t.Helper()

	repo := New()
	for _, chat := range chats {
		repo.chats[chat.ID] = &chat
	}

	return repo
}
