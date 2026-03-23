package solution

func MinEatingSpeed(piles []int, h int) int {
	left, right := 1, maxPile(piles)
	for left < right {
		mid := left + (right-left)/2
		if hoursNeeded(piles, mid) <= h {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
}

func hoursNeeded(piles []int, speed int) int {
	total := 0
	for _, p := range piles {
		total += (p + speed - 1) / speed
	}
	return total
}

func maxPile(piles []int) int {
	m := piles[0]
	for _, p := range piles[1:] {
		if p > m {
			m = p
		}
	}
	return m
}
