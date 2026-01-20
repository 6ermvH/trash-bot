package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/6ermvH/trash-bot/internal/repository"
	_ "modernc.org/sqlite"
)

type RepoSQLite struct {
	db *sql.DB
}

func New(dbPath string) (*RepoSQLite, error) {
	dbConn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	ctx := context.Background()
	if err := dbConn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping sqlite db: %w", err)
	}

	repo := &RepoSQLite{db: dbConn}
	if err := repo.migrate(ctx); err != nil {
		return nil, fmt.Errorf("migrate sqlite db: %w", err)
	}

	return repo, nil
}

func (r *RepoSQLite) Close() error {
	if err := r.db.Close(); err != nil {
		return fmt.Errorf("close sqlite db: %w", err)
	}

	return nil
}

func (r *RepoSQLite) GetChats(ctx context.Context) (_ []repository.Chat, err error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, current, users, notify_time FROM chats")
	if err != nil {
		return nil, fmt.Errorf("query chats: %w", err)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close rows: %w", closeErr)
		}
	}()

	chats := make([]repository.Chat, 0)

	for rows.Next() {
		var (
			chat       repository.Chat
			usersJSON  string
			notifyTime sql.NullString
		)

		if err := rows.Scan(&chat.ID, &chat.Current, &usersJSON, &notifyTime); err != nil {
			return nil, fmt.Errorf("scan chat: %w", err)
		}

		if err := json.Unmarshal([]byte(usersJSON), &chat.Users); err != nil {
			return nil, fmt.Errorf("decode chat users: %w", err)
		}

		if notifyTime.Valid {
			chat.NotifyTime = &notifyTime.String
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate chats: %w", err)
	}

	return chats, nil
}

func (r *RepoSQLite) GetChat(ctx context.Context, chatID int64) (*repository.Chat, error) {
	var (
		chat       repository.Chat
		usersJSON  string
		notifyTime sql.NullString
	)

	err := r.db.QueryRowContext(ctx,
		"SELECT id, current, users, notify_time FROM chats WHERE id = ?", chatID,
	).Scan(&chat.ID, &chat.Current, &usersJSON, &notifyTime)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrChatIsNotInitialize
	}

	if err != nil {
		return nil, fmt.Errorf("query chat: %w", err)
	}

	if err := json.Unmarshal([]byte(usersJSON), &chat.Users); err != nil {
		return nil, fmt.Errorf("decode chat users: %w", err)
	}

	if notifyTime.Valid {
		chat.NotifyTime = &notifyTime.String
	}

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

	if _, err := r.db.ExecContext(
		ctx,
		"UPDATE chats SET current = ? WHERE id = ?",
		newCurrent,
		chatID,
	); err != nil {
		return fmt.Errorf("update current: %w", err)
	}

	return nil
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

	if _, err := r.db.ExecContext(
		ctx,
		"UPDATE chats SET current = ? WHERE id = ?",
		newCurrent,
		chatID,
	); err != nil {
		return fmt.Errorf("update current: %w", err)
	}

	return nil
}

func (r *RepoSQLite) SetEstablish(ctx context.Context, chatID int64, users []string) error {
	usersJSON, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("marshal users: %w", err)
	}

	if _, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO chats (id, current, users) VALUES (?, 0, ?)
		ON CONFLICT(id) DO UPDATE SET current = 0, users = ?
	`,
		chatID,
		string(usersJSON),
		string(usersJSON),
	); err != nil {
		return fmt.Errorf("upsert chat: %w", err)
	}

	return nil
}

func (r *RepoSQLite) Subscribe(ctx context.Context, chatID int64, notifyTime string) error {
	if _, err := r.db.ExecContext(
		ctx,
		"UPDATE chats SET notify_time = ? WHERE id = ?",
		notifyTime,
		chatID,
	); err != nil {
		return fmt.Errorf("update notify_time: %w", err)
	}

	return nil
}

func (r *RepoSQLite) Unsubscribe(ctx context.Context, chatID int64) error {
	if _, err := r.db.ExecContext(
		ctx,
		"UPDATE chats SET notify_time = NULL WHERE id = ?",
		chatID,
	); err != nil {
		return fmt.Errorf("clear notify_time: %w", err)
	}

	return nil
}

func (r *RepoSQLite) GetSubscribedChats(ctx context.Context) (_ []repository.Chat, err error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, current, users, notify_time FROM chats WHERE notify_time IS NOT NULL",
	)
	if err != nil {
		return nil, fmt.Errorf("query subscribed chats: %w", err)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close rows: %w", closeErr)
		}
	}()

	chats := make([]repository.Chat, 0)

	for rows.Next() {
		var (
			chat       repository.Chat
			usersJSON  string
			notifyTime sql.NullString
		)

		if err := rows.Scan(&chat.ID, &chat.Current, &usersJSON, &notifyTime); err != nil {
			return nil, fmt.Errorf("scan chat: %w", err)
		}

		if err := json.Unmarshal([]byte(usersJSON), &chat.Users); err != nil {
			return nil, fmt.Errorf("decode chat users: %w", err)
		}

		if notifyTime.Valid {
			chat.NotifyTime = &notifyTime.String
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate subscribed chats: %w", err)
	}

	return chats, nil
}

func (r *RepoSQLite) migrate(ctx context.Context) error {
	createTable := `
	CREATE TABLE IF NOT EXISTS chats (
		id INTEGER PRIMARY KEY,
		current INTEGER NOT NULL DEFAULT 0,
		users TEXT NOT NULL DEFAULT '[]',
		notify_time TEXT DEFAULT NULL
	);`

	if _, err := r.db.ExecContext(ctx, createTable); err != nil {
		return fmt.Errorf("exec create table migration: %w", err)
	}

	// Добавляем колонку notify_time если она не существует (для существующих БД)
	addColumn := `ALTER TABLE chats ADD COLUMN notify_time TEXT DEFAULT NULL;`
	// Игнорируем ошибку, если колонка уже существует
	_, _ = r.db.ExecContext(ctx, addColumn)

	return nil
}
