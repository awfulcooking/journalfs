package mount

import (
	"bazil.org/fuse"
	bzfs "bazil.org/fuse/fs"
	"github.com/togetherbeer/journalfs/fs"
)

type Mount struct {
	dir string
	fs  *fs.FS
}

func (m *Mount) Dir() string {
	return m.dir
}

func (m *Mount) Serve() error {
	conn, err := fuse.Mount(
		m.dir,
		fuse.FSName("journalfs"),
		fuse.Subtype("journalfs"),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	return bzfs.Serve(conn, m.fs)
}

func NewMount(dir string) *Mount {
	return &Mount{
		dir: dir,
		fs:  fs.NewFS(),
	}
}
