package main

import (
	"flag"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/gosuri/uiprogress"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const helpString = `stanford-dl
Author : https://github.com/coderick14

A dead simple script to download videos or pdfs from Stanford Engineering Everywhere.
USAGE : stanford-dl -course COURSE_CODE [-type {video|pdf}] [-all] [-lec lectures] [--help]

--course    Course name e.g. CS229, EE261
--type 	    Specify whether to download videos or pdfs. Defaults to PDF.
--all       Download for all lectures
--lec       Comma separated list of lectures e.g. 1,3,5,10
--help      Display this help message and quit

Found a bug? Feel free to raise an issue on https://github.com/coderick14/stanford-dl
Contributions welcome :)`

// Wrapper over io.Reader to record progresses
type passThrough struct {
	io.Reader
	index int
	curr  int
	total int
}

var bars = make([]*uiprogress.Bar, 0, 50)
var factor int64

// Override Read method of io.Reader
func (pt *passThrough) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.curr += n

	if err == nil || (err == io.EOF && n > 0) {
		bars[pt.index].Set(int((float64(pt.curr/int(factor)) / float64(pt.total)) * float64(pt.total)))
	}

	return n, err
}

// Goroutine to download a lecture
func downloadLecture(index int, url string, fileName string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Send GET request to required url
	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Error while downloading", fileName)
		return
	}

	defer resp.Body.Close()

	// Open file for writing
	fh, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Error while creating file", fileName)
		return
	}

	defer fh.Close()

	// Initialize the progress bar for this lecture
	bars[index] = uiprogress.AddBar(int(resp.ContentLength / factor)).AppendCompleted()
	bars[index].PrependFunc(func(b *uiprogress.Bar) string {
		return "Downloading " + fileName
	})

	// Create wrapper over io.Reader
	src := &passThrough{Reader: resp.Body, total: int(resp.ContentLength / factor), index: index}
	_, err = io.Copy(fh, src)

	if err != nil {
		bars[index].AppendFunc(func(b *uiprogress.Bar) string {
			return "Failed"
		})
		return
	}

	// Finished downloading
	bars[index].AppendFunc(func(b *uiprogress.Bar) string {
		return "Completed"
	})
}

// Utility function to create a range of numbers
func makeRange(n int) []int {
	var list = make([]int, n)

	for i := 0; i < n; i++ {
		list[i] = i + 1
	}

	return list
}

// Utility function to return a list of formatted lecture ids
func createLectureList(all bool, lectures string, lectureCount int) []int {
	var lectureList []int

	if all {
		lectureList = makeRange(lectureCount)
	} else {
		tempList := strings.Split(lectures, ",")
		for _, num := range tempList {
			val, _ := strconv.Atoi(num)
			lectureList = append(lectureList, val)
		}
	}

	return lectureList
}

func main() {

	// Define flags and base URLs
	var (
		help          = flag.Bool("help", false, "Display help")
		courseName    = flag.String("course", "", "Course name e.g. CS229, EE261")
		typeFlag      = flag.String("type", "pdf", "[video | pdf]. Defaults to pdf.")
		all           = flag.Bool("all", false, "Download material for all lectures for the given course")
		lectures      = flag.String("lec", "", "Specify comma separated list of lectures e.g 1,3,10")
		siteBaseURL   = "https://see.stanford.edu"
		courseBaseURL = "https://see.stanford.edu/Course/"
		videoBaseURL  = "http://html5.stanford.edu/videos/courses/see/"
	)

	// Parse the command line flags
	flag.Parse()

	// Display help and quit
	if *help == true {
		fmt.Println(helpString)
		return
	}

	// Check for required -course flag
	if len(*courseName) == 0 {
		fmt.Println("Please specify a Course code")
		return
	}

	// Check for valid value for -type flag
	if strings.Compare("pdf", *typeFlag) != 0 && strings.Compare("video", *typeFlag) != 0 {
		fmt.Println("[video | pdf] are the only accepted values for -type flag")
		return
	}

	courseURL := courseBaseURL + *courseName

	// Get HTML content of the course page
	resp, err := soup.Get(courseURL)

	if err != nil {
		fmt.Println("Error fetching course details. Check your internet connection!!")
		return
	}

	// Parse HTML content of course page
	doc := soup.HTMLParse(resp)

	// Set BaseURL, path and file extension
	var (
		baseURL, extension string
		paths              []string
		linkParentTags     []soup.Root
		lectureCount       int
		lectureList        []int
	)

	if strings.Compare(*typeFlag, "video") == 0 {
		// For Videos
		baseURL = videoBaseURL
		extension = "mp4"
		linkParentTags = doc.FindAll("table", "class", "table")
		lectureList = createLectureList(*all, *lectures, len(linkParentTags))
		lectureCount = len(lectureList)
		factor = 1000000
		for i := 0; i < lectureCount; i++ {
			paths = append(paths, fmt.Sprintf("%s/%s-lecture%02d.%s", *courseName, *courseName, lectureList[i], extension))
		}

	} else {
		// For PDFs
		baseURL = siteBaseURL
		extension = "pdf"
		linkParentTags = doc.FindAll("ul", "class", "list-inline")
		lectureList = createLectureList(*all, *lectures, len(linkParentTags))
		lectureCount = len(lectureList)
		factor = 1000
		for i := 0; i < lectureCount; i++ {
			lectureId := lectureList[i]
			href := linkParentTags[lectureId-1].FindAll("a")[1].Attrs()["href"]
			paths = append(paths, href)
		}

	}

	// Resize the progress bar array
	bars = bars[:lectureCount]

	fmt.Printf("Found %d lectures for course %s\n", lectureCount, *courseName)
	var wg sync.WaitGroup

	// Listen for download progresses
	uiprogress.Start()

	for i := 0; i < lectureCount; i++ {
		url := baseURL + paths[i]

		fileName := fmt.Sprintf("%s-lecture%02d.%s", *courseName, lectureList[i], extension)

		// fetch lecture concurrently
		wg.Add(1)
		go downloadLecture(i, url, fileName, &wg)
	}

	// Wait for all lectures to be downloaded
	wg.Wait()
}
