package main

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/fogleman/gg"

	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
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

func (f *Forest) shootLightning(toVizualize bool) bool {
	x := rand.Intn(f.rows)
	y := rand.Intn(f.columns)
	if f.forest[x][y] == 0 {
		fmt.Printf("Piorun nie uderzył w drzewo! Nie ma pożaru %v\n", emoji.SmilingFaceWithSunglasses)
		return false
	} else {
		fmt.Printf("Uderzenie pioruna w drzewo (%v, %v)!! Zaczyna się pożar! %v\n", y, x, emoji.AstonishedFace)
		f.forest[x][y] = 2
		if toVizualize {
			f.saveAsImage("forest_pics/forest_1.gif")
		}
		f.spreadFire(x, y, toVizualize)
		return true
	}
}

// z modyfikacją z wiatrem trzeba zmienić tylko neighbours
func (f *Forest) spreadFire(x, y int, toVizualize bool) {
	neighbors := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

	for _, offset := range neighbors {
		newX, newY := x+offset[0], y+offset[1]
		if newX >= 0 && newX < f.rows && newY >= 0 && newY < f.columns && f.forest[newX][newY] == 1 {
			f.forest[newX][newY] = 2
			if toVizualize {
				var fireNumber = f.countFire()
				if fireNumber % 5 == 0 {
					filename := fmt.Sprintf("forest_pics/forest_%d.gif", f.countFire())
					f.saveAsImage(filename)
				}
			}
			f.spreadFire(newX, newY, toVizualize)
		}
	}
}

func (f *Forest) getStatistics() float32{
	var burned = f.countFire()
	return float32(burned)/float32(f.treeNumber)
}

func (f *Forest) countFire() int {
	var burned int = 0
	for i := 0; i < f.rows; i++  {
		for _, place := range f.forest[i] {
			if place == 2 {
				burned++
			}
		}
	}
	return burned
}

func (f *Forest) saveAsImage(filename string) error {
    const cellSize = 10
    width := f.columns * cellSize
    height := f.rows * cellSize

    dc := gg.NewContext(width, height)

    for row := 0; row < f.rows; row++ {
        for col := 0; col < f.columns; col++ {
            switch f.forest[row][col] {
            case 0:
                dc.SetHexColor("#ffffff") // White color for empty space
            case 1:
                dc.SetHexColor("#1f8c1f") // Green color for trees
            case 2:
                dc.SetHexColor("#ff0000") // Red color for fire
            }
            dc.DrawRectangle(float64(col*cellSize), float64(row*cellSize), float64(cellSize), float64(cellSize))
            dc.Fill()
        }
    }

    img := dc.Image()

    paletted := image.NewPaletted(img.Bounds(), palette.Plan9)
    draw.FloydSteinberg.Draw(paletted, img.Bounds(), img, image.ZP)

    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    if err := gif.Encode(file, paletted, nil); err != nil {
        return err
    }

    return nil
}

func createGIFFromImages() error {
    folderPath := "forest_pics"
    files, err := os.ReadDir(folderPath)
    if err != nil {
        return err
    }

    sort.Slice(files, func(i, j int) bool {
        num1, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(files[i].Name(), "forest_"), ".gif"))
        num2, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(files[j].Name(), "forest_"), ".gif"))
        return num1 < num2
    })

    outGif := &gif.GIF{}
    for _, file := range files {
        f, err := os.Open(filepath.Join(folderPath, file.Name()))
        if err != nil {
            return err
        }
        defer f.Close()
        inGif, _, err := image.Decode(f)
        if err != nil {
            return err
        }

        outGif.Image = append(outGif.Image, inGif.(*image.Paletted))
        outGif.Delay = append(outGif.Delay, 0)
    }

    f, err := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        return err
    }
    defer f.Close()

    if err := gif.EncodeAll(f, outGif); err != nil {
        return err
    }

    return nil
}

func runSimulationWithVisualization(rows, columns int, percentage float32) {
	f := newForest(rows, columns, percentage)
	err := f.saveAsImage("forest_pics/forest_0.gif")
	if err != nil {
		fmt.Println("Błąd podczas zapisywania obrazu:", err)
	}
	isBurned := f.shootLightning(true)
	for !isBurned{
		isBurned = f.shootLightning(true)
	}
	f.saveAsImage("forest_pics/forest_10000000000000000000000000000.gif")
	fmt.Println("Tworzymy gifa!")
	err1 := createGIFFromImages()
	if err1 != nil {
		println("Błąd tworzenia GIF:", err1)
	}

}

func main() {
	var rows int = 40
	var columns int = 40
	var percentage float32 = 0.5
	// var toVizualize bool = true
	runSimulationWithVisualization(rows, columns, percentage)


	// f := newForest(rows, columns, percentage)
	// fmt.Printf("Las o wymiarach %v na %v przed uderzeniem pioruna:\n", rows, columns)
	// f.displayForest()

	// fmt.Println("\nWywołujemy uderzenie pioruna...")
	// isBurned := f.shootLightning(toVizualize)
	
	// f.displayForest()
	// var stats float32 = f.getStatistics()
	// fmt.Printf("Spalono %.2f%% drzew\n", stats*100)


}

// 1. Dopytać się o punkty za wizualizację

// Test do znalezienia optymalnego poziomu zalesienia:
// 1. dla każdego percentage od 5% do 95% z krokiem co 5 punktów procentowych wykonać po 1000 prób.
// 2. Dla każdego percentage obliczyć średni poziom spalenia lasu
// 3. Dla każdego percentage obliczyć stosunek: (zalesienie)/(średni poziom spalenia lasu)
// 3. Znaleźć największy (najlepszy) stosunek

// Test można przeprowadzić na lasach o różnych wymiarach, ale o tej samej powierzchni, np.
// powierzchnia 400: 20x20, 16x25, 10x40, 8x50, 5x80, 4x100, 2x200