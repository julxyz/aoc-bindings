package aoclib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type AoCHelper struct {
	sessionCookie string
	httpClient    *http.Client
	cacheDir      string
}

// oh noes something went wrong
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// NewAoCHelper returns a new AoCHelper initialized with the session cookie.
func NewAoCHelper() AoCHelper {
	tmpdir := os.Getenv("TMP")
	if tmpdir == "" {
		tmpdir = "/tmp"
	}
	return AoCHelper{getCookie(), &http.Client{}, filepath.Join(tmpdir, "aoc_cache")}
}

// getCookie returns the session cookie if available
func getCookie() string {
	cookie := os.Getenv("AOC_Cookie")
	if cookie == "" {
		log.Fatalf("No session cookie set.\nPlease set your advent of code session cookie using:\n\nexport AOC_COOKIE = {your session token}")
	}
	return cookie
}

// read stored input from cache
func (h *AoCHelper) readCached(year int, day int) []string {
	readFile, err := os.Open(filepath.Join(h.cacheDir, h.getFileName(year, day)))
	if err != nil {
		// Oh no dangerous recursion watch out this will totally never get stuck here HAHA
		h.cacheFromWeb(year, day)
		return h.readCached(year, day)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	readFile.Close()
	return lines
}

// get puzzle input from interwebz and store to cache
func (h *AoCHelper) cacheFromWeb(year int, day int) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day), nil)
	checkErr(err)
	req.Header.Add("Cookie", "session="+h.sessionCookie)
	resp, err := h.httpClient.Do(req)
	checkErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	err = os.WriteFile(filepath.Join(h.cacheDir, h.getFileName(year, day)), body, os.ModePerm)
}

// get filename for a given day
func (h *AoCHelper) getFileName(year int, day int) string {
	return fmt.Sprintf("%d%20d.txt", year, day)
}

// GetInput retrieves the puzzle input for a certain day as a []string.
//
// Also caches the input, use force=true to force overwrite the cache. Please use caching to avoid unnecessary to the AoC servers
func (h *AoCHelper) GetInput(year int, day int, force bool) []string {
	if force {
		h.cacheFromWeb(year, day)
		return h.readCached(year, day)
	}
	return h.readCached(year, day)
}
