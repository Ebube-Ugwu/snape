package server

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/url"
	"strings"

	"github.com/ebube-ugwu/snape/internal/data"
)

//go:embed static/*
var staticFiles embed.FS

func Handler(store *data.Store) http.Handler {
	mux := http.NewServeMux()
	assets, _ := fs.Sub(staticFiles, "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(assets))))
	mux.HandleFunc("/api/tags", tags(store))
	mux.HandleFunc("/api/tags/", tagSnippets(store))
	mux.HandleFunc("/api/snippets", snippets(store))
	mux.HandleFunc("/api/snippets/", snippet(store))
	mux.HandleFunc("/", index)
	return mux
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFileFS(w, r, staticFiles, "static/index.html")
}

func tags(store *data.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		tags, err := store.ListTags()
		writeJSON(w, tags, err)
	}
}

func tagSnippets(store *data.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/api/tags/")
		tagName, ok := strings.CutSuffix(path, "/snippets")
		if !ok || tagName == "" {
			http.NotFound(w, r)
			return
		}
		tagName, err := url.PathUnescape(tagName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		snippets, err := store.SnippetsByTag(tagName)
		writeJSON(w, snippets, err)
	}
}

func snippets(store *data.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			query := r.URL.Query().Get("q")
			var (
				snippets []data.Snippet
				err      error
			)
			if query == "" {
				snippets, err = store.ListSnippets()
			} else {
				snippets, err = store.SearchSnippets(query)
			}
			writeJSON(w, snippets, err)
		case http.MethodPost:
			var snippet data.Snippet
			if err := json.NewDecoder(r.Body).Decode(&snippet); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			created, err := store.CreateSnippet(snippet)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			writeJSON(w, created, nil)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func snippet(store *data.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/api/snippets/")
		if name == "" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			item, err := store.GetSnippet(name)
			writeJSON(w, item, err)
		case http.MethodPut:
			var update data.Snippet
			if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			updated, err := store.UpdateSnippet(name, update)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			writeJSON(w, updated, nil)
		case http.MethodDelete:
			err := store.DeleteSnippet(name)
			writeJSON(w, map[string]bool{"deleted": err == nil}, err)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func writeJSON(w http.ResponseWriter, value any, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
