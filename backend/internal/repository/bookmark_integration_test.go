package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/hoard-app/backend/internal/model"
)

// setupTestDB connects to the database named by TEST_DATABASE_URL and creates
// a fresh, self-contained `bookmarks` table (no auth.users FK, so it runs
// against any plain PostgreSQL). The test is skipped when the env var is unset.
//
// Run locally with, e.g.:
//
//	TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres" go test ./internal/repository/...
func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping repository integration test")
	}

	ctx := context.Background()

	// Pin search_path to our isolated schema on EVERY pooled connection, so
	// the repository's unqualified `FROM bookmarks` always resolves here
	// regardless of which connection the pool hands out.
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse TEST_DATABASE_URL: %v", err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = "hoard_test"

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	// Isolated schema so we never touch real data and clean up fully.
	if _, err := pool.Exec(ctx, `DROP SCHEMA IF EXISTS hoard_test CASCADE; CREATE SCHEMA hoard_test;`); err != nil {
		pool.Close()
		t.Fatalf("failed to create test schema: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		CREATE TABLE hoard_test.bookmarks (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL,
			url TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			favicon_url TEXT DEFAULT '',
			og_image_url TEXT DEFAULT '',
			site_name TEXT DEFAULT '',
			content_type TEXT DEFAULT 'article',
			is_archived BOOLEAN NOT NULL DEFAULT false,
			is_favorite BOOLEAN NOT NULL DEFAULT false,
			read_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`); err != nil {
		pool.Close()
		t.Fatalf("failed to create bookmarks table: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DROP SCHEMA IF EXISTS hoard_test CASCADE;`)
		pool.Close()
	})

	return pool
}

const (
	userA = "11111111-1111-1111-1111-111111111111"
	userB = "22222222-2222-2222-2222-222222222222"
)

func TestBookmarkRepository_Integration(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewBookmarkRepository(pool)
	ctx := context.Background()

	// Seed a mix of content types for user A and one for user B.
	seed := []struct {
		user string
		req  model.CreateBookmarkRequest
	}{
		{userA, model.CreateBookmarkRequest{URL: "https://a.com/1", Title: "Article One", ContentType: "article"}},
		{userA, model.CreateBookmarkRequest{URL: "https://a.com/2", Title: "Video One", ContentType: "video"}},
		{userA, model.CreateBookmarkRequest{URL: "https://a.com/3", Title: "Video Two", ContentType: "video"}},
		{userA, model.CreateBookmarkRequest{URL: "https://a.com/4", Title: "Default Type"}}, // defaults to article
		{userB, model.CreateBookmarkRequest{URL: "https://b.com/1", Title: "B's Video", ContentType: "video"}},
	}
	for _, s := range seed {
		if _, err := repo.Create(ctx, s.user, s.req); err != nil {
			t.Fatalf("seed create failed: %v", err)
		}
	}

	t.Run("List with no filter returns all of user's bookmarks", func(t *testing.T) {
		items, total, err := repo.List(ctx, userA, 20, 0, "")
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total != 4 {
			t.Errorf("got total %d; want 4", total)
		}
		if len(items) != 4 {
			t.Errorf("got %d items; want 4", len(items))
		}
	})

	t.Run("List filters by content_type=video", func(t *testing.T) {
		items, total, err := repo.List(ctx, userA, 20, 0, "video")
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total != 2 {
			t.Errorf("got total %d; want 2", total)
		}
		for _, b := range items {
			if b.ContentType != "video" {
				t.Errorf("got content_type %q; want video", b.ContentType)
			}
		}
	})

	t.Run("List filters by content_type=article (includes default)", func(t *testing.T) {
		_, total, err := repo.List(ctx, userA, 20, 0, "article")
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total != 2 {
			t.Errorf("got total %d; want 2 (explicit + defaulted)", total)
		}
	})

	t.Run("List is scoped to the requesting user", func(t *testing.T) {
		items, total, err := repo.List(ctx, userB, 20, 0, "")
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total != 1 || len(items) != 1 {
			t.Fatalf("got total=%d len=%d; want 1/1", total, len(items))
		}
		if items[0].UserID != userB {
			t.Errorf("leaked another user's data: %q", items[0].UserID)
		}
	})

	t.Run("List pagination via limit/offset", func(t *testing.T) {
		page1, total, err := repo.List(ctx, userA, 2, 0, "")
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total != 4 {
			t.Errorf("got total %d; want 4", total)
		}
		if len(page1) != 2 {
			t.Errorf("page1 got %d items; want 2", len(page1))
		}
		page2, _, err := repo.List(ctx, userA, 2, 2, "")
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(page2) != 2 {
			t.Errorf("page2 got %d items; want 2", len(page2))
		}
		if page1[0].ID == page2[0].ID {
			t.Error("pagination returned overlapping rows")
		}
	})

	t.Run("Create, GetByID, Update, Delete round trip", func(t *testing.T) {
		created, err := repo.Create(ctx, userA, model.CreateBookmarkRequest{
			URL: "https://a.com/crud", Title: "CRUD Target", ContentType: "social",
		})
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if created.ContentType != "social" {
			t.Errorf("got content_type %q; want social", created.ContentType)
		}

		got, err := repo.GetByID(ctx, userA, created.ID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if got.Title != "CRUD Target" {
			t.Errorf("got title %q; want CRUD Target", got.Title)
		}

		newTitle := "CRUD Updated"
		fav := true
		updated, err := repo.Update(ctx, userA, created.ID, model.UpdateBookmarkRequest{
			Title: &newTitle, IsFavorite: &fav,
		})
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}
		if updated.Title != "CRUD Updated" || !updated.IsFavorite {
			t.Errorf("update not applied: title=%q favorite=%v", updated.Title, updated.IsFavorite)
		}

		if err := repo.Delete(ctx, userA, created.ID); err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
		if _, err := repo.GetByID(ctx, userA, created.ID); err == nil {
			t.Error("expected error fetching deleted bookmark, got nil")
		}
	})

	t.Run("GetByID enforces row-level ownership", func(t *testing.T) {
		created, err := repo.Create(ctx, userA, model.CreateBookmarkRequest{
			URL: "https://a.com/private", Title: "Private",
		})
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		// user B must not be able to read user A's bookmark.
		if _, err := repo.GetByID(ctx, userB, created.ID); err == nil {
			t.Error("expected error when other user fetches bookmark, got nil")
		}
	})
}
