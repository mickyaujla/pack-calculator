package packing

import (
	"errors"
	"slices"
)

type Pack struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

const MaxItems = 10_000_000

func DefaultSizes() []int {
	return []int{250, 500, 1000, 2000, 5000}
}

func Calculate(items int, sizes []int) ([]Pack, error) {
	if items <= 0 {
		return nil, errors.New("items must be greater than 0")
	}

	if items > MaxItems {
		return nil, errors.New("items exceeds the maximum order size")
	}

	if len(sizes) == 0 {
		return nil, errors.New("at least one pack size is required")
	}

	for _, size := range sizes {
		if size <= 0 {
			return nil, errors.New("pack sizes must be greater than 0")
		}
	}

	smallestPack := slices.Min(sizes)
	largestPossibleTotal := items + smallestPack - 1

	fewestPacksFor := make([]int, largestPossibleTotal+1)
	lastPackAddedFor := make([]int, largestPossibleTotal+1)

	for total := 1; total <= largestPossibleTotal; total++ {
		fewestPacksFor[total] = -1

		for _, size := range sizes {
			if size > total {
				continue
			}

			packsBeforeThisOne := fewestPacksFor[total-size]
			if packsBeforeThisOne == -1 {
				continue
			}

			if fewestPacksFor[total] == -1 || packsBeforeThisOne+1 < fewestPacksFor[total] {
				fewestPacksFor[total] = packsBeforeThisOne + 1
				lastPackAddedFor[total] = size
			}
		}
	}

	targetTotal := items

	for fewestPacksFor[targetTotal] == -1 {
		targetTotal++
	}

	counts := map[int]int{}
	remainder := targetTotal

	for remainder > 0 {
		size := lastPackAddedFor[remainder]
		counts[size]++
		remainder -= size
	}

	var packs []Pack
	for size, count := range counts {
		if count > 0 {
			packs = append(packs, Pack{Size: size, Count: count})
		}
	}

	slices.SortFunc(packs, func(a, b Pack) int { return b.Size - a.Size })
	return packs, nil
}
