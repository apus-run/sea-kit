# Some utils for HTTP

## Uploader

Uploader is a tool support large files upload with no extra memory and tmp files.

```go
u := NewUploader().
		AddFile("file", filepath.Base(p), bytes.NewReader(bs)).
		SetParams(map[string][]string{
			"name": []string{"myname"},
        })
```

then

```go
rd, err := u.Body()
req, err := http.NewRequest("POST", url, rd)
req.Header.Add("Content-Type", u.ContentType())
```

If you want to customerize your boundary name,

```go
u := NewUploader().SetBoundary("mytestboundary").
		AddFile("file", filepath.Base(p), bytes.NewReader(bs)).
		SetParams(map[string][]string{
			"name": []string{"myname"},
        })
```
