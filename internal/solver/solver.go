// pkg/solver/solver.go
package solver

import (
	"fmt"
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

	// 遍历所有单元格(简单处理)
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
							nb.State = cell.Flagged
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
						if nb.State == cell.Unknown && !contains(safePoints, nb.Position) {
							safePoints = append(safePoints, image.Point{
								X: nb.Position.X,
								Y: nb.Position.Y,
							})
						}
					}
				} else if flaggedCount > int(ccell.State) {
					for _, nb := range neighbors {
						// 找到被标记的单元格
						if nb.State == cell.Flagged && !contains(minePoints, nb.Position) {
							minePoints = append(minePoints, image.Point{
								X: nb.Position.X,
								Y: nb.Position.Y,
							})
						}
					}
				}
			}
		}
	}
	pointID := NewPointIDMap()
	equations := make([]Equation, 0)
	n := 0
	for i := range rows {
		for j := range cols {
			ccell := s.grid[i][j]
			if ccell.State >= cell.Number1 && ccell.State <= cell.Number8 {
				neighbors := s.getNeighbors(i, j)
				unknowncells := make([]int, 0)
				for _, nb := range neighbors {
					if nb.State == cell.Unknown || nb.State == cell.Flagged {
						id, ok := pointID.GetID(nb.Position)
						if ok {
							unknowncells = append(unknowncells, id)
						} else {
							pointID.Add(nb.Position, n)
							unknowncells = append(unknowncells, n)
							n++
							if n >= 18 {
								goto end
							}
						}
					}
				}
				equations = append(equations, Equation{unknowncells, int(ccell.State)})
				fmt.Println(unknowncells, int(ccell.State))
			}
		}
	}
end:
	res := make([][]int, 0)
	if len(equations) > 10 {
		res = solveBinaryEquations(n, equations[:10])
	} else {
		res = solveBinaryEquations(n, equations)
	}

	fmt.Println("res", res)
	samep := comparePositions(res)
	fmt.Println("samp", samep)
	for id, p := range samep {
		switch p {
		case 0:
			if !contains(safePoints, pointID.idToPoint[id]) {
				safePoints = append(safePoints, pointID.idToPoint[id])
			}
		case 1:
			if !contains(minePoints, pointID.idToPoint[id])&&s.grid[pointID.idToPoint[id].Y][pointID.idToPoint[id].X].State != cell.Flagged {
				minePoints = append(minePoints, pointID.idToPoint[id])
			}
		}
	}

	return safePoints, minePoints
}

type PointIDMap struct {
	pointToID map[image.Point]int
	idToPoint map[int]image.Point
}

func NewPointIDMap() *PointIDMap {
	return &PointIDMap{
		pointToID: make(map[image.Point]int),
		idToPoint: make(map[int]image.Point),
	}
}

func (m *PointIDMap) Add(p image.Point, id int) {
	m.pointToID[p] = id
	m.idToPoint[id] = p
}

func (m *PointIDMap) GetID(p image.Point) (int, bool) {
	id, ok := m.pointToID[p]
	return id, ok
}

func (m *PointIDMap) GetPoint(id int) (image.Point, bool) {
	point, ok := m.idToPoint[id]
	return point, ok
}

func comparePositions(arr [][]int) []int {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return []int{}
	}

	rows := len(arr)
	if rows == 1 {
		return arr[0]
	}
	cols := len(arr[0])
	result := make([]int, cols)

	for col := range cols {
		val := arr[0][col]
		same := true

		for row := 1; row < rows; row++ {
			if arr[row][col] != val {
				same = false
				break
			}
		}

		if same {
			result[col] = val
		} else {
			result[col] = -1
		}
	}

	return result
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
