package common

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

// NewCookieJar is a convenience for creating a jar from key-value pairs
func NewCookieJar(url *url.URL, cookies map[string]string) *cookiejar.Jar {
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

// FetchFile convenience for downloading a file to <dest>
func FetchFile(client *http.Client, url *url.URL, dest string) error {
	resp, error := client.Get(url.String())
	defer resp.Body.Close()
	if error != nil {
		return error
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %s", resp.Status)
	}

	out, error := os.Create(dest)
	defer out.Close()
	if error != nil {
		return error
	}

	_, error = io.Copy(out, resp.Body)
	if error != nil {
		return error
	}

	return nil // success
}
