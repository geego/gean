package geanlib

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/geego/gean/app/helpers"
	"github.com/gostores/gitmap"
)

func (h *HugoSites) assembleGitInfo() {
	if !h.Cfg.GetBool("enableGitInfo") {
		return
	}

	var (
		workingDir = h.Cfg.GetString("workingDir")
		contentDir = h.Cfg.GetString("contentDir")
	)

	gitRepo, err := gitmap.Map(workingDir, "")
	if err != nil {
		h.Log.ERROR.Printf("Got error reading Git log: %s", err)
		return
	}

	gitMap := gitRepo.Files
	repoPath := filepath.FromSlash(gitRepo.TopLevelAbsPath)

	// The Hugo site may be placed in a sub folder in the Git repo,
	// one example being the Hugo docs.
	// We have to find the root folder to the Hugo site below the Git root.
	contentRoot := strings.TrimPrefix(workingDir, repoPath)
	contentRoot = strings.TrimPrefix(contentRoot, helpers.FilePathSeparator)

	s := h.Sites[0]

	for _, p := range s.AllPages {
		if p.Path() == "" {
			// Home page etc. with no content file.
			continue
		}
		// Git normalizes file paths on this form:
		filename := path.Join(filepath.ToSlash(contentRoot), contentDir, filepath.ToSlash(p.Path()))
		g, ok := gitMap[filename]
		if !ok {
			h.Log.WARN.Printf("Failed to find GitInfo for %q", filename)
			return
		}

		p.GitInfo = g
		p.Lastmod = p.GitInfo.AuthorDate
	}

}
