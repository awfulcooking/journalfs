package fs

import (
	"context"
	"fmt"
	"strings"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/coreos/go-systemd/sdjournal"
)

type FileUnitLog struct {
	unit string

	entries []*sdjournal.JournalEntry
}

var _ fs.Node = (*FileUnitLog)(nil)
var _ fs.HandleReadAller = (*FileUnitLog)(nil)

func (f *FileUnitLog) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = 1
	attr.Mode = 0o444
	attr.Size = uint64(len(f.data()))

	return nil
}

func (f *FileUnitLog) ReadAll(ctx context.Context) ([]byte, error) {
	fmt.Println("ReadAll()")
	return f.data(), nil
}

func (f *FileUnitLog) data() []byte {
	var log []string

	for _, entry := range f.entries {
		log = append(log, fmt.Sprintf(
			"%s",
			entry.Fields[sdjournal.SD_JOURNAL_FIELD_MESSAGE],
		))
	}

	return []byte(strings.Join(log, "\n") + "\n")
}
