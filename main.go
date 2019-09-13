// Author: Nate Sales (@nwsnate) nate.cx
// Revision: September 12, 2019

package main

import (
	"bufio"
	"fmt"
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
		entry := strings.Split(scanner.Text(), ",")

		if len(header) == 0 { // Skip the first line since it contains the header.
			header = entry
		} else {
			newEntry := make(map[string]string)

			for i, item := range header {
				//fmt.Println("Assigning " + item + " to " + entry[i])
				newEntry[item] = entry[i]
			}

			if r.Intn(2) == 1 { // Split it up in roughly half
				training = append(training, newEntry)
			} else {
				testing = append(testing, newEntry)
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

	var Sv1 []map[string]string
	var Sv2 []map[string]string

	for _ /* index */, entry := range S { //TODO: where do the 2 possible values come from?
		if entry[A] == "yes" {
			Sv1 = append(Sv1, entry)
		} else if entry[A] == "no" {
			Sv2 = append(Sv2, entry)
		}
	}
	return Entropy(S) - float64(len(Sv1)/len(S))*Entropy(Sv1) - float64(len(Sv2)/len(S))*Entropy(Sv2)
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

	//fmt.Println(Entropy(append(training, testing...))) // ... is to append the slice

	training, testing = readDataSet("data/tennis.txt")
	fmt.Println(Entropy(append(training, testing...)))
	fmt.Println("Training entropy: ", Entropy(training))
	fmt.Println("Testing entropy: ", Entropy(testing))
}

//TODO: Sometimes training entropy is NaN because the RNG can set one array to totally empty.
