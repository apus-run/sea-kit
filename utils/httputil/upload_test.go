package httputil

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fileInfo(p string) (int, string, error) {
	h := sha1.New()
	bs, err := ioutil.ReadFile(p)
	if err != nil {
		return 0, "", err
	}
	h.Write(bs)
	return len(bs), fmt.Sprintf("%x", h.Sum(nil)), nil
}

func TestUpload(t *testing.T) {
	testUploadFile(t, "./upload.go", "./upload_test.go")
}

func testUploadFile(t *testing.T, filenames ...string) {
	u := NewUploader().
		SetBoundary("mytestboundary").
		SetParams(map[string][]string{
			"name": []string{"myname"},
		})
	var files = make(map[string][]byte)
	for i, filename := range filenames {
		f, err := os.Open(filename)
		assert.NoError(t, err)
		defer f.Close()

		bs, err := ioutil.ReadAll(f)
		assert.NoError(t, err)
		files[filename] = bs
		u.AddFile(fmt.Sprintf("file%d", i), filepath.Base(filename), bytes.NewReader(bs))
	}

	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	rd, err := u.Body()
	assert.NoError(t, err)

	content, err := ioutil.ReadAll(rd)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "http://localhost:8000/api/v1/files/upload", bytes.NewReader(content))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", u.ContentType())

	var s = func(resp http.ResponseWriter, req *http.Request) {
		req.ParseMultipartForm(1024 * 1024 * 1024)
		assert.EqualValues(t, req.Form.Get("name"), "myname", string(content))

		for i, filename := range filenames {
			file, head, err := req.FormFile(fmt.Sprintf("file%d", i))
			assert.NoError(t, err, string(content))

			assert.EqualValues(t, filepath.Base(filename), head.Filename, string(content))
			ct, err := ioutil.ReadAll(file)
			assert.NoError(t, err, string(content))
			assert.EqualValues(t, files[filename], ct, string(content))

			fmt.Println(string(content))
		}
	}

	http.HandlerFunc(s).ServeHTTP(recorder, req)
	assert.EqualValues(t, recorder.Code, http.StatusOK)
}
