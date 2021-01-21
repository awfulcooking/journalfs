package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/togetherbeer/journalfs/mount"
)

var mountPath = flag.String("p", "./journal", "mount path")

func init() {
	flag.Parse()
}

func main() {
	mount := mount.NewMount(*mountPath)
	must(mount.Serve())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func printJSON(v interface{}) error {
	if s, err := toJSONString(v); err != nil {
		return err
	} else {
		fmt.Println(s)
		return nil
	}
}

func toJSONString(v interface{}) (string, error) {
	if bytes, err := json.MarshalIndent(v, "", "  "); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}
