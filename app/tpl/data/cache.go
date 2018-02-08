package data

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/geego/gean/app/config"
	"github.com/geego/gean/app/helpers"
	"github.com/gostores/fsintra"
)

var cacheMu sync.RWMutex

// getCacheFileID returns the cache ID for a string.
func getCacheFileID(cfg config.Provider, id string) string {
	hash := md5.Sum([]byte(id))
	return cfg.GetString("cacheDir") + hex.EncodeToString(hash[:])
}

// getCache returns the content for an ID from the file cache or an error.
// If the ID is not found, return nil,nil.
func getCache(id string, fs fsintra.Fs, cfg config.Provider, ignoreCache bool) ([]byte, error) {
	if ignoreCache {
		return nil, nil
	}

	cacheMu.RLock()
	defer cacheMu.RUnlock()

	fID := getCacheFileID(cfg, id)
	isExists, err := helpers.Exists(fID, fs)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, nil
	}

	return fsintra.ReadFile(fs, fID)
}

// writeCache writes bytes associated with an ID into the file cache.
func writeCache(id string, c []byte, fs fsintra.Fs, cfg config.Provider, ignoreCache bool) error {
	if ignoreCache {
		return nil
	}

	cacheMu.Lock()
	defer cacheMu.Unlock()

	fID := getCacheFileID(cfg, id)
	f, err := fs.Create(fID)
	if err != nil {
		return errors.New("Error: " + err.Error() + ". Failed to create file: " + fID)
	}
	defer f.Close()

	n, err := f.Write(c)
	if err != nil {
		return errors.New("Error: " + err.Error() + ". Failed to write to file: " + fID)
	}
	if n == 0 {
		return errors.New("No bytes written to file: " + fID)
	}
	return nil
}

func deleteCache(id string, fs fsintra.Fs, cfg config.Provider) error {
	return fs.Remove(getCacheFileID(cfg, id))
}
