package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/coreos/go-systemd/sdjournal"
)

var unit = flag.String("unit", "", "systemd unit name")
var numPrevious = flag.Uint64("n", 10, "number of previous journal messages to display")

func init() {
	flag.Parse()
}

func main() {
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
	j.PreviousSkip(*numPrevious + 1)

	for {
		count, err := j.Next()
		must(err)

		if count == 1 {
			entry, err := j.GetEntry()
			must(err)

			fmt.Printf(
				"[%s] %s\n",
				entry.Fields[sdjournal.SD_JOURNAL_FIELD_SYSLOG_IDENTIFIER],
				entry.Fields[sdjournal.SD_JOURNAL_FIELD_MESSAGE],
			)
		} else {
			j.Wait(sdjournal.IndefiniteWait)
		}
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
