package common

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func makeCookieJar(url *url.URL, cookies map[string]string) *cookiejar.Jar {
	jar, error := cookiejar.New(nil)
	if error == nil {
		jar.SetCookies(url, func() []*http.Cookie {
			collection := []*http.Cookie{}
			for k, v := range cookies {
				collection = append(collection, &http.Cookie{Name: k, Value: v})
			}
			return collection
		}())

		return jar
	}

	log.Println(error)
	return nil
}

func aocClient(url *url.URL, sessionID string) *http.Client {
	return &http.Client{
		Timeout: time.Second * 30,
		Jar: makeCookieJar(url, map[string]string{
			"session": sessionID,
		}),
	}
}

// fetchInput will retrieve the year and day's input if it doesn't exist
func fetchInput(url *url.URL, sessionID string, outfile string, year int, day int) error {
	// fetch inputs up to current date
	const inputEndpoint = "%d/day/%d/input"

	// now fetch
	endpointURL, error := url.Parse(fmt.Sprintf(inputEndpoint, year, day))
	if error == nil {
		return fetchFile(aocClient(url, sessionID), endpointURL, outfile)
	}

	return error
}

func fetchFile(client *http.Client, url *url.URL, dest string) error {
	resp, error := client.Get(url.String())
	if error == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			out, error := os.Create(dest)
			defer out.Close()
			if error == nil {
				_, error := io.Copy(out, resp.Body)
				if error == nil {
					return nil // success
				}

				return error
			}

			return error
		}

		return fmt.Errorf("status: %s", resp.Status)
	}

	return error
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

// findAOCRoot makes an attempt to find the containing directory for aoc
// that matches [aoc]/[year]/[lang]/<days>
func findAOCRoot(start string, year int) string {
	yearstr := strconv.Itoa(year)
	for prev, cur := "", start; cur != prev; prev, cur = cur, path.Dir(cur) {
		// check if this level has a "year" folder
		curbase := path.Base(cur)
		if strings.EqualFold(curbase, "aoc") {
			return cur
		}

		if files, error := ioutil.ReadDir(cur); error == nil {
			for _, info := range files {
				if info.IsDir() {
					n := info.Name()
					if strings.EqualFold(n, yearstr) {
						return path.Join(cur, n)
					}
				}
			}
		}
	}

	return ""
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

const aocFirstYear = 2015
const defaultAocBaseURL = "https://adventofcode.com"

// AOCInputFile returns "day" input file
// 2018 uses 1 input file per day
func AOCInputFile(day int) (*os.File, error) {
	now := time.Now()

	if !(day > 0 && day < 26) {
		log.Printf("day out of range: %d\n", day)
		os.Exit(2)
	}

	year := now.Year()
	if now.Month() < 12 {
		year--
	}

	if year < aocFirstYear {
		log.Printf("year invalid: %d\n", year)
	}

	cmdpath, error := CommandPath()
	if error == nil {
		aocroot := findAOCRoot(cmdpath, year)
		log.Println("aoc root:", aocroot)
		destFile := fmt.Sprintf("%s/inputs/%02d/input.txt", aocroot, day)

		// verify the path doesn't already exist before fetching
		_, error := os.Stat(destFile)
		if os.IsNotExist(error) {
			destDirectory := filepath.Dir(destFile)
			if error := os.MkdirAll(destDirectory, os.ModePerm); error == nil {
				aocurl := aocBaseURL(defaultAocBaseURL)
				log.Println("using:", aocurl)
				if parsedurl, error := url.Parse(aocurl); error == nil {
					error := fetchInput(parsedurl, aocSessionID(), destFile, year, day)
					if error == nil {
						// successful branch
						log.Println("input created:", destFile)
					} else {
						return nil, error // fetch error
					}
				} else {
					return nil, error // url parse error
				}
			} else {
				return nil, error // MkdirAll error
			}
		} else {
			// successful branch
			log.Println("skipped fetch, input exists:", destFile)
		}

		// return an opened file
		return os.Open(destFile)
	}

	return nil, error
}
