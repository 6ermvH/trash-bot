package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/6ermvH/trash-bot/internal/repository"
	_ "modernc.org/sqlite"
)

type RepoSQLite struct {
	db *sql.DB
}

func New(dbPath string) (*RepoSQLite, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite db: %w", err)
	}

	repo := &RepoSQLite{db: db}
	if err := repo.migrate(); err != nil {
		return nil, fmt.Errorf("migrate sqlite db: %w", err)
	}

	return repo, nil
}

func (r *RepoSQLite) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS chats (
		id INTEGER PRIMARY KEY,
		current INTEGER NOT NULL DEFAULT 0,
		users TEXT NOT NULL DEFAULT '[]'
	);`

	_, err := r.db.Exec(query)
	return err
}

func (r *RepoSQLite) Close() error {
	return r.db.Close()
}

func (r *RepoSQLite) GetChats(ctx context.Context) []repository.Chat {
	rows, err := r.db.QueryContext(ctx, "SELECT id, current, users FROM chats")
	if err != nil {
		return []repository.Chat{}
	}
	defer rows.Close()

	var chats []repository.Chat
	for rows.Next() {
		var chat repository.Chat
		var usersJSON string
		if err := rows.Scan(&chat.ID, &chat.Current, &usersJSON); err != nil {
			continue
		}
		json.Unmarshal([]byte(usersJSON), &chat.Users)
		chats = append(chats, chat)
	}

	if chats == nil {
		return []repository.Chat{}
	}
	return chats
}

func (r *RepoSQLite) GetChat(ctx context.Context, chatID int64) (*repository.Chat, error) {
	var chat repository.Chat
	var usersJSON string

	err := r.db.QueryRowContext(ctx,
		"SELECT id, current, users FROM chats WHERE id = ?", chatID,
	).Scan(&chat.ID, &chat.Current, &usersJSON)

	if err == sql.ErrNoRows {
		return nil, repository.ErrChatIsNotInitialize
	}
	if err != nil {
		return nil, fmt.Errorf("query chat: %w", err)
	}

	json.Unmarshal([]byte(usersJSON), &chat.Users)
	return &chat, nil
}

func (r *RepoSQLite) GetCurrent(ctx context.Context, chatID int64) (string, error) {
	chat, err := r.GetChat(ctx, chatID)
	if err != nil {
		return "", err
	}

	if len(chat.Users) == 0 {
		return "", repository.ErrChatIsEmpty
	}

	return chat.Users[chat.Current], nil
}

func (r *RepoSQLite) SetNext(ctx context.Context, chatID int64) error {
	chat, err := r.GetChat(ctx, chatID)
	if err != nil {
		return err
	}

	if len(chat.Users) == 0 {
		return repository.ErrChatIsEmpty
	}

	newCurrent := (chat.Current + 1) % len(chat.Users)
	_, err = r.db.ExecContext(ctx,
		"UPDATE chats SET current = ? WHERE id = ?", newCurrent, chatID,
	)
	return err
}

func (r *RepoSQLite) SetPrev(ctx context.Context, chatID int64) error {
	chat, err := r.GetChat(ctx, chatID)
	if err != nil {
		return err
	}

	if len(chat.Users) == 0 {
		return repository.ErrChatIsEmpty
	}

	newCurrent := (len(chat.Users) + chat.Current - 1) % len(chat.Users)
	_, err = r.db.ExecContext(ctx,
		"UPDATE chats SET current = ? WHERE id = ?", newCurrent, chatID,
	)
	return err
}

func (r *RepoSQLite) SetEstablish(ctx context.Context, chatID int64, users []string) error {
	usersJSON, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("marshal users: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO chats (id, current, users) VALUES (?, 0, ?)
		ON CONFLICT(id) DO UPDATE SET current = 0, users = ?
	`, chatID, string(usersJSON), string(usersJSON))

	return err
}

func (r *RepoSQLite) Subscribe(ctx context.Context, chatID int64) error {
	return nil
}

func (r *RepoSQLite) Unsubscribe(ctx context.Context, chatID int64) error {
	return nil
}
