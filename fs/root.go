package fs

import (
	"bazil.org/fuse/fs"

	"github.com/togetherbeer/journalfs/journalcache"
)

type Root struct {
	journalCache *journalcache.JournalCache

	*UnitTypesDir
}

var _ fs.Node = (*Root)(nil)
var _ fs.HandleReadDirAller = (*Root)(nil)
var _ fs.NodeStringLookuper = (*Root)(nil)

func NewRoot(journalCache *journalcache.JournalCache) *Root {
	return &Root{
		journalCache: journalCache,
		UnitTypesDir: NewUnitTypesDir(journalCache),
	}
}
