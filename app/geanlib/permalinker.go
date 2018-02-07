package geanlib

var (
	_ Permalinker = (*Page)(nil)
	_ Permalinker = (*OutputFormat)(nil)
)

// Permalinker provides permalinks of both the relative and absolute kind.
type Permalinker interface {
	Permalink() string
	RelPermalink() string
}
