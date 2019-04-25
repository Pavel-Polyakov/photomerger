package main

import (
	"math/rand"
	"time"
)

type interval struct {
	Min, Max int
}

func getIntervals(max, n int) []interval {
	step := max / n

	var result []interval

	prevMin := -1

	for current := 0; current < max; current += step {
		if prevMin == -1 {
			prevMin = current
			continue
		}
		result = append(result, interval{Min: prevMin, Max: current})
		prevMin = current
	}

	if result[len(result)-1].Max < max {
		result = append(result, interval{Min: result[len(result)-1].Max, Max: max})
	}

	return result
}

func shuffle(arr []interval) {
	rand.Seed(int64(time.Now().Nanosecond()))
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i)
		arr[i], arr[j] = arr[j], arr[i]
	}
}
