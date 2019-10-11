// Author: Nate Sales (@nwsnate)

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

var ( // Globals
	categoryName string
	categories   = make(map[string]int) // Category name and number of times it shows up.
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
	var all []map[string]string

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

		categories[entry[0]] += 1

		if len(header) == 0 { // Skip the first line since it contains the header.
			header = entry
			categoryName = header[0]
		} else {
			newEntry := make(map[string]string)

			for i, item := range header { // For each entry in the header...
				// This will fail if the subsequent (not header) entries length != len(header)
				newEntry[item] = entry[i] // Fill up the map
			}

			all = append(all, newEntry)

			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })

			training = all[:len(all)/2]
			testing = all[len(all)/2:]
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

	entryValues := make(map[string]int)

	for _, entry := range entries {
		entryValues[entry[categoryName]] += 1
	}

	final := 0.0

	total := float64(len(entries))

	for _, val := range entryValues {
		percentage := float64(val) / total
		final += -(percentage * math.Log2(percentage))
	}

	return final
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

func sameCategory(entries []map[string]string) (bool, string) {
	/*
	 * Description: Are all entries in the same category?
	 */

	lastCategory := ""
	var currentCategory string

	for _, entry := range entries { // Loop over entries
		currentCategory = entry[categoryName]

		if (lastCategory != "") && (currentCategory != lastCategory) { // If the category is not the same then return, they aren't all in the same category.
			return false, "" // If they're not in the same category, return an empty string
		} else {
			lastCategory = currentCategory
		}
	}

	return true, currentCategory
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
	samecategory, group := sameCategory(entries)
	if samecategory { // If all of the examples belong to the same category, then return a leaf node labeled with that category.
		return Node{ // Return a leaf Node, notated by the Children map being nil.
			Name:     categoryName + "=" + group,
			Children: nil,
		}
	}

	// By this point they are not all in the same category.

	// For some reason len(attributes) goes down to 3 and stays there forever.

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
		Children: make(map[string]Node), // An empty map
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

func indent(indentation int) string {

	/*
	 * Return `indentation` number of tabs.
	 */

	out := ""
	for i := 0; i < indentation; i++ {
		out += "\t"
	}
	return out
}

func printTree(root Node, indentation int) {

	/*
	 * Description: Print out the tree
	 */
	if root.Children == nil { // If leaf
		fmt.Println(indent(indentation), root.Name, root.Description)
	} else {
		fmt.Println(indent(indentation), "What is the", root.Name+"?", root.Description)
	}

	for index, child := range root.Children {
		fmt.Println(indent(indentation+1), index)
		printTree(child, indentation+2)
	}
}

func follow(entry map[string]string, root Node) string {

	/* follow
	 * Description: Follow the tree
	 */

	if root.Children != nil { // If not leaf
		return follow(entry, root.Children[entry[root.Name]])
	} else { // If the base case is reached	return ""
		return root.Name
	}
}

func main() {
	var training []map[string]string
	var testing []map[string]string
	var header []string

	filename := "data/mushrooms.txt"

	log.Print("Reading dataset", filename)

	start := time.Now()
	training, testing, header = readDataSet(filename)
	all := append(training, testing...)
	log.Print("Finished in ", time.Since(start))

	log.Println(len(all), "entries detected.")

	//fmt.Println("Total entropy:", Entropy(all))
	//
	//fmt.Println("Training entropy:", Entropy(training))
	//fmt.Println("Testing entropy:", Entropy(testing))
	//
	//fmt.Println("Outlook Gain:", Gain(all, "outlook"))
	//fmt.Println("Humidity Gain:", Gain(all, "humidity"))
	//fmt.Println("Wind Gain:", Gain(all, "wind"))
	//fmt.Println("Temperature Gain:", Gain(all, "temperature"))

	header = deleteFrom(header, categoryName)

	log.Println("Building tree")

	start = time.Now()
	tree := id3(training, header)
	log.Println("Tree building finished in", time.Since(start))

	//printTree(tree, 0)

	correct := 0
	incorrect := 0
	for _, entry := range testing {
		if strings.Split(follow(entry, tree), "=")[1] == entry[categoryName] {
			correct++
		} else {
			incorrect++
		}
	}

	log.Println(float64(correct)/float64(correct+incorrect)*100, "% Accuracy")
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
