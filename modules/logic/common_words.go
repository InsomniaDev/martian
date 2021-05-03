package logic

import (
	"sort"
	"strings"
)

var commonWords []string

// RemoveCommonWords will remove the common words from the string of words that was provided
func RemoveCommonWords(common []string, stringToParse string) []string {
	// Setup variables for efficiency, do once and use everywhere
	commonWords = common
	sort.Strings(commonWords)

	// Split the provided string by spaces
	parseString := strings.Split(stringToParse, " ")

	// Get a list of the words that we can now search by
	goodWords := []string{}
	for _, a := range parseString {
		// Only if the word is not common will we add it to the list of words to check
		if !checkIfWordIsCommon(a) {
			// Add to list
			goodWords = append(goodWords, a)
		}
	}
	return goodWords
}

// checkIfWordIsCommon does a binary search and returns if it is a common word
func checkIfWordIsCommon(word string) bool {
	i := sort.Search(len(commonWords),
		func(i int) bool { return commonWords[i] >= word })
	if i < len(commonWords) && commonWords[i] == word {
		return true
	}
	return false
}
