package cachekey

import (
	"testing"

	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/go/renderer"
)

func TestForRepository(t *testing.T) {
	repo := &model.Repository{
		Owner:    "owner",
		RepoName: "repo",
	}

	key := ForRepository(repo, "json")
	expected := "repo/owner--repo.json"

	if key != expected {
		t.Errorf("Expected key to be %s, got %s", expected, key)
	}
}

func TestForContributors(t *testing.T) {
	repo := &model.Repository{
		Owner:    "owner",
		RepoName: "repo",
	}

	key := ForContributors(repo, "json")
	expected := "contributors/v1.2/owner--repo.json"

	if key != expected {
		t.Errorf("Expected key to be %s, got %s", expected, key)
	}
}

func TestForImage(t *testing.T) {
	repo := &model.Repository{
		Owner:    "owner",
		RepoName: "repo",
	}

	options := &renderer.RendererOptions{
		MaxCount: 100,
		Columns:  12,
	}

	t.Run("with includeAnonymous=true", func(t *testing.T) {
		key := ForImage(repo, options, "svg", true)
		expected := "image/owner--repo--anon_100_12.svg"

		if key != expected {
			t.Errorf("Expected key to be %s, got %s", expected, key)
		}
	})

	t.Run("with includeAnonymous=false", func(t *testing.T) {
		key := ForImage(repo, options, "svg", false)
		expected := "image/owner--repo--noanon_100_12.svg"

		if key != expected {
			t.Errorf("Expected key to be %s, got %s", expected, key)
		}
	})
}
