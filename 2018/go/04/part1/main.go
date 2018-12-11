package main

import (
	"aoc/common"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const day = 4
const part = 1

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	if file, error := common.AOCInputFile(day); error == nil {
		defer file.Close()
		records := processLogEntries(loadEntries(file))
		sleepRecords := []sleepRecord{}
		for _, r := range records {
			sleepRecords = append(sleepRecords, r.sleepRecord())
		}

		sort.Slice(sleepRecords, func(i, j int) bool {
			return sleepRecords[i].totalMinutesAsleep < sleepRecords[j].totalMinutesAsleep
		})

		sleeprec := sleepRecords[len(sleepRecords)-1]
		log.Println("sleepiest guard:", sleeprec.guardID)
		log.Println("total minutes asleep:", sleeprec.totalMinutesAsleep)
		log.Println("sleepiest minute:", sleeprec.sleepiestMinute)
		log.Println("aoc answer =", int(sleeprec.guardID)*sleeprec.sleepiestMinute)
	} else {
		log.Println(error)
	}
}

var defaultLocation *time.Location

func init() {
	l, error := time.LoadLocation("")
	if error == nil {
		defaultLocation = l
	} else {
		log.Println(error)
	}
}

// convenience for eating parse errors
func safeParseInt(s string, fallback int) int {
	rslt, error := strconv.Atoi(s)
	if error == nil {
		return rslt
	}

	return fallback
}

func timeFromComponents(components ...string) time.Time {
	return time.Date(
		safeParseInt(components[0], 0),
		time.Month(safeParseInt(components[1], 0)),
		safeParseInt(components[2], 0),
		safeParseInt(components[3], 0),
		safeParseInt(components[4], 0),
		0, 0,
		defaultLocation,
	)
}

type guardID int
type logEntryType int

const (
	beginShift logEntryType = iota
	startSleep
	endSleep
)

func logEntryTypeFromString(s string) (logEntryType, error) {
	if strings.EqualFold(s, "begins shift") {
		return beginShift, nil
	} else if strings.EqualFold(s, "falls asleep") {
		return startSleep, nil
	} else if strings.EqualFold(s, "wakes up") {
		return endSleep, nil
	}

	return beginShift, fmt.Errorf("invalid shift event type %s", s)
}

func (t logEntryType) String() string {
	return []string{"begins shift", "falls asleep", "wakes up"}[t]
}

type logEntryInterface interface {
	Timestamp() time.Time
	EntryType() logEntryType
}

type logEntry struct {
	logEntryInterface
	timestamp time.Time
	entryType logEntryType
}

func (e logEntry) Timestamp() time.Time {
	return e.timestamp
}

func (e logEntry) EntryType() logEntryType {
	return e.entryType
}

type guardShiftStartEntry struct {
	logEntry
	guardID guardID
}

func newLogEntry(timestamp time.Time, entryType logEntryType) logEntry {
	return logEntry{timestamp: timestamp, entryType: entryType}
}

func newGuardShiftStartEntry(timestamp time.Time, entryType logEntryType, guardID guardID) guardShiftStartEntry {
	return guardShiftStartEntry{
		newLogEntry(timestamp, entryType),
		guardID,
	}
}

func loadEntries(f *os.File) []logEntryInterface {
	const entryFormat = `\[((\d+)-(\d+)-(\d+) (\d+):(\d+))]\s*(Guard #(\d+)\s*)?(.*)`
	entryRe := regexp.MustCompile(entryFormat)
	entries := []logEntryInterface{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := entryRe.FindAllStringSubmatch(line, -1)
		if fields != nil {
			var entry logEntryInterface
			t := timeFromComponents(fields[0][2], fields[0][3], fields[0][4], fields[0][5], fields[0][6])
			if fields[0][7] == "" {
				entryType, error := logEntryTypeFromString(fields[0][9])
				if error == nil {
					entry = newLogEntry(t, entryType)
				} else {
					log.Println(error)
				}
			} else {
				guardid := safeParseInt(fields[0][8], 0)
				entry = newGuardShiftStartEntry(t, beginShift, guardID(guardid))
			}

			entries = append(entries, entry)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp().Before(entries[j].Timestamp())
	})

	return entries
}

type guardAction struct {
	timestamp  time.Time
	actionType logEntryType
}

type guardShift struct {
	start   time.Time
	actions []guardAction
}

type guardRecord struct {
	guardID guardID
	shifts  []*guardShift
}

type sleepRecord struct {
	guardID            guardID
	minutes            [60]int
	totalMinutesAsleep int
	sleepiestMinute    int
}

func (r guardRecord) sleepRecord() sleepRecord {
	minutes := [60]int{}

	var sstart int
	for _, s := range r.shifts {
		sstart = 0
		for _, a := range s.actions {
			if a.actionType == startSleep {
				sstart = a.timestamp.Minute()
			} else if a.actionType == endSleep {
				for m := sstart; m < a.timestamp.Minute(); m++ {
					minutes[m]++
				}
			}
		}
	}

	// update stats
	mx, minute, sum := 0, 0, 0
	for i, c := range minutes {
		sum += c
		if c > mx {
			minute = i
			mx = c
		}
	}

	sr := sleepRecord{r.guardID, minutes, sum, minute}
	return sr
}

func processLogEntries(entries []logEntryInterface) map[guardID]*guardRecord {
	records := map[guardID]*guardRecord{}

	var curshift *guardShift
	for _, entry := range entries {
		if guardstart, ok := entry.(guardShiftStartEntry); ok {
			// new shift
			gid := guardstart.guardID
			rec, ok := records[gid]
			if !ok {
				rec = &guardRecord{gid, []*guardShift{}}
				records[gid] = rec
			}

			curshift = &guardShift{entry.Timestamp(), []guardAction{}}
			rec.shifts = append(rec.shifts, curshift)
		} else {
			curshift.actions = append(curshift.actions, guardAction{entry.Timestamp(), entry.EntryType()})
		}
	}

	return records
}
