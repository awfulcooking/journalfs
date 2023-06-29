package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/awfulcooking/journalfs/journalcache"
)

var unit = flag.String("unit", "", "systemd unit name")
var wait = flag.Bool("f", false, "tail new messages")
var printCount = flag.Bool("c", false, "print count of entries read")

func init() {
	flag.Parse()
}

func main() {
	jc, err := journalcache.NewJournalCache()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if *unit != "" {
		jc.MatchUnit(*unit)
	}

	count, err := jc.Load()
	must(err)

	for _, entry := range jc.Entries() {
		fmt.Println(entry.Fields["MESSAGE"])
	}

	if *printCount {
		fmt.Printf("Loaded %d entries. Err: %v\n", count, err)
	}
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
