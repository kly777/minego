// pkg/solver/solver.go
package solver

import (
	"image"
	"minego/internal/cell"
)

// solver 扫雷求解器
type solver struct {
	grid [][]cell.GridCell
}

func NewSolver(grid [][]cell.GridCell) *solver {
	return &solver{grid: grid}
}

// Solve 实现扫雷求解逻辑
func (s *solver) Solve() ([]image.Point, []image.Point) {
	var safePoints []image.Point
	var minePoints []image.Point

	rows := len(s.grid)
	if rows == 0 {
		return safePoints, minePoints
	}
	cols := len(s.grid[0])

	// 遍历所有单元格
	for i := range rows {
		for j := range cols {
			ccell := s.grid[i][j]

			// 只处理已打开的数字单元格(1-8)
			if ccell.State >= cell.Number1 && ccell.State <= cell.Number8 {
				// 获取周围单元格
				neighbors := s.getNeighbors(i, j)

				unknownCount := 0
				for _, neighbor := range neighbors {
					if neighbor.State == cell.Unknown || neighbor.State == cell.Flagged {
						unknownCount++
					}
				}
				if unknownCount == int(ccell.State) {
					for _, nb := range neighbors {
						if nb.State == cell.Unknown && !contains(minePoints, nb.Position) {
							minePoints = append(minePoints, image.Point{
								X: nb.Position.X,
								Y: nb.Position.Y,
							})
						}
					}
				}

				// 统计周围标记的地雷数量
				flaggedCount := 0
				for _, nb := range neighbors {
					if nb.State == cell.Flagged {
						flaggedCount++
					}
				}

				// 如果标记数等于当前单元格数字
				if flaggedCount == int(ccell.State) {
					for _, nb := range neighbors {
						// 找到未打开的安全单元格
						if nb.State == cell.Unknown&&!contains(safePoints, nb.Position) {
							safePoints = append(safePoints, image.Point{
								X: nb.Position.X,
								Y: nb.Position.Y,
							})
						}
					}
				}
			}
		}
	}

	return safePoints, minePoints
}

func contains(points []image.Point, p image.Point) bool {
	for _, pt := range points {
		if pt == p {
			return true
		}
	}
	return false
}

// getNeighbors 获取周围8个方向的单元格
func (s *solver) getNeighbors(row, col int) []cell.GridCell {
	var neighbors []cell.GridCell
	rows := len(s.grid)
	cols := len(s.grid[0])

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue // 跳过自身
			}

			r, c := row+i, col+j
			if r >= 0 && r < rows && c >= 0 && c < cols {
				neighbors = append(neighbors, s.grid[r][c])
			}
		}
	}
	return neighbors
}
