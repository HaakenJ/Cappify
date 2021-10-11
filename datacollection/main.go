package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// where to store images
const imgPath string = "./data/images/"

// where the csv containing urls and species is located
const csvPath string = "./data/2019_first_half.csv"
const numberOfWorker int = 20

const urlIndex int = 0  // index of image url in csv
const nameIndex int = 1 // index of scientific name in csv
const reportRate = 100  // report progress every 100 download
// var r, _ = regexp.Compile("[0-9]+$")
// regex used to get url number for naming convention of image files (always unique)
var r, _ = regexp.Compile("\\d{5,}")

type data struct {
	url  string
	name string
}

func main() {
	input := make(chan data, 10)
	output := make(chan bool, 10)
	var counter int = 0 // use to name downloaded image files
	var queue []data

	// open the csv file
	csvFile, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(csvFile)

	// read csv
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// don't store images of lichen
		if !strings.Contains(line[nameIndex], "lichen") {
			queue = append(queue, data{
				url:  line[urlIndex],
				name: line[nameIndex],
			})
		}
	}

	queue = queue[1:] // remove title row

	/* init worker */
	for i := 0; i < numberOfWorker; i++ {
		url := queue[0]
		queue = queue[1:]

		counter++
		input <- url
		go worker(input, output)
	}

	for i := 0; i < len(queue); i++ {
		data := queue[i]

		<-output
		input <- data
		go worker(input, output)

		if i%reportRate == 0 {
			fmt.Println(i)
		}
	}
}

func worker(input chan data, done chan bool) {
	data := <-input
	url := data.url
	name := data.name

	resp, err := http.Get(url) // get the data
	if err != nil {
		log.Println(err)
		return
	}

	// skip failed get requests
	if resp.StatusCode != 200 {
		log.Println("Failed request.")
		done <- true
	}

	// await the response
	defer resp.Body.Close()

	path := imgPath + name + "/" // create path to image directory

	// create directory if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	// find the numerical code in the image url
	index := r.FindStringSubmatch(url)

	// create the image file
	file, err := os.Create(path + strings.Join(index, "") + "2019" + ".jpeg")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	// copy the data from the response to the image file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	done <- true
}
