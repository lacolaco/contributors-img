package image

import (
	"context"
	"testing"

	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/go/renderer"
)

// mockCache is a simple implementation of appcache.AppCache for testing
type mockCache struct {
	getFunc      func(context.Context, string) (model.FileHandle, error)
	getJSONFunc  func(context.Context, string, any) error
	saveFunc     func(context.Context, string, []byte, string) error
	saveJSONFunc func(context.Context, string, any) error
}

func (m *mockCache) Get(ctx context.Context, key string) (model.FileHandle, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, key)
	}
	return nil, nil
}

func (m *mockCache) GetJSON(ctx context.Context, key string, v any) error {
	if m.getJSONFunc != nil {
		return m.getJSONFunc(ctx, key, v)
	}
	return nil
}

func (m *mockCache) Save(ctx context.Context, key string, data []byte, contentType string) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, key, data, contentType)
	}
	return nil
}

func (m *mockCache) SaveJSON(ctx context.Context, key string, v any) error {
	if m.saveJSONFunc != nil {
		return m.saveJSONFunc(ctx, key, v)
	}
	return nil
}

// mockDataURLConverter is a helper function for testing that avoids making real HTTP requests
func mockDataURLConverter(_ context.Context, avatarURL string, _ map[string]string) (string, error) {
	return "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=", nil
}

func TestService_normalizeContributors(t *testing.T) {
	// Store the original converter to restore it after the test
	originalConverter := dataURLConverter
	// Set the mock converter for testing
	dataURLConverter = mockDataURLConverter
	// Restore the original converter after the test
	defer func() { dataURLConverter = originalConverter }()

	service := &Service{
		cache: &mockCache{},
	}

	repo := &model.Repository{
		Owner:    "test-owner",
		RepoName: "test-repo",
	}

	contributors := []*model.Contributor{
		{ID: 1, Login: "user1", AvatarURL: "https://avatar1.com"},
		{ID: 2, Login: "user2", AvatarURL: "https://avatar2.com"},
		{ID: 0, Login: "anonymous", AvatarURL: "https://avatar3.com"}, // Anonymous user
	}

	input := &model.RepositoryContributors{
		Repository:      repo,
		StargazersCount: 100,
		Contributors:    contributors,
	}

	options := &renderer.RendererOptions{
		MaxCount: 2,
		Columns:  4,
		ItemSize: 64,
	}

	t.Run("with includeAnonymous=true", func(t *testing.T) {
		result, err := service.normalizeContributors(context.Background(), input, options, true)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(result.Contributors) != 2 {
			t.Errorf("Expected 2 contributors, got %d", len(result.Contributors))
		}

		if result.Contributors[0].ID != 1 {
			t.Errorf("Expected first contributor ID to be 1, got %d", result.Contributors[0].ID)
		}

		// Verify data URL conversion was applied
		if result.Contributors[0].AvatarURL == "https://avatar1.com" {
			t.Errorf("Avatar URL should have been converted to a data URL")
		}
	})

	t.Run("with includeAnonymous=false", func(t *testing.T) {
		result, err := service.normalizeContributors(context.Background(), input, options, false)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(result.Contributors) != 2 {
			t.Errorf("Expected 2 contributors, got %d", len(result.Contributors))
		}

		// Check that anonymous user was filtered out
		for _, c := range result.Contributors {
			if c.ID == 0 {
				t.Errorf("Anonymous contributor should have been filtered out")
			}
		}
	})
}
