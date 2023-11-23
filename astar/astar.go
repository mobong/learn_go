package main

import (
	"container/heap"
	"fmt"
	"math"
)

// 坐标
type Point struct {
	X, Y int
}

// 节点结构
type Node struct {
	Point
	Parent  *Node   // 父节点
	G, H, F float64 // G：从起点到当前节点的实际代价，H：从当前节点到目标节点的估计代价，F：总代价
}

// 曼哈顿距离
func heuristic(currentNode, targetNode *Node) float64 {
	return math.Abs(float64(currentNode.X-targetNode.X)) + math.Abs(float64(currentNode.Y-targetNode.Y))
}

// A*算法
func astar(startNode, targetNode *Node, obstacles [][]bool) []*Node {
	openList := make(priorityQueue, 0)
	heap.Init(&openList)

	startNode.G = 0
	startNode.H = heuristic(startNode, targetNode)
	startNode.F = startNode.G + startNode.H

	heap.Push(&openList, startNode)

	visited := make(map[Point]bool)
	visited[startNode.Point] = true
	for openList.Len() > 0 {
		currentNode := heap.Pop(&openList).(*Node)
		if currentNode.X == targetNode.X && currentNode.Y == targetNode.Y {
			// 找到路径，回溯父节点
			path := make([]*Node, 0)
			for currentNode != nil {
				path = append(path, currentNode)
				currentNode = currentNode.Parent
			}
			return path
		}

		// 获取当前节点的相邻节点
		neighbors := getNeighbors(currentNode, obstacles, visited)
		if len(neighbors) == 0 {
			break
		}
		for _, neighbor := range neighbors {
			// 计算相邻节点的G值和H值
			g := currentNode.G + 1 // 假设相邻节点间的代价为1
			h := heuristic(neighbor, targetNode)
			f := g + h

			if f < neighbor.F || !visited[neighbor.Point] {
				neighbor.G = g
				neighbor.H = h
				neighbor.F = f
				neighbor.Parent = currentNode

				if !visited[neighbor.Point] {
					heap.Push(&openList, neighbor)
					visited[neighbor.Point] = true
				}
			}
		}
	}
	// 无法找到路径
	return nil
}

// 获取当前节点的相邻节点
func getNeighbors(currentNode *Node, obstacles [][]bool, visited map[Point]bool) []*Node {
	neighbors := make([]*Node, 0)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			newX := currentNode.X + dx
			newY := currentNode.Y + dy
			if newX >= 0 && newX < len(obstacles) && newY >= 0 && newY < len(obstacles[0]) && !obstacles[newX][newY] {
				P := Point{X: newX, Y: newY}
				neighbor := &Node{
					Point:  P,
					Parent: nil,
					G:      0,
					H:      0,
					F:      0,
				}
				if !visited[P] {
					neighbors = append(neighbors, neighbor)
				}
			}
		}
	}
	return neighbors
}

// 反转路径
func reversePath(path []*Node) {
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
}

// 优先队列
type priorityQueue []*Node

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	node := x.(*Node)
	*pq = append(*pq, node)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}

func main() {
	// 创建地图和障碍物
	obstacles := [][]bool{
		{false, false, false, false},
		{false, true, false, true},
		{false, false, true, true},
		{false, false, false, false},
	}

	// 创建起点和终点
	startNode := &Node{Point{0, 0}, nil, 0, 0, 0}
	targetNode := &Node{Point{3, 3}, nil, 0, 0, 0}

	// 使用A*算法寻找路径
	path := astar(startNode, targetNode, obstacles)

	// 打印路径
	if path != nil {
		// 反转路径，使起点在前，终点在后
		reversePath(path)
		for _, node := range path {
			fmt.Printf("{%d, %d} ", node.X, node.Y)
		}
	} else {
		fmt.Println("无法找到路径")
	}
}
