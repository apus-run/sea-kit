package static

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

type Builder struct {
	urlPrefix string
	root      string
}

func NewBuilder() *Builder {
	return &Builder{
		urlPrefix: "/",
		root:      "./static",
	}
}

func (b *Builder) Build() gin.HandlerFunc {
	fs := LocalFile(b.root, false)
	fileServer := http.FileServer(fs)
	if b.urlPrefix != "" {
		fileServer = http.StripPrefix(b.urlPrefix, fileServer)
	}
	return func(c *gin.Context) {
		if fs.Exists(b.urlPrefix, c.Request.URL.Path) {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

const INDEX = "index.html"

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func LocalFile(root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := path.Join(name, INDEX)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

//  static.NewBuilder().Build("/", "./dist")
