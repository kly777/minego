package solver

import (
	"fmt"
)

// Equation 表示一个方程：指定变量的和等于目标值
type Equation struct {
	Indices []int // 变量索引列表
	Sum     int   // 方程的目标和
}

// solveBinaryEquations 求解二进制方程组
func solveBinaryEquations(n int, equations []Equation) [][]int {
	solutions := [][]int{}
	fmt.Println("n", n)
	totalStates := 1 << n // 所有可能的状态数 (2^n)
	fmt.Println("totalStates", totalStates)
	// 枚举所有可能的二进制状态
	for state := range totalStates {
		valid := true

		// 检查是否满足所有方程
		for _, eq := range equations {
			sum := 0
			for _, idx := range eq.Indices {
				if idx < 0 || idx >= n {
					panic(fmt.Sprintf("变量索引 %d 超出范围 [0, %d]", idx, n-1))
				}
				// 使用位运算检查该变量是否为1
				if state&(1<<idx) != 0 {
					sum++
				}
			}

			if sum != eq.Sum {
				valid = false
				break
			}
		}

		// 如果满足所有方程，保存解
		if valid {
			solution := make([]int, n)
			for i := 0; i < n; i++ {
				if state&(1<<i) != 0 {
					solution[i] = 1
				} // 默认为0
			}
			solutions = append(solutions, solution)
		}
	}

	return solutions
}

// func main() {
// 	// 示例1：无解的情况
// 	// 方程:
// 	//   x0 + x1 = 1
// 	//   x0 + x2 = 1
// 	//   x1 + x2 = 1
// 	n1 := 3
// 	equations1 := []Equation{
// 		{Indices: []int{0, 1}, Sum: 1},
// 		{Indices: []int{0, 2}, Sum: 1},
// 		{Indices: []int{1, 2}, Sum: 1},
// 	}

// 	// 示例2：两个解的情况
// 	// 方程:
// 	//   x0 + x1 + x2 = 2
// 	//   x0 + x1 = 1
// 	//   x2 + x3 = 1
// 	n2 := 4
// 	equations2 := []Equation{
// 		{Indices: []int{0, 1, 2}, Sum: 2},
// 		{Indices: []int{0, 1}, Sum: 1},
// 		{Indices: []int{2, 3}, Sum: 1},
// 	}

// 	// 求解示例1
// 	fmt.Println("示例1 求解:")
// 	solutions1 := solveBinaryEquations(n1, equations1)
// 	if len(solutions1) == 0 {
// 		fmt.Println("  无解")
// 	} else {
// 		for i, sol := range solutions1 {
// 			fmt.Printf("  解%d: %v\n", i+1, sol)
// 		}
// 	}

// 	// 求解示例2
// 	fmt.Println("\n示例2 求解:")
// 	solutions2 := solveBinaryEquations(n2, equations2)
// 	if len(solutions2) == 0 {
// 		fmt.Println("  无解")
// 	} else {
// 		for i, sol := range solutions2 {
// 			fmt.Printf("  解%d: %v\n", i+1, sol)
// 		}
// 	}

// 	// === 用户自定义方程组 ===
// 	fmt.Println("\n=== 自定义方程组求解 ===")

// 	// 输入变量数量
// 	var n int
// 	fmt.Print("输入变量数量: ")
// 	_, err := fmt.Scan(&n)
// 	if err != nil || n < 1 {
// 		fmt.Println("无效输入，变量数应为正整数")
// 		return
// 	}

// 	// 输入方程数量
// 	var m int
// 	fmt.Print("输入方程数量: ")
// 	_, err = fmt.Scan(&m)
// 	if err != nil || m < 1 {
// 		fmt.Println("无效输入，方程数应为正整数")
// 		return
// 	}

// 	// 收集方程
// 	equations := make([]Equation, m)
// 	for i := 0; i < m; i++ {
// 		fmt.Printf("\n方程 #%d:\n", i+1)

// 		// 输入变量数量
// 		var k int
// 		fmt.Print("  该方程涉及的变量数量: ")
// 		_, err := fmt.Scan(&k)
// 		if err != nil || k < 1 {
// 			fmt.Println("  无效输入，变量数应为正整数")
// 			return
// 		}

// 		// 输入变量索引
// 		indices := make([]int, k)
// 		fmt.Printf("  输入 %d 个变量索引 (0-%d): ", k, n-1)
// 		for j := 0; j < k; j++ {
// 			_, err := fmt.Scan(&indices[j])
// 			if err != nil || indices[j] < 0 || indices[j] >= n {
// 				fmt.Printf("  无效索引，应在 [0, %d] 范围内\n", n-1)
// 				return
// 			}
// 		}

// 		// 输入目标和
// 		var sum int
// 		fmt.Print("  输入目标和: ")
// 		_, err = fmt.Scan(&sum)
// 		if err != nil || sum < 0 {
// 			fmt.Println("  无效输入，和应为非负整数")
// 			return
// 		}

// 		equations[i] = Equation{Indices: indices, Sum: sum}
// 	}

// 	// 求解并输出结果
// 	fmt.Println("\n求解结果:")
// 	solutions := solveBinaryEquations(n, equations)
// 	if len(solutions) == 0 {
// 		fmt.Println("  无解")
// 	} else {
// 		fmt.Printf("  找到 %d 个解:\n", len(solutions))
// 		for i, sol := range solutions {
// 			fmt.Printf("  解%d: %v\n", i+1, sol)
// 		}
// 	}
// }
