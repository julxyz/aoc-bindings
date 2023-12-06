package aoclib

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type AoCHelper struct {
	Year          int
	Day           int
	puzzleStatus  int
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
func NewAoCHelper(year int, day int) AoCHelper {
	tmpdir := os.Getenv("TMP")
	if tmpdir == "" {
		tmpdir = "/tmp"
	}
	// create cache folder
	dir := filepath.Join(tmpdir, "aoc_cache")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	h := AoCHelper{year, day, 0, getCookie(), &http.Client{}, dir}
	h.puzzleStatus = h.getPuzzleStatus()
	h.printPuzzleStatus()
	return h
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
func (h *AoCHelper) readCached() []string {
	readFile, err := os.Open(filepath.Join(h.cacheDir, h.getFileName()))
	if err != nil {
		// Oh no dangerous recursion watch out this will totally never get stuck here HAHA
		h.cacheFromWeb()
		return h.readCached()
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
func (h *AoCHelper) cacheFromWeb() {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", h.Year, h.Day), nil)
	checkErr(err)
	req.Header.Add("Cookie", "session="+h.sessionCookie)
	resp, err := h.httpClient.Do(req)
	checkErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	err = os.WriteFile(filepath.Join(h.cacheDir, h.getFileName()), body, os.ModePerm)
}

// get filename for a given day
func (h *AoCHelper) getFileName() string {
	return fmt.Sprintf("%d%20d.txt", h.Year, h.Day)
}

// GetInput retrieves the puzzle input as a []string.
//
// Also caches the input, use force=true to force overwrite the cache. Please use caching to avoid unnecessary to the AoC servers
func (h *AoCHelper) GetInput(force bool) []string {
	if force {
		h.cacheFromWeb()
		return h.readCached()
	}
	return h.readCached()
}

// Submits the solution to the AoC website if it hasn't been already solved.
func (h *AoCHelper) Submit(part int, solution int) bool {
	if h.puzzleStatus >= part {
		return true
	}
	body := []byte(fmt.Sprintf("level=%d&answer=%d", part, solution))
	req, err := http.NewRequest("POST", fmt.Sprintf("https://adventofcode.com/%d/day/%d/answer", h.Year, h.Day), bytes.NewBuffer(body))
	checkErr(err)
	req.Header.Add("Cookie", "session="+h.sessionCookie)
	resp, err := h.httpClient.Do(req)
	checkErr(err)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("HTTP Response code != 200")
	}
	body, err = ioutil.ReadAll(resp.Body)
	checkErr(err)
	bodyString := string(body[:])
	if strings.Contains(bodyString, "That's not the right answer") {
		fmt.Println("That answer wasn't correct :c\nPlease wait 60 seconds before submitting the next solution!")
		return false
	}
	fmt.Println("That was the right answer!")
	return true
}

// getPuzzleStatus retrieves which Part of the Puzzle has already been solved
//
// 0 = unsolved, 1 = half solved, 2 = finished
func (h *AoCHelper) getPuzzleStatus() int {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d", h.Year, h.Day), nil)
	checkErr(err)
	req.Header.Add("Cookie", "session="+h.sessionCookie)
	resp, err := h.httpClient.Do(req)
	checkErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	bodyString := string(body[:])
	return strings.Count(bodyString, "Your puzzle answer was")
}

func (h *AoCHelper) printPuzzleStatus() {
	switch h.puzzleStatus {
	case 0:
		fmt.Println("Puzzle has not been solved yet!")
	case 1:
		fmt.Println("Part 1 of the puzzle has been solved. No longer submitting solutions for Part 1!")
	case 2:
		fmt.Println("The puzzle has been fully solved. No longer submitting any solutions to the AoC Website!")
	}
}
