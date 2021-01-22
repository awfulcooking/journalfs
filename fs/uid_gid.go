package fs

import "os"

var UID = uint32(os.Getuid())
var GID = uint32(os.Getgid())
