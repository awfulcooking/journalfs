package mount

import (
	"bazil.org/fuse"
	bzfs "bazil.org/fuse/fs"

	"github.com/awfulcooking/journalfs/fs"
	"github.com/awfulcooking/journalfs/journalcache"
)

type Mount struct {
	dir string
	fs  *fs.FS

	fuseConn *fuse.Conn
}

func (m *Mount) Dir() string {
	return m.dir
}

func (m *Mount) Serve(options ...MountOption) error {
	fuseOptions := []fuse.MountOption{
		fuse.FSName("journalfs"),
		fuse.Subtype("journalfs"),
		fuse.ReadOnly(),
		fuse.DefaultPermissions(), // ask kernel to perform file-mode based access control
	}

	fuseOptions = append(fuseOptions, fuseMountOptions(options)...)

	conn, err := fuse.Mount(m.dir, fuseOptions...)

	if err != nil {
		return err
	}

	m.fuseConn = conn
	defer conn.Close()

	return bzfs.Serve(conn, m.fs)
}

func (m *Mount) Unmount() error {
	return fuse.Unmount(m.dir)
}

func NewMount(dir string, journalCache *journalcache.JournalCache) *Mount {
	return &Mount{
		dir: dir,
		fs:  fs.NewFS(journalCache),
	}
}
