package main

import (
	"context"
	"fmt"

	"github.com/apus-run/sea-kit/got/v2"
)

func main() {
	var req = got.NewRequest(got.Post, "http://127.0.0.1:9090/upload")

	req.FileForm().AddFilePath("file1", "1.jpg", "./1.jpg")
	req.FileForm().AddFilePath("file2", "2.png", "./2.png")

	fmt.Println(req.Do(context.Background()))
}
