package journalcache

import (
	"time"

	"github.com/coreos/go-systemd/v22/sdjournal"
)

type JournalEntry struct {
	Message   string
	Unit      string
	UserUnit  string
	Timestamp time.Time
}

func makeJournalEntry(entry *sdjournal.JournalEntry) JournalEntry {
	return JournalEntry{
		Message:   entry.Fields[sdjournal.SD_JOURNAL_FIELD_MESSAGE],
		Unit:      entry.Fields[sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT],
		UserUnit:  entry.Fields[sdjournal.SD_JOURNAL_FIELD_SYSTEMD_USER_UNIT],
		Timestamp: time.Unix(int64(entry.RealtimeTimestamp/uint64(time.Millisecond)), 0),
	}
}
