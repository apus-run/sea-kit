package ginx

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBind(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{"title":"我是标题", "email":"@163.com"}`))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)
	gc := Context{Context: c}
	var obj struct {
		Title string  `json:"title" binding:"required,max=32,min=4" label:"标题"`
		Email *string `json:"email" binding:"required,email" label:"邮箱"`
	}
	assert.Equal(t, gc.Bind(&obj).Error(), "Email必须是一个有效的邮箱")
	assert.Equal(t, w.Code, 400)
	t.Log("Code:", w.Code, "Body:", w.Body.String())
	assert.Empty(t, c.Errors)
}

func TestShouldBind(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{"title":"t", "email":"@shimo.im"}`))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)
	gc := Context{Context: c}
	var obj struct {
		Title string  `json:"title" binding:"required,max=32,min=4" label:"标题"`
		Email *string `json:"email" binding:"required,email" label:"邮箱"`
	}
	assert.Equal(t, gc.ShouldBind(&obj).Error(), "标题长度必须至少为4个字符|邮箱必须是一个有效的邮箱")
	assert.Equal(t, w.Code, 200)
	t.Log("Code:", w.Code, "Body:", w.Body.String())
	assert.Empty(t, c.Errors)
}

func TestContext_Query(t *testing.T) {
	testCases := []struct {
		name    string
		req     func(t *testing.T) *http.Request
		key     string
		wantVal any
	}{
		{
			name: "获得数据",
			req: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "http://localhost/abc?name=123&age=18", nil)
				require.NoError(t, err)
				return req
			},
			key:     "name",
			wantVal: "123",
		},
		{
			name: "没有数据",
			req: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "http://localhost/abc?name=123&age=18", nil)
				require.NoError(t, err)
				return req
			},
			key:     "nickname",
			wantVal: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &Context{Context: &gin.Context{
				Request: tc.req(t),
			}}
			val := ctx.Query(tc.key)
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestContext_Param(t *testing.T) {
	testCases := []struct {
		name    string
		req     func(t *testing.T) *http.Request
		key     string
		wantErr error
		wantVal any
	}{
		{
			name: "获得数据",
			req: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "http://localhost/hello?name=123&age=18", nil)
				req.Form = url.Values{}
				req.Form.Set("name", "world")
				require.NoError(t, err)
				return req
			},
			key:     "name",
			wantVal: "world",
		},
		{
			name: "没有数据",
			req: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "http://localhost/hello?name=123&age=18", nil)
				require.NoError(t, err)
				return req
			},
			key:     "nickname",
			wantVal: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := gin.Default()
			server.POST("/hello", func(context *gin.Context) {
				ctx := &Context{Context: context}
				val := ctx.Param(tc.key)
				assert.Equal(t, tc.wantVal, val)
			})
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, tc.req(t))
		})
	}
}

func TestContext_Cookie(t *testing.T) {
	testCases := []struct {
		name    string
		req     func(t *testing.T) *http.Request
		key     string
		wantErr error
		wantVal any
	}{
		{
			name: "有cookie",
			req: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "http://localhost/hello?name=123&age=18", nil)
				req.AddCookie(&http.Cookie{
					Name:  "name",
					Value: "world",
				})
				require.NoError(t, err)
				return req
			},
			key:     "name",
			wantVal: "world",
		},
		{
			name: "没有 cookie",
			req: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "http://localhost/hello?name=123&age=18", nil)
				require.NoError(t, err)
				return req
			},
			key:     "nickname",
			wantVal: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := gin.Default()
			server.POST("/hello", func(context *gin.Context) {
				ctx := &Context{Context: context}
				val := ctx.Param(tc.key)
				assert.Equal(t, tc.wantVal, val)
			})
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, tc.req(t))
		})
	}
}
