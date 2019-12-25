package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
)

type Point struct {
	x, y, z int
}

type Maze struct {
	grid   [][]rune
	width  int
	height int
	tele   map[Point]Point
	start  Point
	end    Point
}

func path(c rune) bool {
	return c == '.'
}

func wall(c rune) bool {
	return c == '#'
}

func maze(c rune) bool {
	return path(c) || wall(c)
}

func teleMark(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

func NewMaze(grid [][]rune) *Maze {
	height := len(grid)
	width := len(grid[0])
	midHeight := height / 2
	midWidth := width / 2
	ringHeight := -1
	ringWidth := -1
	for i := 2; i < midHeight; i++ {
		if !maze(grid[i][midWidth]) {
			ringHeight = i - 2
			break
		}
	}
	for i := 2; i < midWidth; i++ {
		if !maze(grid[midHeight][i]) {
			ringWidth = i - 2
			break
		}
	}
	if ringWidth < 0 || ringHeight < 0 {
		log.Fatalln("Invalid map")
	}

	teleStart := make(map[Point]Point)
	teleEnd := make(map[string]Point)
	for i := 0; i < width; i++ {
		if teleMark(grid[0][i]) {
			p := Point{i, 2, 0}
			mark := string([]rune{grid[0][i], grid[1][i]})
			if v, ok := teleEnd[mark]; ok {
				teleStart[p] = v
				teleStart[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if teleMark(grid[2+ringHeight][i]) {
			p := Point{i, 1 + ringHeight, 0}
			mark := string([]rune{grid[2+ringHeight][i], grid[3+ringHeight][i]})
			if v, ok := teleEnd[mark]; ok {
				teleStart[p] = v
				teleStart[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if teleMark(grid[height-2][i]) {
			p := Point{i, height - 3, 0}
			mark := string([]rune{grid[height-2][i], grid[height-1][i]})
			if v, ok := teleEnd[mark]; ok {
				teleStart[p] = v
				teleStart[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if teleMark(grid[height-ringHeight-4][i]) {
			p := Point{i, height - ringHeight - 2, 0}
			mark := string([]rune{grid[height-ringHeight-4][i], grid[height-ringHeight-3][i]})
			if v, ok := teleEnd[mark]; ok {
				teleStart[p] = v
				teleStart[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
	}
	for i := 0; i < height; i++ {
		if teleMark(grid[i][0]) {
			p := Point{2, i, 0}
			mark := string([]rune{grid[i][0], grid[i][1]})
			if v, ok := teleEnd[mark]; ok {
				teleStart[p] = v
				teleStart[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if teleMark(grid[i][2+ringWidth]) {
			p := Point{1 + ringWidth, i, 0}
			mark := string([]rune{grid[i][2+ringWidth], grid[i][3+ringWidth]})
			if v, ok := teleEnd[mark]; ok {
				teleStart[p] = v
				teleStart[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if teleMark(grid[i][width-2]) {
			p := Point{width - 3, i, 0}
			mark := string([]rune{grid[i][width-2], grid[i][width-1]})
			if v, ok := teleEnd[mark]; ok {
				teleStart[p] = v
				teleStart[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if teleMark(grid[i][width-ringWidth-4]) {
			p := Point{width - ringWidth - 2, i, 0}
			mark := string([]rune{grid[i][width-ringWidth-4], grid[i][width-ringWidth-3]})
			if v, ok := teleEnd[mark]; ok {
				teleStart[p] = v
				teleStart[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
	}

	start, ok := teleEnd["AA"]
	if !ok {
		log.Fatalln("No maze start")
	}
	end, ok := teleEnd["ZZ"]
	if !ok {
		log.Fatalln("No maze end")
	}

	return &Maze{
		grid:   grid,
		height: height,
		width:  width,
		tele:   teleStart,
		start:  Point{start.x, start.y, 0},
		end:    Point{end.x, end.y, 0},
	}
}

type Item struct {
	value Point
	g     int
	index int
}

type PriorityQueue []*Item

type OpenSet struct {
	pq     PriorityQueue
	valSet map[Point]Point
}

type ClosedSet map[Point]Point

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].g < pq[j].g
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func NewOpenSet() *OpenSet {
	return &OpenSet{
		pq:     PriorityQueue{},
		valSet: map[Point]Point{},
	}
}

func (os *OpenSet) Empty() bool {
	return os.pq.Len() == 0
}

func (os *OpenSet) Has(val Point) bool {
	_, ok := os.valSet[val]
	return ok
}

func (os *OpenSet) Push(value Point, g int) {
	item := &Item{
		value: value,
		g:     g,
	}
	heap.Push(&os.pq, item)
	os.valSet[value] = Point{}
}

func (os *OpenSet) Pop() (Point, int) {
	item := heap.Pop(&os.pq).(*Item)
	delete(os.valSet, item.value)
	return item.value, item.g
}

func NewClosedSet() ClosedSet {
	return ClosedSet{}
}

func (cs ClosedSet) Has(val Point) bool {
	_, ok := cs[val]
	return ok
}

func (cs ClosedSet) Push(val Point) {
	cs[val] = Point{}
}

func (m *Maze) isPath(pos Point) bool {
	return path(m.grid[pos.y][pos.x])
}

func (m *Maze) neighbors(pos Point) []Point {
	points := make([]Point, 0, 5)
	if k := (Point{pos.x, pos.y - 1, 0}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x, pos.y + 1, 0}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x - 1, pos.y, 0}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x + 1, pos.y, 0}); m.isPath(k) {
		points = append(points, k)
	}
	if v, ok := m.tele[Point{pos.x, pos.y, 0}]; ok {
		points = append(points, Point{v.x, v.y, 0})
	}
	return points
}

func (m *Maze) BFS() int {
	openSet := NewOpenSet()
	openSet.Push(m.start, 0)
	closedSet := NewClosedSet()
	for !openSet.Empty() {
		cur, curg := openSet.Pop()
		closedSet.Push(cur)
		if cur == m.end {
			return curg
		}

		for _, neighbor := range m.neighbors(cur) {
			if closedSet.Has(neighbor) || openSet.Has(neighbor) {
				continue
			}
			openSet.Push(neighbor, curg+1)
		}
	}
	return -1
}

func (m *Maze) isOuter(pos Point) bool {
	if pos.x == 2 || pos.y == 2 {
		return true
	}
	if pos.x == m.width-3 || pos.y == m.height-3 {
		return true
	}
	return false
}

func (m *Maze) neighbors2(pos Point) []Point {
	points := make([]Point, 0, 5)
	if k := (Point{pos.x, pos.y - 1, pos.z}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x, pos.y + 1, pos.z}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x - 1, pos.y, pos.z}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x + 1, pos.y, pos.z}); m.isPath(k) {
		points = append(points, k)
	}
	if v, ok := m.tele[Point{pos.x, pos.y, 0}]; ok {
		if m.isOuter(pos) {
			if pos.z > 0 {
				k := Point{v.x, v.y, pos.z - 1}
				points = append(points, k)
			} else {
			}
		} else {
			k := Point{v.x, v.y, pos.z + 1}
			points = append(points, k)
		}
	}
	return points
}

func (m *Maze) BFS2() int {
	openSet := NewOpenSet()
	openSet.Push(m.start, 0)
	closedSet := NewClosedSet()
	for !openSet.Empty() {
		cur, curg := openSet.Pop()
		closedSet.Push(cur)
		if cur == m.end {
			return curg
		}

		for _, neighbor := range m.neighbors2(cur) {
			if closedSet.Has(neighbor) || openSet.Has(neighbor) {
				continue
			}
			openSet.Push(neighbor, curg+1)
		}
	}
	return -1
}

func (m *Maze) Print() {
	for _, i := range m.grid {
		fmt.Println(string(i))
	}
}

func main() {
	grid := [][]rune{}
	{
		argsWithoutProg := os.Args[1:]
		file, err := os.Open(argsWithoutProg[0])
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			grid = append(grid, []rune(scanner.Text()))
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	m := NewMaze(grid)
	fmt.Println(m.BFS())
	fmt.Println(m.BFS2())
}
