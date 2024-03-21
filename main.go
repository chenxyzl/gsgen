package main

import (
	"fmt"
	"math/rand"
)

type Item struct {
	ID    string
	Count int
}

func getRandomItems(items map[string]int, n int) map[string]int {
	totalCount := 0
	for _, count := range items {
		totalCount += count
	}

	result := make(map[string]int)

	for i := 0; i < n; i++ {
		if len(items) == 0 {
			break
		}

		r := rand.Intn(totalCount)
		countSum := 0

		for id, count := range items {
			countSum += count
			if r < countSum {
				result[id] += 1
				items[id]--
				totalCount--
				if items[id] == 0 {
					delete(items, id)
				}
				break
			}
		}
	}

	return result
}

func main() {
	items := map[string]int{
		"道具1": 2,
		"道具2": 10,
	}

	n := 3
	randomItems := getRandomItems(items, n)

	fmt.Println("随机选择的道具：")
	for k, v := range randomItems {
		fmt.Printf("%s x%d\n", k, v)
	}
}
