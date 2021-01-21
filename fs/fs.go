package fs

import (
	bzfs "bazil.org/fuse/fs"
)

type FS struct {
	bzfs.FS // interface
}

func (fs *FS) Root() (bzfs.Node, error) {
	return NewRoot(), nil
}

func NewFS() *FS {
	return &FS{}
}
