package fs

import (
	"context"
	"os"
	"strings"
	"syscall"

	"bazil.org/fuse"
	bzfs "bazil.org/fuse/fs"

	"github.com/awfulcooking/journalfs/journalcache"
)

type UnitLogsDir struct {
	Type string

	journalCache *journalcache.JournalCache

	logFiles map[string]bzfs.Node
}

var _ bzfs.Node = (*UnitLogsDir)(nil)
var _ bzfs.HandleReadDirAller = (*UnitLogsDir)(nil)
var _ bzfs.NodeStringLookuper = (*UnitLogsDir)(nil)

func NewUnitLogsDir(jc *journalcache.JournalCache, unitType string) *UnitLogsDir {
	return &UnitLogsDir{
		Type: unitType,

		journalCache: jc,
		logFiles:     make(map[string]bzfs.Node),
	}
}

func (d *UnitLogsDir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Uid = UID
	attr.Gid = GID
	attr.Mode = os.ModeDir | 0o550
	attr.Size = uint64(len(d.matchingUnitNames()))

	return nil
}

func (d *UnitLogsDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var response []fuse.Dirent

	for _, name := range d.matchingUnitNames() {
		response = append(response, fuse.Dirent{
			Name: strings.TrimSuffix(name, "."+d.Type) + ".log",
			Type: fuse.DT_File,
		})
	}

	return response, nil
}

func (d *UnitLogsDir) matchingUnitNames() []string {
	var names []string
	for _, name := range d.journalCache.UnitNames() {
		if strings.HasSuffix(name, "."+d.Type) {
			names = append(names, name)
		}
	}
	return names
}

func (d *UnitLogsDir) Lookup(ctx context.Context, name string) (bzfs.Node, error) {
	nameWithoutExtension := strings.TrimSuffix(name, ".log")

	if unitFile := d.unitFile(nameWithoutExtension + "." + d.Type); unitFile != nil {
		return unitFile, nil
	} else {
		return nil, syscall.ENOENT
	}
}

func (d *UnitLogsDir) unitFile(name string) bzfs.Node {
	if node, ok := d.logFiles[name]; ok {
		return node
	} else if entries := d.journalCache.EntriesByUnit(name); len(entries) > 0 {
		d.logFiles[name] = NewUnitLogFile(d.journalCache, name)
		return d.logFiles[name]
	}

	return nil
}
