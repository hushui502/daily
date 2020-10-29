package main

func maxAreaOfIsland(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	res := 0
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			if grid[r][c] == 1 {
				a := area(grid, r, c)
				res = max(a, res)
			}
		}
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func area(grid [][]int, r int, c int) int {
	if !(0 <= r && r < len(grid) && 0 <= c && c < len(grid[0])) {
		return 0
	}
	if grid[r][c] != 1 {
		return 0
	}
	grid[r][c] = 2
	return 1 + area(grid, r-1, c) + area(grid, r+1, c) + area(grid, r, c-1) + area(grid, r, c+1)
}
