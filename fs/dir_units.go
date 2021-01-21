package fs

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"bazil.org/fuse"
	bzfs "bazil.org/fuse/fs"

	"github.com/togetherbeer/journalfs/journalcache"
)

type DirUnits struct {
	journalCache *journalcache.JournalCache

	unitFiles map[string]bzfs.Node
}

var _ bzfs.HandleReadDirAller = (*DirUnits)(nil)

func NewDirUnits(jc *journalcache.JournalCache) *DirUnits {
	return &DirUnits{
		journalCache: jc,
		unitFiles:    make(map[string]bzfs.Node),
	}
}

func (du *DirUnits) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = 1
	attr.Mode = os.ModeDir | 0o555
	attr.Size = 123456
	return nil
}

func (du *DirUnits) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var response []fuse.Dirent

	for unit, _ := range du.journalCache.EntriesByUnit() {
		fmt.Printf("Got unit: %v\n", unit)
		response = append(response, fuse.Dirent{
			Name: unit,
			Type: fuse.DT_File,
		})
	}
	return response, nil
}

func (du *DirUnits) Lookup(ctx context.Context, name string) (bzfs.Node, error) {
	if unitFile := du.unitFile(name); unitFile != nil {
		return unitFile, nil
	} else {
		return nil, syscall.ENOENT
	}
}

func (du *DirUnits) unitFile(name string) bzfs.Node {
	if node, ok := du.unitFiles[name]; ok {
		return node
	} else if entries, ok := du.journalCache.EntriesByUnit()[name]; ok {
		du.unitFiles[name] = &FileUnitLog{
			unit:    name,
			entries: entries,
		}

		return du.unitFiles[name]
	}

	return nil
}
