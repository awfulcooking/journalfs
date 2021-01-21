package fs

import (
	"context"
	"os"

	"bazil.org/fuse"
	bzfs "bazil.org/fuse/fs"
)

type Root struct {
	bzfs.Node
}

func (r *Root) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = 1
	attr.Mode = os.ModeDir | 0o555
	return nil
}

func NewRoot() *Root {
	return &Root{}
}
