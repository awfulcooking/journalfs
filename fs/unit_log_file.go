package fs

import (
	"context"
	"fmt"
	"strings"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/coreos/go-systemd/sdjournal"
)

type UnitLogFile struct {
	unit string

	entries []*sdjournal.JournalEntry
}

var _ fs.Node = (*UnitLogFile)(nil)
var _ fs.HandleReadAller = (*UnitLogFile)(nil)

func (f *UnitLogFile) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = 1
	attr.Uid = UID
	attr.Gid = GID
	attr.Mode = 0o440
	attr.Size = uint64(len(f.data()))

	return nil
}

func (f *UnitLogFile) ReadAll(ctx context.Context) ([]byte, error) {
	return f.data(), nil
}

func (f *UnitLogFile) data() []byte {
	var log []string

	for _, entry := range f.entries {
		log = append(log, fmt.Sprintf(
			"%s",
			entry.Fields[sdjournal.SD_JOURNAL_FIELD_MESSAGE],
		))
	}

	return []byte(strings.Join(log, "\n") + "\n")
}
