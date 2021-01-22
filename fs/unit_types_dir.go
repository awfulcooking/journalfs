package fs

import (
	"context"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/togetherbeer/journalfs/journalcache"
)

type UnitTypesDir struct {
	journalCache *journalcache.JournalCache

	typeDirs map[string]*UnitLogsDir
}

var _ fs.Node = (*UnitTypesDir)(nil)
var _ fs.HandleReadDirAller = (*UnitTypesDir)(nil)
var _ fs.NodeStringLookuper = (*UnitTypesDir)(nil)

func (d *UnitTypesDir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Mode = os.ModeDir | 0o550
	attr.Size = uint64(len(d.typeDirs))

	return nil
}

func (d *UnitTypesDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var dirs []fuse.Dirent

	for name, _ := range d.typeDirs {
		dirs = append(dirs, fuse.Dirent{
			Name: name,
			Type: fuse.DT_Dir,
		})
	}

	return dirs, nil
}

func (d *UnitTypesDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if node, ok := d.typeDirs[name]; ok {
		return node, nil
	}

	return nil, syscall.ENOENT
}

func NewUnitTypesDir(journalCache *journalcache.JournalCache) *UnitTypesDir {
	return &UnitTypesDir{
		journalCache: journalCache,
		typeDirs:     makeUnitTypeDirs(journalCache),
	}
}

var UNIT_TYPES = []string{
	"automount",
	"device",
	"mount",
	"path",
	"scope",
	"service",
	"slice",
	"socket",
	"swap",
	"target",
	"timer",
}

func makeUnitTypeDirs(journalCache *journalcache.JournalCache) map[string]*UnitLogsDir {
	var dirs = map[string]*UnitLogsDir{}

	for _, name := range UNIT_TYPES {
		dirs[name] = NewUnitLogsDir(journalCache, name)
	}

	return dirs
}
