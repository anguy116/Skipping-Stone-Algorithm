// GoSolution
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type matrix struct {
	height, length     int
	cost               [][]int
	demand, supply     []int
	warehouse, factory []string
}

type Cell struct {
	x, y int
	item int //qty in initial solution
	cost int //cost of each route

}

//A combination structure between the cost and quantity of item matrix

type WorkingMatrix struct {
	height, length     int
	matrix             [][]Cell
	demand, supply     []int
	warehouse, factory []string
}

//Structure for closed path found in the marginal cost method
//List of tuples of the coordinates from the closed path
//Cost of the route
type PathSolution struct {
	path [][2]int
	cost int
}

//Makes the matrix of cells, containing both cost of each cell as well as quantity in the cell
func newWorkMatrix(initialMatrix *matrix, costMatrix *matrix) (m WorkingMatrix) {

	height := costMatrix.height
	length := costMatrix.length

	m.height = height
	m.length = length

	matrix := make([][]Cell, height-2)

	for i := 0; i < height-2; i++ {
		matrix[i] = make([]Cell, length-2)
		for j := 0; j < length-2; j++ {
			//adding cost
			matrix[i][j].cost = costMatrix.cost[i][j]
			//adding quantity
			matrix[i][j].item = initialMatrix.cost[i][j]

			//setting the coordinates
			matrix[i][j].x = j
			matrix[i][j].y = i

		}
	}

	m.matrix = matrix

	m.supply = initialMatrix.supply
	m.demand = initialMatrix.demand
	m.warehouse = initialMatrix.warehouse
	m.factory = initialMatrix.factory

	return m

}

func marginalCost(matrix *WorkingMatrix, cell Cell, path chan PathSolution) {

	fmt.Println("Finding path for", cell.x, cell.y)

	//Finds the path of the empty cell
	solutionPath := findPath(matrix.matrix[cell.y][cell.x], matrix)

	//Calculate the cost of the solution

	sum := 0

	for i, coord := range solutionPath {
		x := coord[0]
		y := coord[1]
		if i%2 == 0 {
			sum += (matrix.matrix[y][x].cost) * 1

		} else {
			sum += (matrix.matrix[y][x].cost) * -1
		}

	}

	if solutionPath != nil {

		var entry PathSolution

		entry.cost = sum
		entry.path = solutionPath

		fmt.Printf("Finished path: (Cost: %v) ", sum)

		for i, coord := range solutionPath {
			x := coord[0]
			y := coord[1]
			if i%2 == 0 {
				fmt.Printf("+ Cell [%v,%v] %v", x, y, matrix.matrix[y][x].cost)

			} else {
				fmt.Printf("- Cell [%v,%v] %v", x, y, matrix.matrix[y][x].cost)
			}

		}
		fmt.Println()
		fmt.Println()

		path <- entry
	}
}

func findPath(cell Cell, matrix *WorkingMatrix) (path [][2]int) {

	var list [][2]int

	path = findPathH(cell, list, matrix)

	return path

}

func findPathH(cell Cell, list [][2]int, matrix *WorkingMatrix) (solution [][2]int) {

	list = append(list, [2]int{cell.x, cell.y})

	for i := 0; i < matrix.length-2; i++ {

		if i != cell.x && matrix.matrix[cell.y][i].item > 0 && !contains(list, [2]int{matrix.matrix[cell.y][i].x, matrix.matrix[cell.y][i].y}) {

			newpath := findPathV(matrix.matrix[cell.y][i], list, matrix) //finish thi call

			if newpath != nil {
				return newpath
			}
		}

	}

	return nil

}

func findPathV(cell Cell, list [][2]int, matrix *WorkingMatrix) (solution [][2]int) {
	list = append(list, [2]int{cell.x, cell.y})

	if list[0][0] == cell.x {
		return list
	}

	for i := 0; i < matrix.height-2; i++ {
		if i != cell.y && matrix.matrix[i][cell.x].item > 0 && !contains(list, [2]int{matrix.matrix[i][cell.x].x, matrix.matrix[i][cell.x].y}) {

			newpath := findPathH(matrix.matrix[i][cell.x], list, matrix)

			if newpath != nil {
				return newpath
			}
		}
	}

	return nil
}

