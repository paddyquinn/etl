package models

import "sort"

type intSlice struct {
  data []int
}

func (is intSlice) Len() int {
  return len(is.data)
}

func (is intSlice) Less(i, j int) bool {
  return is.data[i] < is.data[j]
}

func (is intSlice) Swap(i, j int) {
  is.data[i], is.data[j] = is.data[j], is.data[i]
}

func (is intSlice) median() float64 {
  sort.Sort(is)
  length := len(is.data)
  midpoint := length/2
  if length%2 != 0 {
    return float64(is.data[midpoint])
  }

  return float64(is.data[midpoint-1]+is.data[midpoint])/2
}

func findMostCommon(appearanceCount map[string]int) []string {
  // Find the most appearances for a single word.
  max := 0
  for _, count := range appearanceCount {
    if count > max {
      max = count
    }
  }

  // Create a list of the words with the most appearances.
  var mostCommonWords []string
  for word, count := range appearanceCount {
    if count == max {
      mostCommonWords = append(mostCommonWords, word)
    }
  }

  return mostCommonWords
}