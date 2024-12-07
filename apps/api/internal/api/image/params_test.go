package image

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_GetImageParams_BindQuery_AllParams(t *testing.T) {
	uri := url.URL{Path: "/image"}
	q := uri.Query()
	q.Set("repo", "angular/angular-ja")
	q.Set("max", "100")
	q.Set("columns", "12")
	q.Set("itemSize", "64")
	q.Set("anon", "1")
	q.Set("preview", "1")
	uri.RawQuery = q.Encode()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", uri.String(), nil)

	var params GetImageParams
	err := params.bind(c)

	if err != nil {
		t.Fatalf(err.Error())
	}
	if params.Repository != "angular/angular-ja" {
		t.Fatalf("Bound params: %v", params)
	}
	if params.MaxCount != 100 {
		t.Fatalf("Bound params: %v", params)
	}
	if params.Columns != 12 {
		t.Fatalf("Bound params: %v", params)
	}
	if params.ItemSize != 64 {
		t.Fatalf("Bound params: %v", params)
	}
	if !params.IncludeAnonymous {
		t.Fatalf("Bound params: %v", params)
	}
	if !params.Preview {
		t.Fatalf("Bound params: %v", params)
	}
}

func Test_GetImageParams_BingQuery_NoRepository(t *testing.T) {
	uri := url.URL{Path: "/image"}
	q := uri.Query()
	uri.RawQuery = q.Encode()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", uri.String(), nil)

	var params GetImageParams
	err := params.bind(c)

	if err == nil {
		t.Fatalf(err.Error())
	}
}

func Test_GetImageParams_BingQuery_InvalidRepository(t *testing.T) {
	uri := url.URL{Path: "/image"}
	q := uri.Query()
	q.Set("repo", "angular-ja")
	uri.RawQuery = q.Encode()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", uri.String(), nil)

	var params GetImageParams
	err := params.bind(c)

	if err == nil {
		t.Fatalf(err.Error())
	}
}

func Test_GetImageParams_BingQuery_Anon_ZeroIsFalse(t *testing.T) {
	uri := url.URL{Path: "/image"}
	q := uri.Query()
	q.Set("anon", "0")
	uri.RawQuery = q.Encode()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", uri.String(), nil)

	var params GetImageParams
	err := params.bind(c)

	if err == nil {
		t.Fatalf(err.Error())
	}
	if params.IncludeAnonymous {
		t.Fatalf("Bound params: %v", params)
	}
}

func Test_GetImageParams_BingQuery_Anon_FalseIsFalse(t *testing.T) {
	uri := url.URL{Path: "/image"}
	q := uri.Query()
	q.Set("anon", "false")
	uri.RawQuery = q.Encode()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", uri.String(), nil)

	var params GetImageParams
	err := params.bind(c)

	if err == nil {
		t.Fatalf(err.Error())
	}
	if params.IncludeAnonymous {
		t.Fatalf("Bound params: %v", params)
	}
}