func contains(path [][2]int, e [2]int) bool {
	for _, a := range path {
		if a[0] == e[0] && a[1] == e[1] {
			return true
		}
	}
	return false
}

func ReadFile(fileName string) (m matrix) {

	var cost [][]int

	file, err := os.Open(fileName) // just pass the file name
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	height := 0
	length := 0

	for scanner.Scan() {
		var s string
		s = scanner.Text()
		temp := strings.Split(s, " ")
		if len(temp) > length {
			length = len(temp)
		}

		height += 1
	}

	//Scanns the file for height and length of the whole matrix
	m.height = height
	m.length = length
	input, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Println(err)
	}

	//creates the matrix
	cost = make([][]int, height-1)

	warehouse := strings.Split(string(input), " ")

	warehouseList := make([]string, length-2)

	for k := 1; k < length-2; k++ {
		warehouseList[k] = warehouse[k]
	}

	m.warehouse = warehouseList

	factoryrow := strings.Split(string(input), "\n")

	factoryList := make([]string, height-1)
	for i := 1; i < height-1; i++ {

		factorycol := strings.Split(string(factoryrow[i]), " ")
		factoryList[i] = factorycol[0]
	}

	m.factory = factoryList

	row := strings.Split(string(input), "\n")
	_ = row
	for i := 1; i < height; i++ {
		cost[i-1] = make([]int, length-1)
		col := strings.Split(string(row[i]), " ")
		_ = col
		for k := 1; k < length; k++ {
			//checks if the index of the for loop isn't on the last element
			if i != height-1 || k != length-1 {
				entry, err := strconv.Atoi(col[k])
				if err == nil {
					//enters the numbers
					cost[i-1][k-1] = entry
				}
			}

		}
	}
	cost[height-2][length-2] = 0
	//makes the n by m matrix
	cutCost := make([][]int, height-2)
	for i := 0; i < height-2; i++ {
		cutCost[i] = make([]int, length-2)
		for j := 0; j < length-2; j++ {
			cutCost[i][j] = cost[i][j]

		}
	}

	m.cost = cutCost
	//makes the demand matrix
	demand := make([]int, length-2)

	for i := 0; i < length-2; i++ {
		demand[i] = cost[height-2][i]
	}
	m.demand = demand
	//makes the supply matrix
	supply := make([]int, height-2)
	for i := 0; i < height-2; i++ {
		supply[i] = cost[i][length-2]
	}
	m.supply = supply

	return m

}

func reallocate(minimumCost PathSolution, matrix *WorkingMatrix) (result *WorkingMatrix) {

	minimumQty := matrix.matrix[minimumCost.path[1][1]][minimumCost.path[1][0]].item

	for i, coord := range minimumCost.path {
		x := coord[0]
		y := coord[1]
		if i != 0 && i%2 == 1 {
			if matrix.matrix[y][x].item < minimumQty {
				minimumQty = matrix.matrix[y][x].item
			}
		}

	}

	for i, coord := range minimumCost.path {
		x := coord[0]
		y := coord[1]
		if i%2 == 0 {
			matrix.matrix[y][x].item = matrix.matrix[y][x].item + minimumQty

		} else {
			matrix.matrix[y][x].item = matrix.matrix[y][x].item - minimumQty

		}

	}

	result = matrix
	return result

}

func matrixValidity(initialMatrix *matrix, costMatrix *matrix) bool {
	if initialMatrix.length != costMatrix.length || initialMatrix.height != costMatrix.height {
		fmt.Println("The two given matrices are not identical")
		return false
	}

	count := 0

	for i := 0; i < initialMatrix.height-2; i++ {
		for j := 0; j < initialMatrix.length-2; j++ {

			if initialMatrix.cost[i][j] > 0 {
				count++
			}

		}

	}
	if count > 0 {
		return true
	}
	return false
}

