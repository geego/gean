package command

import (
	"os"
	"path/filepath"

	"github.com/gostores/fsnotify"
	"github.com/gostores/fsync"

	"github.com/geego/gean/app/helpers"

	src "github.com/geego/gean/app/source"
)

type staticSyncer struct {
	c *commandeer
	d *src.Dirs
}

func newStaticSyncer(c *commandeer) (*staticSyncer, error) {
	dirs, err := src.NewDirs(c.Fs, c.Cfg, c.DepsCfg.Logger)
	if err != nil {
		return nil, err
	}

	return &staticSyncer{c: c, d: dirs}, nil
}

func (s *staticSyncer) isStatic(path string) bool {
	return s.d.IsStatic(path)
}

func (s *staticSyncer) syncsStaticEvents(staticEvents []fsnotify.Event) error {
	c := s.c

	syncFn := func(dirs *src.Dirs, publishDir string) error {
		staticSourceFs, err := dirs.CreateStaticFs()
		if err != nil {
			return err
		}

		if staticSourceFs == nil {
			c.Logger.WARN.Println("No static directories found to sync")
			return nil
		}

		syncer := fsync.NewSyncer()
		syncer.NoTimes = c.Cfg.GetBool("noTimes")
		syncer.NoChmod = c.Cfg.GetBool("noChmod")
		syncer.SrcFs = staticSourceFs
		syncer.DestFs = c.Fs.Destination

		// prevent spamming the log on changes
		logger := helpers.NewDistinctFeedbackLogger()

		for _, ev := range staticEvents {
			// Due to our approach of layering both directories and the content's rendered output
			// into one we can't accurately remove a file not in one of the source directories.
			// If a file is in the local static dir and also in the theme static dir and we remove
			// it from one of those locations we expect it to still exist in the destination
			//
			// If gean generates a file (from the content dir) over a static file
			// the content generated file should take precedence.
			//
			// Because we are now watching and handling individual events it is possible that a static
			// event that occupies the same path as a content generated file will take precedence
			// until a regeneration of the content takes places.
			//
			// gean assumes that these cases are very rare and will permit this bad behavior
			// The alternative is to track every single file and which pipeline rendered it
			// and then to handle conflict resolution on every event.

			fromPath := ev.Name

			// If we are here we already know the event took place in a static dir
			relPath := dirs.MakeStaticPathRelative(fromPath)
			if relPath == "" {
				// Not member of this virtual host.
				continue
			}

			// Remove || rename is harder and will require an assumption.
			// gean takes the following approach:
			// If the static file exists in any of the static source directories after this event
			// gean will re-sync it.
			// If it does not exist in all of the static directories gean will remove it.
			//
			// This assumes that gean has not generated content on top of a static file and then removed
			// the source of that static file. In this case gean will incorrectly remove that file
			// from the published directory.
			if ev.Op&fsnotify.Rename == fsnotify.Rename || ev.Op&fsnotify.Remove == fsnotify.Remove {
				if _, err := staticSourceFs.Stat(relPath); os.IsNotExist(err) {
					// If file doesn't exist in any static dir, remove it
					toRemove := filepath.Join(publishDir, relPath)

					logger.Println("File no longer exists in static dir, removing", toRemove)
					_ = c.Fs.Destination.RemoveAll(toRemove)
				} else if err == nil {
					// If file still exists, sync it
					logger.Println("Syncing", relPath, "to", publishDir)

					if err := syncer.Sync(filepath.Join(publishDir, relPath), relPath); err != nil {
						c.Logger.ERROR.Println(err)
					}
				} else {
					c.Logger.ERROR.Println(err)
				}

				continue
			}

			// For all other event operations gean will sync static.
			logger.Println("Syncing", relPath, "to", publishDir)
			if err := syncer.Sync(filepath.Join(publishDir, relPath), relPath); err != nil {
				c.Logger.ERROR.Println(err)
			}
		}

		return nil
	}

	return c.doWithPublishDirs(syncFn)

}
