package solution

func LargestRectangleArea(heights []int) int {
	type entry struct{ index, height int }
	stack := []entry{}
	maxArea := 0
	for i, h := range heights {
		start := i
		for len(stack) > 0 && stack[len(stack)-1].height > h {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			area := top.height * (i - top.index)
			if area > maxArea {
				maxArea = area
			}
			start = top.index
		}
		stack = append(stack, entry{start, h})
	}
	for _, e := range stack {
		area := e.height * (len(heights) - e.index)
		if area > maxArea {
			maxArea = area
		}
	}
	return maxArea
}
