package fs

import (
	bzfs "bazil.org/fuse/fs"

	"github.com/togetherbeer/journalfs/journalcache"
)

type FS struct {
	journalCache *journalcache.JournalCache
}

var _ bzfs.FS = (*FS)(nil)

func (fs *FS) Root() (bzfs.Node, error) {
	return NewRoot(fs.journalCache), nil
}

func NewFS(journalCache *journalcache.JournalCache) *FS {
	return &FS{
		journalCache: journalCache,
	}
}
