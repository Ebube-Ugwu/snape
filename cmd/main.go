package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ebube-ugwu/snape/internal/data"
	"github.com/ebube-ugwu/snape/internal/server"
)

func main() {
	log.SetFlags(0)
	data.MigrateDB(data.DbURI)
	db := data.OpenDB(data.DiskDbPath)
	defer db.Close()

	store := data.NewStore(db)
	if err := run(os.Args[1:], store); err != nil {
		log.Fatal(err)
	}
}

func run(args []string, store *data.Store) error {
	if len(args) == 0 {
		usage()
		return nil
	}

	switch args[0] {
	case "add":
		return add(args[1:], store)
	case "get":
		return get(args[1:], store)
	case "list":
		return list(args[1:], store)
	case "search":
		return search(args[1:], store)
	case "delete":
		return remove(args[1:], store)
	case "tag":
		return tag(args[1:], store)
	case "insert":
		return insert(args[1:], store)
	case "import":
		return importSnippets(args[1:], store)
	case "export":
		return exportSnippets(args[1:], store)
	case "serve":
		return serve(args[1:], store)
	default:
		usage()
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func add(args []string, store *data.Store) error {
	name, flagArgs := leadingArg(args)
	flags := flag.NewFlagSet("add", flag.ExitOnError)
	description := flags.String("description", "", "snippet description")
	language := flags.String("language", "", "snippet language")
	snippetType := flags.String("type", "inline", "snippet type")
	content := flags.String("content", "", "snippet content")
	file := flags.String("file", "", "read snippet content from file")
	tags := stringSliceFlag{}
	flags.Var(&tags, "tag", "tag to attach; repeat for multiple tags")
	if err := flags.Parse(flagArgs); err != nil {
		return err
	}
	if name == "" && flags.NArg() == 1 {
		name = flags.Arg(0)
	}
	if name == "" || flags.NArg() > 0 {
		return fmt.Errorf("usage: snape add <name> [--tag tag] [--content text|--file path]")
	}

	body, err := resolveContent(*content, *file)
	if err != nil {
		return err
	}
	_, err = store.CreateSnippet(data.Snippet{
		Name:        name,
		Description: *description,
		Content:     body,
		Language:    *language,
		Type:        *snippetType,
		Tags:        tags,
	})
	return err
}

func get(args []string, store *data.Store) error {
	name, flagArgs := leadingArg(args)
	flags := flag.NewFlagSet("get", flag.ExitOnError)
	vars := varFlags{}
	flags.Var(&vars, "var", "template variable in KEY=VALUE form")
	if err := flags.Parse(flagArgs); err != nil {
		return err
	}
	if name == "" && flags.NArg() == 1 {
		name = flags.Arg(0)
	}
	if name == "" || flags.NArg() > 0 {
		return fmt.Errorf("usage: snape get <name> [--var KEY=VALUE]")
	}
	snippet, err := store.UseSnippet(name)
	if err != nil {
		return err
	}
	defaults, err := store.VariableDefaults(snippet.ID)
	if err != nil {
		return err
	}
	rendered, err := data.RenderTemplate(snippet.Content, vars, defaults)
	if err != nil {
		return err
	}
	fmt.Println(rendered)
	return nil
}

func list(args []string, store *data.Store) error {
	flags := flag.NewFlagSet("list", flag.ExitOnError)
	verbose := false
	flags.BoolVar(&verbose, "v", false, "show descriptions and tags")
	flags.BoolVar(&verbose, "verbose", false, "show descriptions and tags")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("usage: snape list [-v|--verbose]")
	}

	snippets, err := store.ListSnippets()
	if err != nil {
		return err
	}
	if verbose {
		rows := [][]tableCell{}
		for _, snippet := range snippets {
			rows = append(rows, []tableCell{
				{text: snippet.Name, color: nameColor},
				{text: emptyDash(snippet.Language)},
				{text: emptyDash(snippet.Type)},
				{text: emptyDash(strings.Join(snippet.Tags, ", ")), color: tagColor},
				{text: emptyDash(snippet.Description)},
			})
		}
		printTable([]tableCell{
			{text: "NAME", color: headerColor},
			{text: "LANGUAGE", color: headerColor},
			{text: "TYPE", color: headerColor},
			{text: "TAGS", color: headerColor},
			{text: "DESCRIPTION", color: headerColor},
		}, rows)
		return nil
	}

	rows := [][]tableCell{}
	for _, snippet := range snippets {
		rows = append(rows, []tableCell{
			{text: snippet.Name, color: nameColor},
			{text: emptyDash(snippet.Language)},
			{text: emptyDash(strings.Join(snippet.Tags, ", ")), color: tagColor},
		})
	}
	printTable([]tableCell{
		{text: "NAME", color: headerColor},
		{text: "LANGUAGE", color: headerColor},
		{text: "TAGS", color: headerColor},
	}, rows)
	return nil
}

func search(args []string, store *data.Store) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: snape search <query>")
	}
	snippets, err := store.SearchSnippets(args[0])
	if err != nil {
		return err
	}
	for _, snippet := range snippets {
		fmt.Printf("%s\t%s\t%s\n", snippet.Name, snippet.Language, snippet.Description)
	}
	return nil
}

