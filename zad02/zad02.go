package main

import (
	"fmt"
	"math/rand"

	"github.com/enescakir/emoji"
)

type Forest struct {
	forest  [][]int
	rows    int
	columns int
	treeNumber int
}

func newForest(rows, columns int, percentage float32) *Forest {
	f := &Forest{
		forest:  make([][]int, rows),
		rows:    rows,
		columns: columns,
		treeNumber: int(float32(rows*columns) * percentage),
	}
	for i := range f.forest {
		f.forest[i] = make([]int, columns)
	}
	f.populateForest()
	return f
}

func (f *Forest) populateForest() {
	for i := 0; i < f.treeNumber; i++ {
		x := rand.Intn(f.rows)
		y := rand.Intn(f.columns)
		for f.forest[x][y] == 1 {
			x = rand.Intn(f.rows)
			y = rand.Intn(f.columns)
		}
		f.forest[x][y] = 1
	}
}

func (f *Forest) displayForest() {
	for row := 0; row < f.rows; row++ {
		for col := 0; col < f.columns; col++ {
			switch f.forest[row][col] {
			case 0:
				fmt.Printf("%v  ", emoji.Multiply)
			case 1:
				fmt.Printf("%v ", emoji.EvergreenTree)
			case 2:
				fmt.Printf("%v ", emoji.Fire)
			}
		}
		fmt.Println()
	}
}

func (f *Forest) shootLightning() bool {
	x := rand.Intn(f.rows)
	y := rand.Intn(f.columns)
	if f.forest[x][y] == 0 {
		fmt.Printf("Piorun nie uderzył w drzewo! Nie ma pożaru %v\n", emoji.SmilingFaceWithSunglasses)
		return false
	} else {
		fmt.Printf("Uderzenie pioruna w drzewo (%v, %v)!! Zaczyna się pożar! %v\n", y, x, emoji.AstonishedFace)
		f.forest[x][y] = 2
		f.spreadFire(x, y)
		return true
	}
}

// z modyfikacją z wiatrem trzeba zmienić tylko neighbours
func (f *Forest) spreadFire(x, y int) {
	neighbors := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

	for _, offset := range neighbors {
		newX, newY := x+offset[0], y+offset[1]
		if newX >= 0 && newX < f.rows && newY >= 0 && newY < f.columns && f.forest[newX][newY] == 1 {
			f.forest[newX][newY] = 2
			// f.DisplayForest()
			// fmt.Println()
			f.spreadFire(newX, newY)
		}
	}
}

func (f *Forest) getStatistics() float32{
	var burned int = 0
	for i := 0; i < f.rows; i++  {
		for _, place := range f.forest[i] {
			if place == 2 {
				burned++
			}
		}
	}
	return float32(burned)/float32(f.treeNumber)
}

func main() {
	var rows int = 10
	var columns int = 40
	var percentage float32 = 0.5

	f := newForest(rows, columns, percentage)
	fmt.Printf("Las o wymiarach %v na %v przed uderzeniem pioruna:\n", rows, columns)
	f.displayForest()

	fmt.Println("\nWywołujemy uderzenie pioruna...")
	f.shootLightning()
	f.displayForest()
	var stats float32 = f.getStatistics()
	fmt.Printf("Spalono %.2f%% drzew\n", stats*100)

}

// 1. Dopytać się o punkty za wizualizację

// Test do znalezienia optymalnego poziomu zalesienia:
// 1. dla każdego percentage od 5% do 95% z krokiem co 5 punktów procentowych wykonać po 1000 prób.
// 2. Dla każdego percentage obliczyć średni poziom spalenia lasu
// 3. Dla każdego percentage obliczyć stosunek: (zalesienie)/(średni poziom spalenia lasu)
// 3. Znaleźć największy (najlepszy) stosunek
// maks zalesienie, min straty

// Test można przeprowadzić na lasach o różnych wymiarach, ale o tej samej powierzchni, np.
// powierzchnia 400: 20x20, 16x25, 10x40, 8x50