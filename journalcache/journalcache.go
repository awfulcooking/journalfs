package journalcache

import (
	"errors"
	"log"

	"github.com/coreos/go-systemd/sdjournal"
)

// todo: locking

type JournalCache struct {
	journal *sdjournal.Journal

	entries       []*JournalEntry
	entriesByUnit map[string][]*JournalEntry

	following bool
	Debug bool
}

func (jc *JournalCache) Reset() {
	jc.entries = make([]*JournalEntry, 0)
	jc.entriesByUnit = make(map[string][]*JournalEntry)
}

func (jc *JournalCache) Entries() []*JournalEntry {
	return jc.entries
}

func (jc *JournalCache) EntriesByUnit(unit string) []*JournalEntry {
	return jc.entriesByUnit[unit]
}

func (jc *JournalCache) UnitNames() []string {
	names := make([]string, len(jc.entriesByUnit))
	for name, _ := range jc.entriesByUnit {
		names = append(names, name)
	}
	return names
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

func (jc *JournalCache) Follow() error {
	if jc.Debug {
		log.Println("Follow()")
	}
	if jc.following {
		return errors.New("already following")
	}
	jc.following = true
	go jc.follow()
	return nil
}

func (jc *JournalCache) follow() {
	for {
		jc.journal.Wait(sdjournal.IndefiniteWait)

		if count, err := jc.journal.Next(); err != nil {
			panic(err)
		} else if count == 1 {
			entry, err := jc.journal.GetEntry()
			if err != nil {
				panic(err)
			}

			jc.addEntry(entry)
		}
	}
}

func (jc *JournalCache) addEntry(sdjournalEntry *sdjournal.JournalEntry) {
	entry := makeJournalEntry(sdjournalEntry)

	if jc.Debug {
		log.Println("New entry", entry)
	}

	jc.entries = append(jc.entries, &entry)

	if entry.Unit != "" {
		if jc.entriesByUnit[entry.Unit] == nil {
			jc.entriesByUnit[entry.Unit] = make([]*JournalEntry, 0)
		}
		jc.entriesByUnit[entry.Unit] = append(jc.entriesByUnit[entry.Unit], &entry)
	}
}

func NewJournalCache() (*JournalCache, error) {
	var journal *sdjournal.Journal
	var err error

	if journal, err = sdjournal.NewJournal(); err == nil {
		cache := &JournalCache{
			journal:       journal,
			entries:       make([]*JournalEntry, 0, 200),
			entriesByUnit: make(map[string][]*JournalEntry),
		}
		return cache, nil
	}

	return nil, err
}