func remove(args []string, store *data.Store) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: snape delete <name>")
	}
	return store.DeleteSnippet(args[0])
}

func tag(args []string, store *data.Store) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: snape tag <name> <tag>")
	}
	return store.TagSnippet(args[0], args[1])
}

func insert(args []string, store *data.Store) error {
	positionals, flagArgs := leadingArgs(args, 2)
	flags := flag.NewFlagSet("insert", flag.ExitOnError)
	line := flags.Int("line", 0, "1-based line number to insert before")
	vars := varFlags{}
	flags.Var(&vars, "var", "template variable in KEY=VALUE form")
	if err := flags.Parse(flagArgs); err != nil {
		return err
	}
	positionals = append(positionals, flags.Args()...)
	if len(positionals) != 2 {
		return fmt.Errorf("usage: snape insert <name> <file> [--line N] [--var KEY=VALUE]")
	}
	snippet, err := store.UseSnippet(positionals[0])
	if err != nil {
		return err
	}
	defaults, err := store.VariableDefaults(snippet.ID)
	if err != nil {
		return err
	}
	rendered, err := data.RenderTemplate(snippet.Content, vars, defaults)
	if err != nil {
		return err
	}
	return insertIntoFile(positionals[1], rendered, *line)
}

func serve(args []string, store *data.Store) error {
	flags := flag.NewFlagSet("serve", flag.ExitOnError)
	addr := flags.String("addr", ":7777", "listen address")
	if err := flags.Parse(args); err != nil {
		return err
	}
	fmt.Println("serving snape on http://localhost" + *addr)
	return http.ListenAndServe(*addr, server.Handler(store))
}

func exportSnippets(args []string, store *data.Store) error {
	flags := flag.NewFlagSet("export", flag.ExitOnError)
	output := flags.String("output", "", "output JSON file")
	flags.StringVar(output, "o", "", "output JSON file")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *output == "" && flags.NArg() == 1 {
		*output = flags.Arg(0)
	}
	if *output == "" || flags.NArg() > 1 {
		return fmt.Errorf("usage: snape export --output snippets.json")
	}

	snippets, err := store.ListSnippets()
	if err != nil {
		return err
	}
	payload := snippetExport{
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		Snippets:   snippets,
	}
	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	bytes = append(bytes, '\n')
	return os.WriteFile(*output, bytes, 0644)
}

func importSnippets(args []string, store *data.Store) error {
	flags := flag.NewFlagSet("import", flag.ExitOnError)
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("usage: snape import snippets.json")
	}

	bytes, err := os.ReadFile(flags.Arg(0))
	if err != nil {
		return err
	}
	snippets, err := decodeSnippetImport(bytes)
	if err != nil {
		return err
	}
	for _, snippet := range snippets {
		if _, err := store.UpsertSnippet(snippet); err != nil {
			return fmt.Errorf("import %q: %w", snippet.Name, err)
		}
	}
	fmt.Printf("imported %d snippets\n", len(snippets))
	return nil
}

func resolveContent(content string, file string) (string, error) {
	if file != "" {
		bytes, err := os.ReadFile(file)
		return string(bytes), err
	}
	if content != "" {
		return content, nil
	}
	if pipedStdin() {
		bytes, err := io.ReadAll(os.Stdin)
		return string(bytes), err
	}
	return editContent()
}

func editContent() (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "", fmt.Errorf("provide --content, --file, stdin, or set EDITOR")
	}
	tmp, err := os.CreateTemp("", "snape-*.md")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())
	if err := tmp.Close(); err != nil {
		return "", err
	}
	cmd := exec.Command(editor, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	bytes, err := os.ReadFile(tmp.Name())
	return string(bytes), err
}

