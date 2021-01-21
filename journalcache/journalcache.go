package journalcache

import (
	"github.com/coreos/go-systemd/sdjournal"
)

// todo: locking

type JournalCache struct {
	journal *sdjournal.Journal

	entries       []*sdjournal.JournalEntry
	entriesByUnit map[string][]*sdjournal.JournalEntry
}

func (jc *JournalCache) Reset() {
	jc.entries = make([]*sdjournal.JournalEntry, 0)
	jc.entriesByUnit = make(map[string][]*sdjournal.JournalEntry)
}

func (jc *JournalCache) Entries() []*sdjournal.JournalEntry {
	return jc.entries
}

func (jc *JournalCache) EntriesByUnit() map[string][]*sdjournal.JournalEntry {
	return jc.entriesByUnit
}

func (jc *JournalCache) MatchUnit(unit string) error {
	match := sdjournal.Match{
		Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
		Value: unit,
	}
	return jc.journal.AddMatch(match.String())
}

func (jc *JournalCache) FlushMatches() {
	jc.journal.FlushMatches()
}

func (jc *JournalCache) Load() (int, error) {
	jc.Reset()

	jc.journal.SeekHead()

	for {
		count, err := jc.journal.Next()
		if count == 0 || err != nil {
			return len(jc.entries), err
		}

		entry, err := jc.journal.GetEntry()
		if err != nil {
			return 0, err
		}

		jc.addEntry(entry)
	}
}

func (jc *JournalCache) addEntry(entry *sdjournal.JournalEntry) {
	jc.entries = append(jc.entries, entry)

	if unit, ok := entry.Fields[sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT]; ok {
		if jc.entriesByUnit[unit] == nil {
			jc.entriesByUnit[unit] = make([]*sdjournal.JournalEntry, 0)
		}

		jc.entriesByUnit[unit] = append(jc.entriesByUnit[unit], entry)
	}
}

func NewJournalCache() (*JournalCache, error) {
	var journal *sdjournal.Journal
	var err error

	if journal, err = sdjournal.NewJournal(); err == nil {
		cache := &JournalCache{
			journal: journal,
			entries: make([]*sdjournal.JournalEntry, 0, 200),
		}
		return cache, nil
	}

	return nil, err
}
