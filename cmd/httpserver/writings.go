package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"slices"
	"time"

	"github.com/igormichalak/website"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const WritingsDirPath = "./writings"

var Writings []Writing

type Writing struct {
	Title        string
	Slug         string
	Featured     bool
	Description  string
	Body         template.HTML
	PublishedAt  time.Time
	LastModified time.Time
}

func loadWriting(file fs.File, md goldmark.Markdown) (*Writing, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(data) < 3 || string(data[:3]) != "---" {
		return nil, nil
	}

	var buf bytes.Buffer
	ctx := parser.NewContext()

	if err := md.Convert(data, &buf, parser.WithContext(ctx)); err != nil {
		return nil, err
	}

	metadata := meta.Get(ctx)
	writing := &Writing{
		Featured: false,
		Body:     template.HTML(buf.String()),
	}

	title, ok := metadata["title"].(string)
	if !ok {
		return nil, fmt.Errorf("no valid title meta for %s", fileInfo.Name())
	}
	writing.Title = title

	slug, ok := metadata["slug"].(string)
	if !ok {
		return nil, fmt.Errorf("no valid slug meta for %s", fileInfo.Name())
	}
	writing.Slug = slug

	featured, ok := metadata["featured"].(bool)
	writing.Featured = ok && featured

	description, ok := metadata["description"].(string)
	if !ok {
		return nil, fmt.Errorf("no valid description meta for %s", fileInfo.Name())
	}
	writing.Description = description

	publishedAtStr, ok := metadata["published_at"].(string)
	if !ok {
		return nil, fmt.Errorf("no valid published_at meta for %s", fileInfo.Name())
	}
	publishedAt, err := time.Parse("2006-01-02", publishedAtStr)
	if err != nil {
		return nil, fmt.Errorf("can't parse published_at: %w", err)
	}
	writing.PublishedAt = publishedAt

	lastModifiedStr, ok := metadata["last_modified"].(string)
	if ok {
		lastModified, err := time.Parse("2006-01-02", lastModifiedStr)
		if err != nil {
			return nil, fmt.Errorf("can't parse last_modified: %w", err)
		}
		writing.LastModified = lastModified
	} else {
		writing.LastModified = publishedAt
	}

	return writing, nil
}

func init() {
	entries, err := fs.ReadDir(website.WritingsFS, ".")
	if err != nil {
		panic(fmt.Errorf("can't open %s, error: %w", WritingsDirPath, err))
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Footnote,
			extension.Typographer,
			meta.Meta,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	for _, entry := range entries {
		file, err := website.WritingsFS.Open(entry.Name())
		if err != nil {
			panic(err)
		}

		if entry.IsDir() {
			continue
		}

		writing, err := loadWriting(file, md)
		if err != nil {
			panic(fmt.Errorf("failed to load a writing: %w", err))
		}
		if writing == nil {
			fmt.Printf("Ignoring %s/%s\n", WritingsDirPath, entry.Name())
		} else {
			Writings = append(Writings, *writing)
		}

		if err := file.Close(); err != nil {
			panic(err)
		}
	}

	slices.SortFunc(Writings, func(a, b Writing) int {
		return b.PublishedAt.Compare(a.PublishedAt)
	})

	fmt.Printf("Loaded %d writing(s).\n", len(Writings))
}
