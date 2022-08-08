package renderer

import (
	"strings"
	"testing"

	"contrib.rocks/libs/goutils/model"
	"github.com/bradleyjkemp/cupaloy"
)

func TestRender_Snapshot(t *testing.T) {
	r := NewRenderer(&RendererOptions{MaxCount: 12, Columns: 4, ItemSize: 64})
	ret := r.Render(&model.RepositoryContributors{
		Repository:      &model.Repository{Owner: "owner", RepoName: "name"},
		StargazersCount: 100,
		Contributors: []*model.Contributor{
			{Login: "login1", AvatarURL: "avatar1", Contributions: 1},
			{Login: "login2", AvatarURL: "avatar2", Contributions: 2},
		},
	})
	svgStr := string(ret.Bytes())
	cupaloy.New(cupaloy.SnapshotFileExtension(".svg")).SnapshotT(t, svgStr)
}

func TestRender_Title(t *testing.T) {
	r := NewRenderer(&RendererOptions{MaxCount: 12, Columns: 4, ItemSize: 64})
	ret := r.Render(&model.RepositoryContributors{
		Repository:      &model.Repository{Owner: "owner", RepoName: "name"},
		StargazersCount: 100,
		Contributors: []*model.Contributor{
			{Login: "login1", AvatarURL: "avatar1", Contributions: 1},
			{Login: "login2", AvatarURL: "avatar2", Contributions: 2},
		},
	})
	svgStr := string(ret.Bytes())
	for _, s := range []string{"login1", "login2"} {
		if !strings.Contains(svgStr, s) {
			t.Log(svgStr)
			t.Fatal("title not found")
		}
	}
}

func TestSvgImage_Size(t *testing.T) {
	img := svgImage([]byte("<svg></svg>"))
	if img.Size() != 11 {
		t.Fatalf("size not correct: %d", img.Size())
	}
}

func TestSvgImage_ContentType(t *testing.T) {
	img := svgImage([]byte("<svg></svg>"))
	if img.ContentType() != "image/svg+xml" {
		t.Fatalf("content type not correct: %s", img.ContentType())
	}
}

func TestSvgImage_Etag(t *testing.T) {
	img := svgImage([]byte("<svg></svg>"))
	if img.ContentType() != "image/svg+xml" {
		t.Fatalf("content type not correct: %s", img.ContentType())
	}
}
