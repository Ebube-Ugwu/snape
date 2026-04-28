package data

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtractVariables(t *testing.T) {
	got := ExtractVariables("copyright {{YEAR}} {{AUTHOR}} {{YEAR}} {{bad-key}}")
	want := []string{"AUTHOR", "YEAR"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExtractVariables() = %v, want %v", got, want)
	}
}

func TestRenderTemplate(t *testing.T) {
	got, err := RenderTemplate(
		"package {{NAME}} by {{AUTHOR}}",
		map[string]string{"NAME": "snape"},
		map[string]string{"AUTHOR": "Ebube"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if got != "package snape by Ebube" {
		t.Fatalf("RenderTemplate() = %q", got)
	}
}

func TestRenderTemplateReportsMissingValues(t *testing.T) {
	_, err := RenderTemplate("hello {{NAME}}", nil, nil)
	if err == nil {
		t.Fatal("expected missing variable error")
	}
}

func TestStoreSnippetLifecycle(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	created, err := store.CreateSnippet(Snippet{
		Name:        "license",
		Description: "MIT license header",
		Content:     "Copyright {{YEAR}} {{AUTHOR}}",
		Language:    "text",
		Type:        "license",
		Tags:        []string{"legal"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 {
		t.Fatal("expected created snippet ID")
	}

	got, err := store.GetSnippet("license")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "license" || got.Language != "text" || got.Type != "license" {
		t.Fatalf("unexpected snippet: %#v", got)
	}
	if !reflect.DeepEqual(got.Tags, []string{"legal"}) {
		t.Fatalf("tags = %v, want [legal]", got.Tags)
	}

	defaults, err := store.VariableDefaults(got.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := defaults["YEAR"]; !ok {
		t.Fatalf("expected YEAR variable, got %v", defaults)
	}

	results, err := store.SearchSnippets("MIT")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Name != "license" {
		t.Fatalf("SearchSnippets() = %#v", results)
	}

	if err := store.DeleteSnippet("license"); err != nil {
		t.Fatal(err)
	}
}

func TestUpsertSnippetReplacesExistingSnippetAndTags(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	if _, err := store.CreateSnippet(Snippet{
		Name:    "hello",
		Content: "old",
		Tags:    []string{"old"},
	}); err != nil {
		t.Fatal(err)
	}

	updated, err := store.UpsertSnippet(Snippet{
		Name:        "hello",
		Description: "updated snippet",
		Content:     "new",
		Language:    "go",
		Type:        "inline",
		Tags:        []string{"new", "test"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Content != "new" || updated.Description != "updated snippet" {
		t.Fatalf("unexpected updated snippet: %#v", updated)
	}
	if !reflect.DeepEqual(updated.Tags, []string{"new", "test"}) {
		t.Fatalf("tags = %v, want [new test]", updated.Tags)
	}
}

func TestListTagsAndSnippetsByTag(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	if _, err := store.CreateSnippet(Snippet{Name: "go-one", Content: "one", Tags: []string{"go", "cli"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.CreateSnippet(Snippet{Name: "go-two", Content: "two", Tags: []string{"go"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.CreateSnippet(Snippet{Name: "web-one", Content: "web", Tags: []string{"web"}}); err != nil {
		t.Fatal(err)
	}

	tags, err := store.ListTags()
	if err != nil {
		t.Fatal(err)
	}
	gotTags := []string{}
	for _, tag := range tags {
		gotTags = append(gotTags, tag.Name)
	}
	if !reflect.DeepEqual(gotTags, []string{"cli", "go", "web"}) {
		t.Fatalf("tags = %v, want [cli go web]", gotTags)
	}

	snippets, err := store.SnippetsByTag("go")
	if err != nil {
		t.Fatal(err)
	}
	gotNames := []string{}
	for _, snippet := range snippets {
		gotNames = append(gotNames, snippet.Name)
	}
	if !reflect.DeepEqual(gotNames, []string{"go-two", "go-one"}) {
		t.Fatalf("snippets = %v, want [go-two go-one]", gotNames)
	}
}

func newTestStore(t *testing.T) *Store {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "snape-test.db")
	MigrateDB("sqlite://" + dbPath)
	db := OpenDB(dbPath)
	return NewStore(db)
}
