package models

import (
  "fmt"
  "strconv"
  "strings"
)

type Object struct {
  FavoriteFruit string `json:"favoriteFruit"`
  Posts []Post `json:"posts"`
  Name Name `json:"name"`
  Age int `json:"age"`
  Balance string `json:"balance"`
  Active bool `json:"isActive"`
}

func (o *Object) Transform() *TransformedObject {
  // Count each time a word appears in a post.
  wordCount := make(map[string]int)
  for _, post := range o.Posts {
    words := strings.Fields(post.Text)
    for _, word := range words {
      lastIndex := len(word)-1

      // Remove the last character of the word if it is a period. Period was the only punctuation found in the sample
      // data but in production we would need to check for more types of punctuation.
      if word[lastIndex] == '.' {
        word = word[:lastIndex]
      }

      normalizedWord := strings.ToLower(word)
      count := wordCount[normalizedWord]
      wordCount[normalizedWord] = count+1
    }
  }

  // Find the most appearances for a single word.
  max := 0
  for _, count := range wordCount {
    if count > max {
      max = count
    }
  }

  // Create a list of the words with the most appearances.
  var mostCommonWords []string
  for word, count := range wordCount {
    if count == max {
      mostCommonWords = append(mostCommonWords, word)
    }
  }

  // Errors resulting from parsing the float value are ignored here. If an error occurs the default float value of 0
  // will be used.
  balance, _ := strconv.ParseFloat(strings.Replace(o.Balance[1:], ",", "", -1), 64)

  return &TransformedObject{
    FullName: o.Name.fullName(),
    PostCount: len(o.Posts),
    MostCommonWords: mostCommonWords,
    Age: o.Age,
    Active: o.Active,
    FavoriteFruit: o.FavoriteFruit,
    Balance: balance,
  }
}

type Post struct {
  Text string `json:"post"`
}

type Name struct {
  First string `json:"first"`
  Last string `json:"last"`
}

func (n *Name) fullName() string {
  return fmt.Sprintf("%s %s", n.First, n.Last)
}

type TransformedObject struct {
  FullName string `json:"full_name"`
  PostCount int `json:"post_count"`
  MostCommonWords []string `json:"most_common_word_in_posts"`
  Age int `json:"age"`
  Active bool `json:"is_active"`
  FavoriteFruit string `json:"favorite_fruit"`
  Balance float64 `json:"balance"`
}