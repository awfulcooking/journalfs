package fs

import (
	"context"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"github.com/togetherbeer/journalfs/journalcache"
)

type Root struct {
	journalCache *journalcache.JournalCache

	dirEntries []fuse.Dirent
}

var _ fs.Node = (*Root)(nil)
var _ fs.HandleReadDirAller = (*Root)(nil)
var _ fs.NodeStringLookuper = (*Root)(nil)

func (r *Root) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = 1
	attr.Mode = os.ModeDir | 0o555
	attr.Size = 0

	return nil
}

func (r *Root) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return []fuse.Dirent{
		{Name: "by-unit", Type: fuse.DT_Dir},
	}, nil
}

func (r *Root) Lookup(ctx context.Context, name string) (fs.Node, error) {
	switch name {
	case "by-unit":
		return NewDirUnits(r.journalCache), nil
	}
	return nil, syscall.ENOENT
}

func NewRoot(journalCache *journalcache.JournalCache) *Root {
	return &Root{
		journalCache: journalCache,
	}
}
