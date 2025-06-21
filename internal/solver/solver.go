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
	safeSet := make(map[image.Point]struct{})
	mineSet := make(map[image.Point]struct{})

	addSafe := func(p image.Point) {
		if _, exists := safeSet[p]; !exists {
			safePoints = append(safePoints, p)
			safeSet[p] = struct{}{}
		}
	}
	addMine := func(p image.Point) {
		if _, exists := mineSet[p]; !exists {
			minePoints = append(minePoints, p)
			mineSet[p] = struct{}{}
		}
	}

	rows := len(s.grid)
	if rows == 0 {
		return safePoints, minePoints
	}
	cols := len(s.grid[0])

	// 遍历所有单元格
	for i := range rows {
		for j := range cols {
			ccell := s.grid[i][j]
			if ccell.State < cell.Number1 || ccell.State > cell.Number8 {
				continue // 只处理数字单元格
			}

			neighbors := s.getNeighborPointers(i, j) // 返回指向网格单元格的指针切片
			unknownCount := 0
			flaggedCount := 0

			// 单次遍历统计未知和标记数量
			for _, nb := range neighbors {
				switch nb.State {
				case cell.Unknown, cell.Flagged:
					unknownCount++
					if nb.State == cell.Flagged {
						flaggedCount++
					}
				}
			}

			// 标记所有未知为地雷
			if unknownCount == int(ccell.State) {
				for _, nb := range neighbors {
					if nb.State == cell.Unknown {
						nb.State = cell.Flagged // 同步到原始网格
						addMine(nb.Position)
					}
				}
				continue
			}

			// 标记所有未知为安全
			if flaggedCount == int(ccell.State) {
				for _, nb := range neighbors {
					if nb.State == cell.Unknown {
						addSafe(nb.Position)
					}
				}
			} else if flaggedCount > int(ccell.State) {
				for _, nb := range neighbors {
					if nb.State == cell.Flagged {
						addMine(nb.Position)
					}
				}
			}
		}
	}

	// 构建方程组
	pointID := NewPointIDMap()
	equations := make([]Equation, 0)
	n := 0
Loop:
	for i := range rows {
		for j := range cols {
			flaggedCount := 0
			ccell := s.grid[i][j]
			if ccell.State < cell.Number1 || ccell.State > cell.Number8 {
				continue
			}
			neighbors := s.getNeighbors(i, j)
			s := 0
			for _, nb := range neighbors {
				if nb.State == cell.Flagged || nb.State == cell.Empty {
					s += 1
				}
			}
			if s == len(neighbors) {
				continue
			}
			unknowncells := make([]int, 0)
			for _, nb := range neighbors {
				if nb.State == cell.Flagged {
					flaggedCount += 1
				}
				if nb.State == cell.Unknown {
					id, ok := pointID.GetID(nb.Position)
					if !ok {
						pointID.Add(nb.Position, n)
						unknowncells = append(unknowncells, n)
						n++
						if n >= 21 {
							break Loop // 替代goto
						}
					} else {
						unknowncells = append(unknowncells, id)
					}
				}
			}
			if len(unknowncells) > 0 {
				equations = append(equations, Equation{unknowncells, int(ccell.State) - flaggedCount})
			}
		}
	}

	// 求解方程组
	res := make([][]int, 0)
	res = solveBinaryEquations(n, equations)
	// fmt.Println("res", res)
	// 处理结果
	samep := comparePositions(res)
	if usefulEle(samep) == 0 && len(safePoints) == 0 && len(minePoints) == 0 {
		if len(res) >= 1 && len(res[0]) >= 1 {
			samep = []int{res[0][0]}
		}

	}
	for id, p := range samep {
		point := pointID.idToPoint[id]
		switch p {
		case 0:
			addSafe(point)
		case 1:
			if s.grid[point.Y][point.X].State != cell.Flagged {
				addMine(point)
			}
		}
	}

	return safePoints, minePoints
}

func usefulEle(arr []int) int {
	count := 0
	for _, v := range arr {
		if v != -1 {
			count++
		}
	}
	return count
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

// getNeighborPointers 获取周围8个方向单元格的指针切片
func (s *solver) getNeighborPointers(row, col int) []*cell.GridCell {
	var neighbors []*cell.GridCell
	rows := len(s.grid)
	if rows == 0 {
		return neighbors
	}
	cols := len(s.grid[0])

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue // 跳过自身
			}

			r, c := row+i, col+j
			if r >= 0 && r < rows && c >= 0 && c < cols {
				neighbors = append(neighbors, &s.grid[r][c]) // 取地址
			}
		}
	}
	return neighbors
}
