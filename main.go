// Author: Nate Sales (@nwsnate) nate.cx
// Revision: September 12, 2019

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

func handle(e error) {
	/**
	 * handle
	 * Description: Handle an error. Well actually just panic and quit.
	 */
	if e != nil {
		panic(e)
	}
}

type Entry struct {
	PlayTennis  string
	Outlook     string
	Temperature string
	Humidity    string
	Wind        string
}

func readDataSet(filename string) ([]Entry, []Entry, []Entry) {

	/**
	 * readDataSet
	 * Description: Read data from text file randomly into slices of Entry objects.
	 * filename: string Path to text file for reading
	 * returns 2 []Entry objects containing random values from filename as well as a third with all the Entry objects
	 */

	var training []Entry
	var testing []Entry
	var total []Entry

	r := rand.New(rand.NewSource(time.Now().Unix())) // Not a good seed. TODO

	file, err := os.Open(filename)
	handle(err)
	defer func() {
		err := file.Close()
		handle(err)
	}()

	first := true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if first { // Skip the first line since it contains the header. TODO: Better way to do this?
			first = false
		} else {
			entry := strings.Split(scanner.Text(), ",")
			newEntry := Entry{
				PlayTennis:  entry[0],
				Outlook:     entry[1],
				Temperature: entry[2],
				Humidity:    entry[3],
				Wind:        entry[4],
			}

			if r.Intn(2) == 1 {
				training = append(training, newEntry)
			} else {
				testing = append(testing, newEntry)
			}
			total = append(total, newEntry)
		}

	}

	return training, testing, total
}

func entropy(entries []Entry) float64 { // TODO: Should this be float64? Its what Log2 returns.
	var pYes float64
	var pNo float64
	for _, entry := range entries {
		if entry.PlayTennis == "yes" {
			pYes++
		} else if entry.PlayTennis == "no" {
			pNo++
		} else {
			log.Fatalf("PlayTennis is neither yes or no. It's " + entry.PlayTennis)
		}
	}
	pYes /= pYes + pNo
	pNo /= pYes + pNo
	return -(pYes * math.Log2(pYes)) - (pNo * math.Log2(pNo))
}

func main() {
	var training []Entry
	var testing []Entry
	var total []Entry

	training, testing, total = readDataSet("data/tennis.txt")

	fmt.Println("Training entropy: ", entropy(training))
	fmt.Println("Testing entropy: ", entropy(testing))

	fmt.Println("Total entropy: ", entropy(total))

	//for i, line := range training {
	//	fmt.Println(i, line)
	//}
	//
	//fmt.Println("---------------")
	//
	//for i, line := range testing {
	//	fmt.Println(i, line)
	//}
}
