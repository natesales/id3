// Author: Nate Sales (@nwsnate)

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

	Name        string
	Description string
	Children    map[string]Node
}

func handle(e error) {

	/* handle
	 * Description: Handle an error. Well actually just panic and quit.
	 */

	if e != nil {
		panic(e)
	}
}

func readDataSet(filename string) ([]map[string]string, []map[string]string, []string) {

	/* readDataSet
	 * Description: Read data from text file randomly into slices of Entry objects.
	 * filename: string Path to text file for reading
	 * returns: 2 []map[string]string objects containing random values from filename and the header
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
			// with training and testing being empty,
			// training = all[:len(all)/2]
			// testing = all[len(all)/2:]
		}
	}

	return training, testing, header
}

func Entropy(entries []map[string]string) float64 {

	/* Entropy
	 * Description: Calculate entropy for a given array of maps.
	 * entries: Array of maps containing all the entries
	 * returns: float64, entropy value
	 */
	//TODO: Allow for more than one category.
	var pYes float64
	var pNo float64
	for i, entry := range entries {
		if entry[categoryName] == outcomeTrue {
			pYes++
		} else if entry[categoryName] == outcomeFalse {
			pNo++
		} else {
			fmt.Println(entry)
			fmt.Println(entry)
			fmt.Println(categoryName + " is neither yes or no. It's '" + entry[categoryName] + "' with index " + strconv.Itoa(i))
			os.Exit(1)
		}
	}

	total := pYes + pNo
	pYes /= total
	pNo /= total

	if pYes == 0 {
		return -(pNo * math.Log2(pNo))
	}
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

// Begin id3 helpers

func sameCategory(entries []map[string]string) bool {
	/*
	 * Description: Are all entries in the same category?
	 */

	lastCategory := ""

	for _, entry := range entries { // Loop over entries
		currentCategory := entry[categoryName]

		if (lastCategory != "") && (currentCategory != lastCategory) { // If the category is not the same then return, they aren't all in the same category.
			return false
		} else {
			lastCategory = currentCategory
		}
	}

	return true
}

func uniqueValuesOf(entries []map[string]string, attribute string) []string {
	/*
	 * Description: Return array of unique values of attribute. (Essentially a set)
	 */

	valueMap := make(map[string]bool) // This is the recommended way to make a set in Go.
	unique := make([]string, 0, len(valueMap))

	for _, entry := range entries {
		valueMap[entry[attribute]] = true
	}

	for key := range valueMap {
		unique = append(unique, key)
	}

	return unique
}

func attribWithLargestGain(entries []map[string]string, attributes []string) string {
	var attribLargestGainSoFar = ""
	var largestGain = 0.0

	for _, attribute := range attributes { // Loop through the attributes
		currentAttribGain := Gain(entries, attribute) // Compute the gain

		if currentAttribGain >= largestGain { // If the current gain is larger, then update the currently known values and continue
			attribLargestGainSoFar = attribute
			largestGain = currentAttribGain
		}
	}

	return attribLargestGainSoFar
}

func mostCommon(entries []map[string]string, attribute string) (string, int) {

	/*
	 * Description: Compute the most common value of attribute in a given list of attributes.
	 * Returns: string, attribute that is the most common and a certainty percent which is the percentage of that one out of all of them.
	 */

	valueMap := make(map[string]int) // map[attribute value]number of times

	for _, entry := range entries {
		valueMap[entry[attribute]]++
	}

	mostCommon := ""
	valueMostCommon := 0

	for value := range valueMap {
		if valueMap[value] > valueMostCommon {
			mostCommon = value
			valueMostCommon = valueMap[value]
		}
	}

	return mostCommon, 0 // todo valueMostCommon / len(valueMap)
}

func deleteFrom(slice []string, item string) []string {

	/*
	 * Description: Delete item from slice
	 * Note: This is the recommended way to remove from a slice in Go. collection/list should have been used in this case.
	 * https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
	 */

	indexOfItem := indexOf(item, slice)
	if indexOfItem != -1 { // If item is in slice
		//return append(slice[:indexOfItem], slice[indexOfItem+1:]...)
		return append(slice[:indexOfItem], slice[indexOfItem+1:]...)
		//slice = append(slice[:0], slice[1:]...)
	} // If item is not in slice, then there's nothing to do.

	return slice // If it comes to this point, then item isnt in slice. So return.
}

func indexOf(element string, slice []string) int {
	/*
	 * Description: Find first occurrence of element in array.
	 * Returns: index (int)
	 */
	for index, item := range slice {
		if item == element {
			return index
		}
	}

	return -1 // Then element isn't in array.
}

// End id3 helpers

func id3(entries []map[string]string, attributes []string) Node {

	/* id3
	 * Description: Main id3 recursive function.
	 * entries: Array of entries.
	 * attributes: String array of attributes.
	 * Returns: Node
	 */

	// The ID3 function is given a list of examples and a list of possible attributes.
	//    If all of the examples belong to the same category, then return a leaf node labeled with that category.
	//    If there are no more attributes, return a leaf node labeled with the most common category in the examples.
	//    otherwise,
	//        select the attribute that results in the greatest information gain
	//        create (and eventually return) a non-leaf node that is labeled with that attribute
	//        For each value v of that attribute:
	//            create a child for that value by applying one of the following two options:
	//                If there are no examples with the value v, then the child is a leaf labeled with the most common category in the current examples
	//                otherwise, the child is the result of running ID3 recursively with the examples that have value v and all the remaining attributes
	if sameCategory(entries) { // If all of the examples belong to the same category, then return a leaf node labeled with that category.
		return Node{ // Return a leaf Node, notated by the Children map being nil.
			Name:     categoryName,
			Children: nil,
		}
	}

	// By this point they are not all in the same category.

	// For some reason len(attributes) goes down to 3 and stays there forever.

	//fmt.Println(attributes)
	if len(attributes) == 0 { // If there are no attributes, return a leaf node labeled with the most common category in the examples.
		mostCommon, _ := mostCommon(entries, categoryName)
		return Node{
			Name:     mostCommon,
			Children: nil, // This is a leaf.
		}
	}

	// If there are no examples with the value v, then the child is a leaf labeled with the most common category in the current examples
	// otherwise, the child is the result of running ID3 recursively with the examples that have value v and all the remaining attributes

	largestGain := attribWithLargestGain(entries, attributes) // select the attribute that results in the greatest information gain
	node := Node{                                             // create (and eventually return) a non-leaf node that is labeled with that attribute
		Name:     largestGain,
		Children: map[string]Node{}, // An empty map
	}

	// If there are no examples with the value v, then the child is a leaf labeled with the most common category in the current examples
	// otherwise, the child is the result of running ID3 recursively with the examples that have value v and all the remaining attributes

	for _, v := range uniqueValuesOf(entries, largestGain) { // For each value v of that attribute,
		// create a child for that value by applying one of the following two options:

		var subset []map[string]string // subset of examples that are share a given property (wind=weak)

		for _, e := range entries { // Fill up the subset of entries that have v as their value of largestGain
			if e[largestGain] == v {
				subset = append(subset, e)
			}
		}

		if len(subset) == 0 { // If there are no examples with the value v,
			mostCommon, _ := mostCommon(entries, categoryName)
			node.Children[v] = Node{ // Child is a leaf labeled with the most common category in the current examples
				Name:        mostCommon,
				Description: "I give up",
				Children:    nil,
			}

		} else { // There is a subset...
			// Otherwise, the child is the result of running ID3 recursively with the examples that have value v and all the remaining attributes
			node.Children[v] = id3(subset, deleteFrom(attributes, largestGain)) // delete largestGain from attributes
		}
	}

	return node
}

func printTree(root Node, indentation int) {

	/*
	 * Description: Print out the tree
	 */
	for i := 0; i < indentation; i++ {
		fmt.Print("\t")
	}

	fmt.Println(root.Name, root.Description)

	for _, child := range root.Children {
		newIndentation := indentation + 1
		printTree(child, newIndentation)
	}
}

func main() {
	var training []map[string]string
	var testing []map[string]string
	var header []string

	training, testing, header = readDataSet(filename)
	//_, _, header = readDataSet(filename)

	all := append(training, testing...)
	//fmt.Println("Total entropy:", Entropy(append(training, testing...)))
	//
	//fmt.Println("Training entropy:", Entropy(training))
	//fmt.Println("Testing entropy:", Entropy(testing))
	//
	//fmt.Println("Outlook Gain:", Gain(append(training, testing...), "outlook"))
	//fmt.Println("Humidity Gain:", Gain(append(training, testing...), "humidity"))
	//fmt.Println("Wind Gain:", Gain(append(training, testing...), "wind"))
	//fmt.Println("Temperature Gain:", Gain(append(training, testing...), "temperature"))

	header = deleteFrom(header, categoryName)
	//printTree(
	id3(all, header) //, 0)
}

//func findAttributes(entries []map[string]string) []string { // Why is this here
//	var attributes []string
//
//	for _, entry := range entries {
//		attributes = append(attributes, entry[])
//	}
//
//	return attributes
//}
