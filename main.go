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

func readDataSet(filename string) ([]map[string]string, []map[string]string) {

	/**
	 * readDataSet
	 * Description: Read data from text file randomly into slices of Entry objects.
	 * filename: string Path to text file for reading
	 * returns 2 []map[string]string objects containing random values from filename
	 */

	var training []map[string]string
	var testing []map[string]string

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
		if first { // Skip the first line since it contains the header. TODO: Dont use struct, use a map
			first = false
		} else {
			entry := strings.Split(scanner.Text(), ",")
			newEntry := map[string]string{
				"PlayTennis":  entry[0],
				"Outlook":     entry[1],
				"Temperature": entry[2],
				"Humidity":    entry[3],
				"Wind":        entry[4],
			}

			if r.Intn(2) == 1 {
				training = append(training, newEntry)
			} else {
				testing = append(testing, newEntry)
			}
		}

	}

	return training, testing
}

func Entropy(entries []map[string]string) float64 { // TODO: Should this be float64? Its what Log2 returns.
	var pYes float64
	var pNo float64
	for _, entry := range entries {
		if entry["PlayTennis"] == "yes" {
			pYes++
		} else if entry["PlayTennis"] == "no" {
			pNo++
		} else {
			log.Fatalf("PlayTennis is neither yes or no. It's " + entry["PlayTennis"])
		}
	}
	pYes /= pYes + pNo
	pNo /= pYes + pNo
	return -(pYes * math.Log2(pYes)) - (pNo * math.Log2(pNo))
}

//func Gain(S []map[string]string, A string) float64 { // TODO: A is going to be problematic
//	return Entropy(S) - (math.Abs(Sv1)/float64(len(S)))*Entropy(Sv1) - (math.Abs(Sv2)/float64(len(S)))*Entropy(Sv2)
//}

func main() {
	var training []map[string]string
	var testing []map[string]string

	training, testing = readDataSet("data/tennis.txt")
	//
	//fmt.Println("Training entropy: ", Entropy(training))
	//fmt.Println("Testing entropy: ", Entropy(testing))

	fmt.Println(Entropy(append(training, testing...))) // ... is to

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
