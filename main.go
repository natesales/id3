// Author: Nate Sales (@nwsnate) nate.cx
// Revision: September 12, 2019

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func handle(e error) {

	/* handle
	 * Description: Handle an error. Well actually just panic and quit.
	 */

	if e != nil {
		panic(e)
	}
}

func readDataSet(filename string) ([]map[string]string, []map[string]string) {

	/* readDataSet
	 * Description: Read data from text file randomly into slices of Entry objects.
	 * filename: string Path to text file for reading
	 * returns: 2 []map[string]string objects containing random values from filename
	 */

	var training []map[string]string
	var testing []map[string]string

	r := rand.New(rand.NewSource(time.Now().Unix())) // Not a good seed.

	file, err := os.Open(filename)
	handle(err)
	defer func() {
		err := file.Close()
		handle(err)
	}()

	var header []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		reader := csv.NewReader(strings.NewReader(scanner.Text()))

		for {
			entry, err := reader.Read()
			if err == io.EOF {
				break
			}
			handle(err)

			if len(header) == 0 { // Skip the first line since it contains the header.
				header = entry
			} else {
				newEntry := make(map[string]string)

				for i, item := range header {
					newEntry[item] = entry[i]
				}

				if r.Intn(2) == 1 { // Split it up in roughly half
					training = append(training, newEntry)
				} else {
					testing = append(testing, newEntry)
				}
			}
		}
	}

	return training, testing
}

func Entropy(entries []map[string]string) float64 {

	/* Entropy
	 * Description: Calculate entropy for a given array of maps.
	 * entries: Array of maps containing all the entries
	 * returns: float64, entropy value
	 */

	var pYes float64
	var pNo float64
	for i, entry := range entries {
		if entry["play tennis"] == "yes" {
			pYes++
		} else if entry["play tennis"] == "no" {
			pNo++
		} else {
			fmt.Println(entry)
			fmt.Println("PlayTennis is neither yes or no. It's '" + entry["PlayTennis"] + "' with index " + strconv.Itoa(i))
			os.Exit(1)
		}
	}

	total := pYes + pNo
	pYes /= total
	pNo /= total
	return -(pYes * math.Log2(pYes)) - (pNo * math.Log2(pNo))
}

func Gain(S []map[string]string, A string) float64 {

	/* Gain
	 * Description: Calculate entropy gain for a given attribute.
	 * returns: float64, entropy gain value
	 */

	gain := Entropy(S)
	values := make(map[string][]map[string]string) // All the possible values of A

	for _, entry := range S { // Fill values map up with how many of each value of A there is in S
		values[entry[A]] = append(values[entry[A]], entry)
	}

	for value := range values {
		Sv1 := values[value]
		gain -= (float64(len(Sv1)) / float64(len(S))) * Entropy(Sv1)
	}

	return gain
}

func id3(entries []map[string]string, attributes []string) {

	/* id3
	 * Description: Main id3 recursive function.
	 */

	/*
		If all of the examples belong to the same category, then return a leaf node labeled with that category.
		  If there are no more attributes, return a leaf node labeled with the most common category in the examples.
		  otherwise,
		      select the attribute that results in the greatest information gain
		      create (and eventually return) a non-leaf node that is labeled with that attribute
		      For each value v of that attribute:
		          create a child for that value by applying one of the following two options:
		              If there are no examples with the value v, then the child is a leaf labeled with the most common category in the current examples
		              otherwise, the child is the result of running ID3 recursively with the examples that have value v and all the remaining attributes
	*/
}

func main() {
	var training []map[string]string
	var testing []map[string]string

	training, testing = readDataSet("data/tennis.txt")
	fmt.Println("Total entropy:", Entropy(append(training, testing...)))
	fmt.Println("Training entropy:", Entropy(training))
	fmt.Println("Testing entropy:", Entropy(testing))

	fmt.Println("Outlook Gain:", Gain(append(training, testing...), "outlook"))
}

//TODO: Sometimes training entropy is NaN because the RNG can set one array to totally empty.
