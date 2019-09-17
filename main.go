// Author: Nate Sales (@nwsnate) nate.cx
// Revision: September 17, 2019

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

const ( // Global constants
	categoryName = "play tennis"
	outcomeTrue  = "yes"
	outcomeFalse = "no" // TODO: There could be more than 2 outcomes. Take this from the outcome column

	filename = "data/tennis.txt"
)

type Node struct {
	/* Node
	 * Name (string) Name of value
	 * Children (map[string]Node) Keys: Value of attribute, Values: Child node
	 */

	Name     string
	Children map[string]Node
}

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
	defer func() { // Deallocate
		err := file.Close()
		handle(err)
	}()

	var header []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		reader := csv.NewReader(strings.NewReader(scanner.Text()))
		entry, err := reader.Read() // Read in one line
		if err == io.EOF {
			break
		}
		handle(err)

		if len(header) == 0 { // Skip the first line since it contains the header.
			header = entry
		} else {
			newEntry := make(map[string]string)

			for i, item := range header { // For each entry in the header...
				// This will fail if the subsequent (not header) entries length != len(header)
				newEntry[item] = entry[i] // Fill up the map
			}

			// TODO Sometimes training entropy is NaN because the RNG can set one array to totally empty.
			if r.Intn(2) == 1 { // Split it up in roughly half.
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
		if entry[categoryName] == outcomeTrue {
			pYes++
		} else if entry[categoryName] == outcomeFalse {
			pNo++
		} else {
			fmt.Println(entry)
			fmt.Println(categoryName + " is neither yes or no. It's '" + entry[categoryName] + "' with index " + strconv.Itoa(i))
			os.Exit(1)
		}
	}

	total := pYes + pNo
	pYes /= total
	pNo /= total

	if pNo == 0 {
		return -(pYes * math.Log2(pYes))
	} else {
		return -(pYes * math.Log2(pYes)) - (pNo * math.Log2(pNo))
	}
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
		gain -= (float64(len(Sv1)) / float64(len(S))) * Entropy(Sv1) // Compute Entropy gain for given attribute
	}

	return gain
}

func id3(entries []map[string]string, attributes []string) Node {

	/* id3
	 * Description: Main id3 recursive function.
	 */

	var lastCategory string
	for _, entry := range entries { // Check if they are all the same category (outcome)
		currentCategory := entry[categoryName]
		if currentCategory != lastCategory { // If the category is not the same, then there is more work to do.
			/*
				select the attribute that results in the greatest information gain
					create (and eventually return) a non-leaf node that is labeled with that attribute
					For each value v of that attribute:
					create a child for that value by applying one of the following two options:
					If there are no examples with the value v, then the child is a leaf labeled with the most common category in the current examples
					otherwise, the child is the result of running ID3 recursively with the examples that have value v and all the remaining attributes
			*/
		} else {
			lastCategory = currentCategory // If its the same, then keep going.
		}
	}

	// If nothing has been returned by now, then they are all in the same category. (Yes/No)
	return Node{ // Return a leaf Node, notated by the Children map being nil.
		Name:     categoryName,
		Children: nil,
	}
	/*
		If all of the examples belong to the same category, then return a leaf node labeled with that category.
		  If there are no more attributes, return a leaf node labeled with the most common category in the examples.
	*/
}

func main() {
	var training []map[string]string
	var testing []map[string]string

	training, testing = readDataSet(filename)
	fmt.Println("Total entropy:", Entropy(append(training, testing...)))
	fmt.Println("Training entropy:", Entropy(training))
	fmt.Println("Testing entropy:", Entropy(testing))

	fmt.Println()

	fmt.Println("Outlook Gain:", Gain(append(training, testing...), "outlook"))
	fmt.Println("Humidity Gain:", Gain(append(training, testing...), "humidity"))
	fmt.Println("Wind Gain:", Gain(append(training, testing...), "wind"))
	fmt.Println("Temperature Gain:", Gain(append(training, testing...), "temperature"))

	//id3(append(training, testing...), []string{"outlook"})

}
