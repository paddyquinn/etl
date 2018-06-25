package models

import "math"

const apple = "apple"

type Analysis struct {
  TotalPostCount int                     `json:"total_post_count"`
  MostCommonWord []string                `json:"most_common_word_overall"`
  Balance        balanceStatistics       `json:"account_balance"`
  Age            ageStatistics           `json:"age"`
  FavoriteFruit  favoriteFruitStatistics `json:"favorite_fruit"`
}

type balanceStatistics struct {
  Total          float64 `json:"total"`
  Mean           float64 `json:"mean"`
  ActiveMean     float64 `json:"active_user_mean"`
  StrawberryMean float64 `json:"strawberry_lovers_mean"`
}

type ageStatistics struct {
  Min                int     `json:"min"`
  Max                int     `json:"max"`
  Mean               float64 `json:"mean"`
  Median             float64 `json:"median"`
  AppleAge           int     `json:"age_with_most_apple_lovers"`
  YoungestAppleHater int     `json:"youngest_age_hating_apples"`
  OldestAppleHater   int     `json:"oldest_age_hating_apples"`
}

type favoriteFruitStatistics struct {
  Active []string `json:"active_users"`
  Median []string `json:"median_age"`
  Rich   []string `json:"acct_balance_gt_mean"`
}

type Analyzer struct {
  transformedObjects []*TransformedObject
}

func NewAnalyzer(transformedObjects []*TransformedObject) *Analyzer {
  return &Analyzer{transformedObjects: transformedObjects}
}

func (a *Analyzer) Analyze() *Analysis {
  // Count the total posts and create a count of all of the most common words.
  totalPostCount := 0
  wordCount := make(map[string]int)
  for _, obj := range a.transformedObjects {
    totalPostCount += obj.PostCount
    for _, word := range obj.MostCommonWords {
      count := wordCount[word]
      wordCount[word] = count+1
    }
  }

  balanceStatistics := a.calculateBalanceStatistics()
  ageStatistics := a.calculateAgeStatistics()

  return &Analysis{
    TotalPostCount: totalPostCount,
    MostCommonWord: findMostCommon(wordCount),
    Balance:        balanceStatistics,
    Age:            ageStatistics,
    FavoriteFruit:  a.calculateFavoriteFruitStatistics(ageStatistics.Median, balanceStatistics.Mean),
  }
}

func (a *Analyzer) calculateBalanceStatistics() balanceStatistics {
  var total, activeTotal, strawberryTotal, activeLength, strawberryLength float64
  for _, obj := range a.transformedObjects {
    total += obj.Balance

    if obj.Active {
      activeTotal += obj.Balance
      activeLength++
    }

    if obj.FavoriteFruit == "strawberry" {
      strawberryTotal += obj.Balance
      strawberryLength++
    }
  }

  return balanceStatistics{
    Total:          total,
    Mean:           total / float64(len(a.transformedObjects)),
    ActiveMean:     activeTotal / activeLength,
    StrawberryMean: strawberryTotal / strawberryLength,
  }
}

func (a *Analyzer) calculateAgeStatistics() ageStatistics {
  min := math.MaxInt64
  max := -1
  var (
    total float64
    ages  intSlice
  )
  appleAge := -1
  appleCount := make(map[int]int)
  youngestAppleHater := math.MaxInt64
  oldestAppleHater := -1
  for _, obj := range a.transformedObjects {
    if obj.Age < min {
      min = obj.Age
    }

    if obj.Age > max {
      max = obj.Age
    }

    total += float64(obj.Age)

    ages.data = append(ages.data, obj.Age)

    if obj.FavoriteFruit == apple {
      count := appleCount[obj.Age]
      count += 1
      appleCount[obj.Age] = count
      if count > appleAge {
        appleAge = obj.Age
      }
    }

    if obj.FavoriteFruit != apple && obj.Age < youngestAppleHater {
      youngestAppleHater = obj.Age
    }

    if obj.FavoriteFruit != apple && obj.Age > oldestAppleHater {
      oldestAppleHater = obj.Age
    }
  }

  return ageStatistics{
    Min:                min,
    Max:                max,
    Mean:               total / float64(len(a.transformedObjects)),
    Median:             ages.median(),
    AppleAge:           appleAge,
    YoungestAppleHater: youngestAppleHater,
    OldestAppleHater:   oldestAppleHater,
  }
}

func (a *Analyzer) calculateFavoriteFruitStatistics(medianAge, meanBalance float64) favoriteFruitStatistics {
  medianAgeInt := int(medianAge)
  activeCount := make(map[string]int)
  medianCount := make(map[string]int)
  richCount := make(map[string]int)
  for _, obj := range a.transformedObjects {
    if obj.Active {
      count := activeCount[obj.FavoriteFruit]
      activeCount[obj.FavoriteFruit] = count+1
    }

    if obj.Age == medianAgeInt {
      count := medianCount[obj.FavoriteFruit]
      medianCount[obj.FavoriteFruit] = count+1
    }

    if obj.Balance > meanBalance {
      count := richCount[obj.FavoriteFruit]
      richCount[obj.FavoriteFruit] = count+1
    }
  }

  return favoriteFruitStatistics{
    Active: findMostCommon(activeCount),
    Median: findMostCommon(medianCount),
    Rich:   findMostCommon(richCount),
  }
}