package data

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/gostores/fsintra"
	"github.com/gostores/notepad"

	"github.com/geego/gean/app/config"
	"github.com/geego/gean/app/helpers"
)

var (
	remoteURLLock = &remoteLock{m: make(map[string]*sync.Mutex)}
	resSleep      = time.Second * 2 // if JSON decoding failed sleep for n seconds before retrying
	resRetries    = 1               // number of retries to load the JSON from URL or local file system
)

type remoteLock struct {
	sync.RWMutex
	m map[string]*sync.Mutex
}

// URLLock locks an URL during download
func (l *remoteLock) URLLock(url string) {
	var (
		lock *sync.Mutex
		ok   bool
	)
	l.Lock()
	if lock, ok = l.m[url]; !ok {
		lock = &sync.Mutex{}
		l.m[url] = lock
	}
	l.Unlock()
	lock.Lock()
}

// URLUnlock unlocks an URL when the download has been finished. Use only in defer calls.
func (l *remoteLock) URLUnlock(url string) {
	l.RLock()
	defer l.RUnlock()
	if um, ok := l.m[url]; ok {
		um.Unlock()
	}
}

// getRemote loads the content of a remote file. This method is thread safe.
func getRemote(req *http.Request, fs fsintra.Fs, cfg config.Provider, hc *http.Client) ([]byte, error) {
	url := req.URL.String()

	c, err := getCache(url, fs, cfg, cfg.GetBool("ignoreCache"))
	if err != nil {
		return nil, err
	}
	if c != nil {
		return c, nil
	}

	// avoid race condition with locks, block other goroutines if the current url is processing
	remoteURLLock.URLLock(url)
	defer func() { remoteURLLock.URLUnlock(url) }()

	// avoid multiple locks due to calling getCache twice
	c, err = getCache(url, fs, cfg, cfg.GetBool("ignoreCache"))
	if err != nil {
		return nil, err
	}
	if c != nil {
		return c, nil
	}

	notepad.INFO.Printf("Downloading: %s ...", url)
	res, err := hc.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("Failed to retrieve remote file: %s", http.StatusText(res.StatusCode))
	}

	c, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	err = writeCache(url, c, fs, cfg, cfg.GetBool("ignoreCache"))
	if err != nil {
		return nil, err
	}

	notepad.INFO.Printf("... and cached to: %s", getCacheFileID(cfg, url))
	return c, nil
}

// getLocal loads the content of a local file
func getLocal(url string, fs fsintra.Fs, cfg config.Provider) ([]byte, error) {
	filename := filepath.Join(cfg.GetString("workingDir"), url)
	if e, err := helpers.Exists(filename, fs); !e {
		return nil, err
	}

	return fsintra.ReadFile(fs, filename)

}

// getResource loads the content of a local or remote file
func (ns *Namespace) getResource(req *http.Request) ([]byte, error) {
	switch req.URL.Scheme {
	case "":
		return getLocal(req.URL.String(), ns.deps.Fs.Source, ns.deps.Cfg)
	default:
		return getRemote(req, ns.deps.Fs.Source, ns.deps.Cfg, ns.client)
	}
}
