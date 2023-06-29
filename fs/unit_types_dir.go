package fs

import (
	"context"
	"os"
	"syscall"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/awfulcooking/journalfs/journalcache"
)

// UnitTypesDir is a filesystem node that lists directories
// which correspond to systemd unit types for which there
// may be journal entries.
//
// These directories in turn list log files which map to journal
// entries for units of their given unit type (see UnitLogsDir)
//
// e.g  service/ssh.log  => "ssh.service" logs
// e.g  timer/some.log   => "some.timer"  logs
type UnitTypesDir struct {
	journalCache *journalcache.JournalCache

	typeDirs map[string]*UnitLogsDir
}

var _ fs.Node = (*UnitTypesDir)(nil)
var _ fs.HandleReadDirAller = (*UnitTypesDir)(nil)
var _ fs.NodeStringLookuper = (*UnitTypesDir)(nil)

// Attr populates filesystem metadata for a UnitTypesDir
func (d *UnitTypesDir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Valid = 0

	attr.Uid = UID
	attr.Gid = GID
	attr.Mode = os.ModeDir | 0o550
	attr.Size = uint64(len(d.typeDirs))
	attr.Ctime = time.Now()
	attr.Mtime = time.Now()

	return nil
}

// ReadDirAll enumerates directory entries for a UnitTypesDir
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

// Lookup looks up a node by name within a UnitTypesDir
func (d *UnitTypesDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if node, ok := d.typeDirs[name]; ok {
		return node, nil
	}

	return nil, syscall.ENOENT
}

// NewUnitTypesDir returns a new UnitTypesDir
//
// These directories in turn list log files which map to journal
// entries for units of their given unit type (see UnitLogsDir)
//
// e.g services/ssh.log corresponding to "ssh.service" logs
func NewUnitTypesDir(journalCache *journalcache.JournalCache) *UnitTypesDir {
	return &UnitTypesDir{
		journalCache: journalCache,
		typeDirs:     makeUnitTypeDirs(journalCache),
	}
}

// UNIT_TYPES represents the list of systemd unit types
// for which there can be journal entries
//
// This forms the list of directories under a UnitTypesDir
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
