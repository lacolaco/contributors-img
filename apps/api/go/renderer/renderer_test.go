package renderer

import (
	"crypto/md5"
	"fmt"
	"strings"
	"testing"

	"contrib.rocks/apps/api/go/model"
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
	t.Run("should be image/svg+xml", func(tt *testing.T) {
		img := svgImage([]byte("<svg></svg>"))
		if img.ContentType() != "image/svg+xml" {
			tt.Fatalf("content type not correct: %s", img.ContentType())
		}
	})
}

func TestSvgImage_ETag(t *testing.T) {
	t.Run("should be a MD5 hash", func(tt *testing.T) {
		img := svgImage([]byte("<svg></svg>"))
		if img.ETag() != fmt.Sprintf("%x", md5.Sum([]byte("<svg></svg>"))) {
			tt.Fatalf("etag not correct: %s", img.ETag())
		}
	})
}
