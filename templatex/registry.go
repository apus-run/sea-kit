// Package templatex is a thin wrapper around the standard html/template
// and text/template packages that implements a convenient registry to
// load and cache templates on the fly concurrently.
//
// It was created to assist the JSVM plugin HTML rendering, but could be used in other Go code.
//
// Example:
//
//	registry := templatex.NewRegistry()
//
//	html1, err := registry.LoadFiles(
//		// the files set wil be parsed only once and then cached
//		"layout.html",
//		"content.html",
//	).Render(map[string]any{"name": "John"})
//
//	html2, err := registry.LoadFiles(
//		// reuse the already parsed and cached files set
//		"layout.html",
//		"content.html",
//	).Render(map[string]any{"name": "Jane"})
package templatex

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hashicorp/golang-lru/v2"
)

// NewRegistry creates and initializes a new templates registry with
// some defaults (eg. global "raw" template function for unescaped HTML).
//
// Use the Registry.Load* methods to load templates into the registry.
func NewRegistry() *Registry {
	l, err := lru.New[string, *Renderer](65536)
	if err != nil {
		return nil
	}
	return &Registry{
		cache: l,
		funcs: template.FuncMap{
			"raw": func(str string) template.HTML {
				return template.HTML(str)
			},
		},
	}
}

// Registry defines a templates registry that is safe to be used by multiple goroutines.
//
// Use the Registry.Load* methods to load templates into the registry.
type Registry struct {
	cache *lru.Cache[string, *Renderer]
	funcs template.FuncMap
}

// AddFuncs registers new global template functions.
//
// The key of each map entry is the function name that will be used in the templates.
// If a function with the map entry name already exists it will be replaced with the new one.
//
// The value of each map entry is a function that must have either a
// single return value, or two return values of which the second has type error.
//
// Example:
//
//	r.AddFuncs(map[string]any{
//	  "toUpper": func(str string) string {
//	      return strings.ToUppser(str)
//	  },
//	  ...
//	})
func (r *Registry) AddFuncs(funcs map[string]any) *Registry {
	for name, f := range funcs {
		r.funcs[name] = f
	}

	return r
}

// LoadFiles caches (if not already) the specified filenames set as a
// single template and returns a ready to use Renderer instance.
//
// There must be at least 1 filename specified.
func (r *Registry) LoadFiles(filenames ...string) *Renderer {
	key := strings.Join(filenames, ",")

	v, ok := r.cache.Get(key)
	if ok {
		return v
	}

	// parse and cache
	tpl, err := template.New(filepath.Base(filenames[0])).Funcs(r.funcs).ParseFiles(filenames...)
	v = &Renderer{template: tpl, parseError: err}
	r.cache.Add(key, v)

	return v
}

// LoadString caches (if not already) the specified inline string as a
// single template and returns a ready to use Renderer instance.
func (r *Registry) LoadString(text string) *Renderer {
	v, ok := r.cache.Get(text)
	if ok {
		return v
	}

	// parse and cache (using the text as key)
	tpl, err := template.New("").Funcs(r.funcs).Parse(text)
	v = &Renderer{template: tpl, parseError: err}
	r.cache.Add(text, v)

	return v
}

// LoadFS caches (if not already) the specified fs and globPatterns
// pair as single template and returns a ready to use Renderer instance.
//
// There must be at least 1 file matching the provided globPattern(s)
// (note that most file names serves as glob patterns matching themselves).
func (r *Registry) LoadFS(fsys fs.FS, globPatterns ...string) *Renderer {
	key := fmt.Sprintf("%v%v", fsys, globPatterns)

	v, ok := r.cache.Get(key)
	if ok {
		return v
	}

	// find the first file to use as template name (it is required when specifying Funcs)
	var firstFilename string
	if len(globPatterns) > 0 {
		list, _ := fs.Glob(fsys, globPatterns[0])
		if len(list) > 0 {
			firstFilename = filepath.Base(list[0])
		}
	}

	tpl, err := template.New(firstFilename).Funcs(r.funcs).ParseFS(fsys, globPatterns...)
	v = &Renderer{template: tpl, parseError: err}
	r.cache.Add(key, v)

	return v
}