func main() {

	fmt.Println("Enter Cost Matrix file: ")
	var cost string
	fmt.Scanln(&cost)

	costm := ReadFile(cost)
	_ = costm

	fmt.Println("Enter Initial Matrix file: ")
	var init string
	fmt.Scanln(&init)

	initial := ReadFile(init)
	_ = initial

	if matrixValidity(&initial, &costm) == true {

		x := newWorkMatrix(&initial, &costm)
		_ = x

		//Prints off initial values of the matrix
		fmt.Println("Costs:")

		for i := 0; i < x.height-2; i++ {
			for j := 0; j < x.length-2; j++ {
				fmt.Print(x.matrix[i][j].cost, " ")
			}
			fmt.Println()
		}

		fmt.Println("Supply: ")

		for i := 0; i < x.height-2; i++ {
			fmt.Print(x.supply[i], " ")
		}

		fmt.Println("\nDemand: ")

		for j := 0; j < x.length-2; j++ {
			fmt.Print(x.demand[j], " ")
		}

		fmt.Println("\nUnits: ")

		for i := 0; i < x.height-2; i++ {
			for j := 0; j < x.length-2; j++ {
				fmt.Print(x.matrix[i][j].item, " ")
			}
			fmt.Println()
		}

		empty := 0

		for i := 0; i < x.height-2; i++ {
			for j := 0; j < x.length-2; j++ {
				if x.matrix[i][j].item == 0 {

					empty += 1
				}
			}
		}

		path := make(chan PathSolution, empty)

		var minimumCost PathSolution
		minimumCost.cost = 0

		iteration := 1

		fmt.Println("Iteration ", iteration)

		for i := 0; i < x.height-2; i++ {
			for j := 0; j < x.length-2; j++ {
				if x.matrix[i][j].item == 0 {

					go marginalCost(&x, x.matrix[i][j], path)

					currentSolution := <-path

					if currentSolution.cost < minimumCost.cost {
						minimumCost = currentSolution
					}

					time.Sleep(1 * 1e9)

				}
			}
		}

		result := reallocate(minimumCost, &x)
		_ = result

		iteration++

		fmt.Println("Iteration ", iteration)

		for {

			empty := 0

			for i := 0; i < result.height-2; i++ {
				for j := 0; j < result.length-2; j++ {
					if result.matrix[i][j].item == 0 {

						empty += 1
					}
				}
			}

			path := make(chan PathSolution, empty)

			var minimumCost PathSolution
			minimumCost.cost = 0

			for i := 0; i < result.height-2; i++ {
				for j := 0; j < result.length-2; j++ {
					if result.matrix[i][j].item == 0 {

						go marginalCost(&x, x.matrix[i][j], path)

						currentSolution := <-path

						if currentSolution.cost < minimumCost.cost {
							minimumCost = currentSolution
						}

						time.Sleep(1 * 1e9)

					}
				}
			}

			if minimumCost.cost == 0 {
				break
			} else {
				iteration++
				result = reallocate(minimumCost, result)
			}

		}

		fmt.Println("Costs:")

		for i := 0; i < result.height-2; i++ {
			for j := 0; j < result.length-2; j++ {
				fmt.Print(result.matrix[i][j].cost, " ")
			}
			fmt.Println()
		}

		fmt.Println("\nUnits: ")

		for i := 0; i < result.height-2; i++ {
			for j := 0; j < result.length-2; j++ {
				fmt.Print(result.matrix[i][j].item, " ")
			}
			fmt.Println()
		}

		tcost := 0

		for i := 0; i < result.height-2; i++ {
			for j := 0; j < result.length-2; j++ {
				tcost += result.matrix[i][j].cost * result.matrix[i][j].item

			}

		}

		fmt.Println("Total Cost: ", tcost)

	}
}
