package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/coreos/go-systemd/sdjournal"
)

func main() {
	unit := flag.String("unit", "", "systemd unit name")
	flag.Parse()

	j, err := sdjournal.NewJournal()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if *unit != "" {
		match := sdjournal.Match{
			Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
			Value: *unit,
		}

		j.AddMatch(match.String())
	}

	j.SeekTail()
	j.Previous()

	for {
		entry, _ := j.GetEntry()

		// printJSON(entry.Fields)

		fmt.Printf("[%s] %s\n", entry.Fields["SYSLOG_IDENTIFIER"], entry.Fields["MESSAGE"])

		size, _ := j.Next()
		if size == 0 {
			j.Wait(sdjournal.IndefiniteWait)
		}
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
