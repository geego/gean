package geanlib

import (
	"fmt"
	"math/rand"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/gostores/require"
	"yiqilai.tech/gean/app/deps"
)

const pageCollectionsPageTemplate = `---
title: "%s"
categories:
- Hugo
---
# Doc
`

func BenchmarkGetPage(b *testing.B) {
	var (
		cfg, fs = newTestCfg()
		r       = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	for i := 0; i < 10; i++ {
		for j := 0; j < 100; j++ {
			writeSource(b, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), fmt.Sprintf("page%d.md", j)), "CONTENT")
		}
	}

	s := buildSingleSite(b, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	pagePaths := make([]string, b.N)

	for i := 0; i < b.N; i++ {
		pagePaths[i] = fmt.Sprintf("sect%d", r.Intn(10))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		home := s.getPage(KindHome)
		if home == nil {
			b.Fatal("Home is nil")
		}

		p := s.getPage(KindSection, pagePaths[i])
		if p == nil {
			b.Fatal("Section is nil")
		}

	}
}

func BenchmarkGetPageRegular(b *testing.B) {
	var (
		cfg, fs = newTestCfg()
		r       = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	for i := 0; i < 10; i++ {
		for j := 0; j < 100; j++ {
			content := fmt.Sprintf(pageCollectionsPageTemplate, fmt.Sprintf("Title%d_%d", i, j))
			writeSource(b, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), fmt.Sprintf("page%d.md", j)), content)
		}
	}

	s := buildSingleSite(b, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	pagePaths := make([]string, b.N)

	for i := 0; i < b.N; i++ {
		pagePaths[i] = path.Join(fmt.Sprintf("sect%d", r.Intn(10)), fmt.Sprintf("page%d.md", r.Intn(100)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		page := s.getPage(KindPage, pagePaths[i])
		require.NotNil(b, page)
	}
}

func TestGetPage(t *testing.T) {

	var (
		assert  = require.New(t)
		cfg, fs = newTestCfg()
	)

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			content := fmt.Sprintf(pageCollectionsPageTemplate, fmt.Sprintf("Title%d_%d", i, j))
			writeSource(t, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), fmt.Sprintf("page%d.md", j)), content)
		}
	}

	content := fmt.Sprintf(pageCollectionsPageTemplate, "UniqueBase")
	writeSource(t, fs, filepath.Join("content", "sect3", "unique.md"), content)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	tests := []struct {
		kind          string
		path          []string
		expectedTitle string
	}{
		{KindHome, []string{}, ""},
		{KindSection, []string{"sect3"}, "Sect3s"},
		{KindPage, []string{"sect3", "page1.md"}, "Title3_1"},
		{KindPage, []string{"sect4/page2.md"}, "Title4_2"},
		{KindPage, []string{filepath.FromSlash("sect5/page3.md")}, "Title5_3"},
		// Ref/Relref supports this potentially ambiguous lookup.
		{KindPage, []string{"unique.md"}, "UniqueBase"},
	}

	for i, test := range tests {
		errorMsg := fmt.Sprintf("Test %d", i)
		page := s.getPage(test.kind, test.path...)
		assert.NotNil(page, errorMsg)
		assert.Equal(test.kind, page.Kind)
		assert.Equal(test.expectedTitle, page.Title)
	}

}
