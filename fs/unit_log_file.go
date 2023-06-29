package fs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"github.com/awfulcooking/journalfs/journalcache"
)

type UnitLogFile struct {
	Unit         string
	journalCache *journalcache.JournalCache
}

var _ fs.Node = (*UnitLogFile)(nil)
var _ fs.HandleReader = (*UnitLogFile)(nil)

func NewUnitLogFile(jc *journalcache.JournalCache, unit string) *UnitLogFile {
	return &UnitLogFile{
		journalCache: jc,
		Unit:         unit,
	}
}

func (f *UnitLogFile) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	resp.Flags |= fuse.OpenDirectIO
	return f, nil
}

func (f *UnitLogFile) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Valid = 0

	attr.Uid = UID
	attr.Gid = GID
	attr.Mode = 0o440
	attr.Size = uint64(len(f.data()))
	attr.Mtime = f.modifiedTime()
	attr.Atime = time.Now()
	attr.Ctime = time.Now()

	if f.journalCache.Debug {
		log.Println(f, "Attr", "|", "modified", attr.Mtime, "size", attr.Size)
	}

	return nil
}

func (f *UnitLogFile) entries() []*journalcache.JournalEntry {
	return f.journalCache.EntriesByUnit(f.Unit)
}

func (f *UnitLogFile) modifiedTime() time.Time {
	entries := f.entries()
	if n := len(entries); n > 0 {
		return f.entries()[len(entries)-1].Timestamp
	}
	return time.Unix(0, 0)
}

func (f *UnitLogFile) data() []byte {
	var log []string

	for _, entry := range f.entries() {
		log = append(log, fmt.Sprintf(
			"%s",
			entry.Message,
		))
	}

	return []byte(strings.Join(log, "\n") + "\n")
}

var _ fs.HandleReader = (*UnitLogFile)(nil)

func (f *UnitLogFile) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	if f.journalCache.Debug {
		log.Println(f, "Read", "|", req)
	}
	fuseutil.HandleRead(req, resp, f.data())
	return nil
}

func (f *UnitLogFile) String() string {
	return fmt.Sprintf("%s", f.Unit)
}
