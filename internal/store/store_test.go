package store

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

var (
	testStore *Store
	testRDB   *redis.Client
)

// getRedisAddr reads the Redis address from the environment or returns a default.
func getRedisAddr() string {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return addr
}

// TestMain sets up the test environment by connecting to a test Redis DB
// and cleaning it up after tests have run.
func TestMain(m *testing.M) {
	// Use a different Redis DB for testing to isolate from development data
	testRDB = redis.NewClient(&redis.Options{
		Addr: getRedisAddr(),
		DB:   15, // Test-specific DB
	})

	// Ping to ensure the connection is alive
	if _, err := testRDB.Ping(context.Background()).Result(); err != nil {
		slog.Error("failed to connect to test redis", "error", err)
		os.Exit(1)
	}

	// Initialize a logger that discards output for clean test runs
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	testStore = NewStore(testRDB, logger)

	// Run the tests
	code := m.Run()

	// Clean up the test database and close the connection
	if err := testRDB.FlushDB(context.Background()).Err(); err != nil {
		slog.Error("failed to flush test db", "error", err)
	}
	if err := testRDB.Close(); err != nil {
		slog.Error("failed to close test redis connection", "error", err)
	}

	os.Exit(code)
}

// setupTest is a helper to ensure a clean state before each test.
func setupTest(t *testing.T) {
	// Flush the DB to ensure no data leaks between tests
	err := testRDB.FlushDB(context.Background()).Err()
	require.NoError(t, err, "Failed to flush test DB")
}

func TestSetEstablish(t *testing.T) {
	setupTest(t)
	chatID := int64(1)
	users := []string{"user1", "user2", "user3"}

	err := testStore.SetEstablish(chatID, users)
	require.NoError(t, err)

	storedUsers, err := testStore.GetUsers(chatID)
	require.NoError(t, err)
	require.Equal(t, users, storedUsers)

	idx, err := testStore.GetActiveIndex(chatID)
	require.NoError(t, err)
	require.Equal(t, 0, idx)
}

func TestAddUser(t *testing.T) {
	setupTest(t)
	chatID := int64(2)
	initialUsers := []string{"userA"}
	err := testStore.SetEstablish(chatID, initialUsers)
	require.NoError(t, err)

	// Add a new user
	err = testStore.AddUser(chatID, "userB")
	require.NoError(t, err)

	storedUsers, err := testStore.GetUsers(chatID)
	require.NoError(t, err)
	require.Equal(t, []string{"userA", "userB"}, storedUsers)

	// Try to add an existing user (should fail)
	err = testStore.AddUser(chatID, "userA")
	require.Error(t, err, "Expected error when adding duplicate user")
}

func TestRemoveUser(t *testing.T) {
	setupTest(t)
	chatID := int64(3)
	users := []string{"user1", "user2", "user3", "user2"}
	err := testStore.SetEstablish(chatID, users)
	require.NoError(t, err)

	// Remove all instances of "user2"
	err = testStore.RemoveUser(chatID, "user2")
	require.NoError(t, err)

	storedUsers, err := testStore.GetUsers(chatID)
	require.NoError(t, err)
	require.Equal(t, []string{"user1", "user3"}, storedUsers)
}

func TestWhoNextPrev(t *testing.T) {
	setupTest(t)
	chatID := int64(4)
	users := []string{"a", "b", "c"}
	err := testStore.SetEstablish(chatID, users)
	require.NoError(t, err)

	// Test Who
	name, err := testStore.Who(chatID)
	require.NoError(t, err)
	require.Equal(t, "a", name, "Initial user should be 'a'")

	// Test Next
	name, err = testStore.Next(chatID)
	require.NoError(t, err)
	require.Equal(t, "b", name, "Next user should be 'b'")

	name, err = testStore.Next(chatID)
	require.NoError(t, err)
	require.Equal(t, "c", name, "Next user should be 'c'")

	// Test Next (looping)
	name, err = testStore.Next(chatID)
	require.NoError(t, err)
	require.Equal(t, "a", name, "Next user should loop back to 'a'")

	// Test Prev
	name, err = testStore.Prev(chatID)
	require.NoError(t, err)
	require.Equal(t, "c", name, "Prev user should be 'c'")

	name, err = testStore.Prev(chatID)
	require.NoError(t, err)
	require.Equal(t, "b", name, "Prev user should be 'b'")

	// Test Prev (looping)
	name, err = testStore.Prev(chatID)
	require.NoError(t, err)
	require.Equal(t, "a", name, "Prev user should be 'a'")

	name, err = testStore.Prev(chatID)
	require.NoError(t, err)
	require.Equal(t, "c", name, "Prev user should loop back to 'c'")
}
