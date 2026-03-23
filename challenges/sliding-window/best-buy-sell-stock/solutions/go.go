package solution

func MaxProfit(prices []int) int {
	minPrice := prices[0]
	best := 0
	for _, p := range prices {
		if p < minPrice {
			minPrice = p
		}
		if profit := p - minPrice; profit > best {
			best = profit
		}
	}
	return best
}
