package mount

import "bazil.org/fuse"

type MountOption func() fuse.MountOption

// AllowOther allows other users to access the filesystem.
func AllowOther() fuse.MountOption {
	return fuse.AllowOther()
}

// Name sets the name of the filesystem
func Name(name string) fuse.MountOption {
	return fuse.FSName(name)
}

func fuseMountOptions(options []MountOption) []fuse.MountOption {
	var fuseOptions []fuse.MountOption
	for _, opt := range options {
		fuseOptions = append(fuseOptions, opt())
	}
	return fuseOptions
}
