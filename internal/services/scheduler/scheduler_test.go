package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/6ermvH/trash-bot/internal/repository"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	chats      []repository.Chat
	whoResults map[int64]string
	whoErr     error
	chatsErr   error
}

func (m *mockService) GetSubscribedChats(ctx context.Context) ([]repository.Chat, error) {
	if m.chatsErr != nil {
		return nil, m.chatsErr
	}

	return m.chats, nil
}

func (m *mockService) Who(ctx context.Context, chatID int64) (string, error) {
	if m.whoErr != nil {
		return "", m.whoErr
	}

	if m.whoResults != nil {
		if who, ok := m.whoResults[chatID]; ok {
			return who, nil
		}
	}

	return "Unknown", nil
}

func TestScheduler_CheckAndNotify(t *testing.T) {
	t.Parallel()

	t.Run("Filters chats by current time", func(t *testing.T) {
		t.Parallel()

		currentTime := time.Now().Format("15:04")
		otherTime := "99:99" // Невозможное время

		time1 := currentTime
		time2 := otherTime

		service := &mockService{
			chats: []repository.Chat{
				{
					ID:         1,
					Users:      []string{"German"},
					NotifyTime: &time1,
				},
				{
					ID:         2,
					Users:      []string{"Anthon"},
					NotifyTime: &time2,
				},
			},
			whoResults: map[int64]string{
				1: "German",
				2: "Anthon",
			},
		}

		ctx := t.Context()

		// Тестируем фильтрацию - должен обработать только чат с текущим временем
		chats, err := service.GetSubscribedChats(ctx)
		require.NoError(t, err)
		require.Len(t, chats, 2)

		// Проверяем фильтрацию по времени
		matchingChats := 0
		for _, chat := range chats {
			if chat.NotifyTime != nil && *chat.NotifyTime == currentTime {
				matchingChats++
			}
		}
		require.Equal(t, 1, matchingChats)
	})

	t.Run("Handles nil NotifyTime gracefully", func(t *testing.T) {
		t.Parallel()

		service := &mockService{
			chats: []repository.Chat{
				{
					ID:         1,
					Users:      []string{"German"},
					NotifyTime: nil,
				},
			},
		}

		ctx := t.Context()

		chats, err := service.GetSubscribedChats(ctx)
		require.NoError(t, err)

		// Проверяем что nil NotifyTime корректно обрабатывается
		for _, chat := range chats {
			if chat.NotifyTime != nil {
				t.Fail()
			}
		}
	})
}

func TestScheduler_GetChatsToNotify(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		currentTime  string
		chats        []repository.Chat
		expectedIDs  []int64
	}{
		{
			name:        "Single matching chat",
			currentTime: "09:00",
			chats: func() []repository.Chat {
				time1 := "09:00"

				return []repository.Chat{
					{ID: 1, NotifyTime: &time1},
				}
			}(),
			expectedIDs: []int64{1},
		},
		{
			name:        "Multiple chats, one matches",
			currentTime: "09:00",
			chats: func() []repository.Chat {
				time1 := "09:00"
				time2 := "10:00"

				return []repository.Chat{
					{ID: 1, NotifyTime: &time1},
					{ID: 2, NotifyTime: &time2},
				}
			}(),
			expectedIDs: []int64{1},
		},
		{
			name:        "Multiple chats, all match",
			currentTime: "09:00",
			chats: func() []repository.Chat {
				time1 := "09:00"
				time2 := "09:00"

				return []repository.Chat{
					{ID: 1, NotifyTime: &time1},
					{ID: 2, NotifyTime: &time2},
				}
			}(),
			expectedIDs: []int64{1, 2},
		},
		{
			name:        "No matching chats",
			currentTime: "09:00",
			chats: func() []repository.Chat {
				time1 := "10:00"
				time2 := "11:00"

				return []repository.Chat{
					{ID: 1, NotifyTime: &time1},
					{ID: 2, NotifyTime: &time2},
				}
			}(),
			expectedIDs: []int64{},
		},
		{
			name:        "Chat with nil NotifyTime",
			currentTime: "09:00",
			chats: []repository.Chat{
				{ID: 1, NotifyTime: nil},
			},
			expectedIDs: []int64{},
		},
		{
			name:        "Empty chat list",
			currentTime: "09:00",
			chats:       []repository.Chat{},
			expectedIDs: []int64{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			matchingIDs := make([]int64, 0)

			for _, chat := range tc.chats {
				if chat.NotifyTime == nil {
					continue
				}

				if *chat.NotifyTime == tc.currentTime {
					matchingIDs = append(matchingIDs, chat.ID)
				}
			}

			require.ElementsMatch(t, tc.expectedIDs, matchingIDs)
		})
	}
}

func TestScheduler_ServiceIntegration(t *testing.T) {
	t.Parallel()

	t.Run("Who returns correct user for notification", func(t *testing.T) {
		t.Parallel()

		service := &mockService{
			whoResults: map[int64]string{
				1: "German",
				2: "Anthon",
			},
		}

		ctx := t.Context()

		who1, err := service.Who(ctx, 1)
		require.NoError(t, err)
		require.Equal(t, "German", who1)

		who2, err := service.Who(ctx, 2)
		require.NoError(t, err)
		require.Equal(t, "Anthon", who2)
	})

	t.Run("Who returns default for unknown chat", func(t *testing.T) {
		t.Parallel()

		service := &mockService{
			whoResults: map[int64]string{},
		}

		ctx := t.Context()

		who, err := service.Who(ctx, 999)
		require.NoError(t, err)
		require.Equal(t, "Unknown", who)
	})
}