func insertIntoFile(path string, content string, line int) error {
	original, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if line <= 0 {
		return os.WriteFile(path, appendWithNewline(string(original), content), 0644)
	}
	lines := splitLines(string(original))
	if line > len(lines)+1 {
		return fmt.Errorf("line %d is outside file with %d lines", line, len(lines))
	}
	index := line - 1
	lines = append(lines[:index], append([]string{content}, lines[index:]...)...)
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

func appendWithNewline(original string, content string) []byte {
	if original == "" || strings.HasSuffix(original, "\n") {
		return []byte(original + content)
	}
	return []byte(original + "\n" + content)
}

func splitLines(content string) []string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func pipedStdin() bool {
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice == 0
}

func usage() {
	fmt.Println(`snape commands:
  add <name> [--tag tag] [--content text|--file path]
  get <name> [--var KEY=VALUE]
  list [-v|--verbose]
  search <query>
  delete <name>
  tag <name> <tag>
  insert <name> <file> [--line N] [--var KEY=VALUE]
  import <file.json>
  export --output <file.json>
  serve [--addr :7777]`)
}

func leadingArg(args []string) (string, []string) {
	positionals, rest := leadingArgs(args, 1)
	if len(positionals) == 0 {
		return "", rest
	}
	return positionals[0], rest
}

func leadingArgs(args []string, count int) ([]string, []string) {
	positionals := []string{}
	index := 0
	for index < len(args) && len(positionals) < count && !strings.HasPrefix(args[index], "-") {
		positionals = append(positionals, args[index])
		index++
	}
	return positionals, args[index:]
}

type varFlags map[string]string

func (v varFlags) String() string {
	return fmt.Sprint(map[string]string(v))
}

func (v varFlags) Set(value string) error {
	key, val, ok := strings.Cut(value, "=")
	if !ok || key == "" {
		return fmt.Errorf("expected KEY=VALUE, got %q", value)
	}
	v[key] = val
	return nil
}

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	*s = append(*s, value)
	return nil
}

const (
	headerColor = "\033[1;36m"
	nameColor   = "\033[1;32m"
	tagColor    = "\033[35m"
	resetColor  = "\033[0m"
)

func color(code string, value string) string {
	if value == "" {
		return "-"
	}
	if code == "" {
		return value
	}
	return code + value + resetColor
}

func emptyDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

type snippetExport struct {
	ExportedAt string         `json:"exported_at"`
	Snippets   []data.Snippet `json:"snippets"`
}

func decodeSnippetImport(bytes []byte) ([]data.Snippet, error) {
	var wrapped snippetExport
	if err := json.Unmarshal(bytes, &wrapped); err == nil && wrapped.Snippets != nil {
		return wrapped.Snippets, nil
	}
	var snippets []data.Snippet
	if err := json.Unmarshal(bytes, &snippets); err != nil {
		return nil, err
	}
	return snippets, nil
}

type tableCell struct {
	text  string
	color string
}

func printTable(headers []tableCell, rows [][]tableCell) {
	widths := tableWidths(headers, rows)
	border := tableBorder(widths)
	fmt.Println(border)
	printTableRow(headers, widths)
	fmt.Println(border)
	for _, row := range rows {
		printTableRow(row, widths)
	}
	fmt.Println(border)
}

func tableWidths(headers []tableCell, rows [][]tableCell) []int {
	widths := make([]int, len(headers))
	for index, header := range headers {
		widths[index] = visibleLen(header.text)
	}
	for _, row := range rows {
		for index, cell := range row {
			if index < len(widths) && visibleLen(cell.text) > widths[index] {
				widths[index] = visibleLen(cell.text)
			}
		}
	}
	return widths
}

func tableBorder(widths []int) string {
	var builder strings.Builder
	builder.WriteByte('+')
	for _, width := range widths {
		builder.WriteString(strings.Repeat("-", width+2))
		builder.WriteByte('+')
	}
	return builder.String()
}

func printTableRow(row []tableCell, widths []int) {
	var builder strings.Builder
	builder.WriteByte('|')
	for index, width := range widths {
		cell := tableCell{text: "-"}
		if index < len(row) {
			cell = row[index]
		}
		text := color(cell.color, cell.text)
		builder.WriteByte(' ')
		builder.WriteString(text)
		builder.WriteString(strings.Repeat(" ", width-visibleLen(cell.text)))
		builder.WriteByte(' ')
		builder.WriteByte('|')
	}
	fmt.Println(builder.String())
}

func visibleLen(value string) int {
	return utf8.RuneCountInString(stripANSI(value))
}

func stripANSI(value string) string {
	var builder strings.Builder
	inEscape := false
	for index := 0; index < len(value); index++ {
		char := value[index]
		if inEscape {
			if char >= '@' && char <= '~' {
				inEscape = false
			}
			continue
		}
		if char == 0x1b {
			inEscape = true
			continue
		}
		builder.WriteByte(char)
	}
	return builder.String()
}
