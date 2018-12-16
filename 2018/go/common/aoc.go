package common

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

// aoc defaults
const aocFirstYear = 2015
const defaultAocBaseURL = "https://adventofcode.com"
const inputEndpoint = "%d/day/%d/input" // [year]/day/[day]/input

// CommandPath *kludge ahead* gets the directory that called the command command directory
// used to support "go run" invocations -- obviously another solution is required
// if we need to support "built" environments.
//
// If there is a more consistent way to do this, I haven't been able to figure it out.
// 2, matches being called directly from main > AOCInputFile > CommandPath
func CommandPath() (string, error) {
	if _, file, _, ok := runtime.Caller(2); ok {
		return path.Dir(file), nil
	}

	return "", fmt.Errorf("could not determine command path")
}

// fetchInput will retrieve the year and day's input if it doesn't exist
func fetchInput(url *url.URL, sessionID string, outfile string, year int, day int) error {
	// now fetch
	endpointURL, error := url.Parse(fmt.Sprintf(inputEndpoint, year, day))
	if error == nil {
		return FetchFile(aocClient(url, sessionID), endpointURL, outfile)
	}

	return error
}

func aocClient(url *url.URL, sessionID string) *http.Client {
	return &http.Client{
		Timeout: time.Second * 30,
		Jar: NewCookieJar(url, map[string]string{
			"session": sessionID,
		}),
	}
}

// aocBaseURL
func aocBaseURL(fallback string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Base URL: [%s]", fallback)
	if scanner.Scan() {
		if ans := scanner.Text(); len(ans) > 0 {
			return ans
		}
	}

	return fallback
}

// AOCProblem wraps data related to an AOC problem
type AOCProblem struct {
	Language  string
	Year      int
	Day       int
	Part      int
	InputFile string
}

// IsValid simple validation of AOC meta data
func (p AOCProblem) IsValid() error {
	if !(p.Day > 0 && p.Day < 26) {
		return fmt.Errorf("day out of range: %d", p.Day)
	}

	if p.Year < aocFirstYear {
		return fmt.Errorf("year invalid: %d", p.Year)
	}

	return nil
}

// 2018 uses 1 input file per day
func (p AOCProblem) fetchInput(serverURL *url.URL) error {
	destDirectory := filepath.Dir(p.InputFile)
	if error := os.MkdirAll(destDirectory, os.ModePerm); error != nil {
		return error // MkdirAll error
	}

	log.Println("using url:", serverURL)
	error := fetchInput(serverURL, aocSessionID(), p.InputFile, p.Year, p.Day)
	if error != nil {
		return error // fetch error
	}

	log.Println("input created:", p.InputFile)
	return nil
}

// OpenInput convenience for fetching and opening the problem input
func (p AOCProblem) OpenInput() (*os.File, error) {
	// verify the path doesn't already exist before fetching
	_, error := os.Stat(p.InputFile)
	if os.IsNotExist(error) {
		baseurl := aocBaseURL(defaultAocBaseURL)
		parsedurl, error := url.Parse(baseurl)
		if error != nil {
			return nil, error // url parse error
		}

		if error = p.fetchInput(parsedurl); error != nil {
			return nil, error
		}
	} else {
		log.Println("skipped fetch, input exists:", p.InputFile)
	}

	return os.Open(p.InputFile) // implicitly returns os.Open args
}

// AOCProblemFromPath attempts to determine [year]/[lang]/[day]/part[part] from path
// typically p should be the "module" or "caller" path
func AOCProblemFromPath(p string) *AOCProblem {
	lang, day, part := "", -1, -1

	// climb up from p
	for prev, cur := "", p; cur != prev; prev, cur = cur, path.Dir(cur) {
		// check if this level has a "year" folder
		curbase := path.Base(cur)

		if part > 0 {
			if day > 0 {
				if lang != "" {
					year := (-1)
					rslt, error := strconv.Atoi(curbase)
					if error == nil {
						year = rslt
					} else {
						log.Println("invalid year", error)
						// make attempt to "guess" remaining problem details
						now := time.Now()
						year = now.Year()
						if now.Day() > day {
							year--
						}
					}

					inputfile := fmt.Sprintf("%s/inputs/%02d/input.txt", cur, day)
					return &AOCProblem{lang, year, day, part, inputfile}
				} else {
					// assume current node is the language
					lang = curbase
				}
			} else {
				rslt, error := strconv.Atoi(curbase)
				if error == nil {
					day = rslt
				} else {
					log.Println("day error", error)
					return nil
				}
			}
		} else {
			partRe := regexp.MustCompile(`part(\d+)`)
			fields := partRe.FindAllStringSubmatch(curbase, -1)
			if fields != nil {
				rslt, error := strconv.Atoi(fields[0][1])
				if error == nil {
					part = rslt
				} else {
					log.Println("part error", error)
					return nil
				}
			}
		}
	}

	return nil
}

func aocSessionID() string {
	const sessionIDEnvKey = "AOC_SESSIONID"
	sessionID := os.Getenv(sessionIDEnvKey)
	if sessionID == "" {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("Session ID (%s not set): ", sessionIDEnvKey)
		scanner.Scan()
		sessionID = scanner.Text()
	}

	return sessionID
}

// GetAOCProblem returns a struct with problem vitals
func GetAOCProblem() (*AOCProblem, error) {
	cmdpath, error := CommandPath()
	if error != nil {
		return nil, error
	}

	vitals := AOCProblemFromPath(cmdpath)
	if vitals == nil {
		return nil, fmt.Errorf("invalid problem")
	}

	if error = vitals.IsValid(); error != nil {
		return nil, fmt.Errorf("could not determine problem meta: %+v", vitals)
	}

	return vitals, nil
}
