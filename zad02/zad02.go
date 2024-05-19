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
	"github.com/wcharczuk/go-chart"

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

func newForest(rows, columns int, percentage float64) *Forest {
	f := &Forest{
		forest:  make([][]int, rows),
		rows:    rows,
		columns: columns,
		treeNumber: int(float64(rows*columns) * percentage),
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
		// fmt.Printf("Piorun nie uderzył w drzewo! Nie ma pożaru %v\n", emoji.SmilingFaceWithSunglasses)
		return false
	} else {
		// fmt.Printf("Uderzenie pioruna w drzewo (%v, %v)!! Zaczyna się pożar! %v\n", y, x, emoji.AstonishedFace)
		f.forest[x][y] = 2
		if toVizualize {
			f.saveAsImage("forest_pics/forest_1.gif")
		}
		f.spreadFire(x, y, toVizualize)
		return true
	}
}

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

func (f *Forest) getStatistics() float64{
	var burned = f.countFire()
	return float64(burned)/float64(f.treeNumber)
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
                dc.SetHexColor("#ffffff")
            case 1:
                dc.SetHexColor("#1f8c1f")
            case 2:
                dc.SetHexColor("#ff0000")
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

    f, err := os.OpenFile("visuals/out.gif", os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        return err
    }
    defer f.Close()

    if err := gif.EncodeAll(f, outGif); err != nil {
        return err
    }

    return nil
}

func runSimulationWithVisualization(rows, columns int, percentage float64) {
	f := newForest(rows, columns, percentage)
	f.displayForest()
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
	fmt.Println()
	f.displayForest()
	fmt.Println(f.getStatistics())

}


func findOptimum(m, n int, toPlot bool) (float64, float64) {
	var percentages []float64
	var maxRatio float64
	var bestPercentage float64

	var percentage float64 = 0.05
	for percentage < 1 {
		percentages = append(percentages, percentage)
		percentage += 0.05
	}

	var ratioData []float64
	for _, p := range percentages {
		sumBurned := float64(0)
		for i := 0; i < 5000; i++ {
			f := newForest(m, n, p)
			f.shootLightning(false)
			burnedPercentage := f.getStatistics()
			sumBurned += burnedPercentage
		}
		avgBurned := sumBurned / 5000
		ratio := p * (1 - avgBurned)
		ratioData = append(ratioData, ratio)

		if ratio > maxRatio {
			maxRatio = ratio
			bestPercentage = p
		}
	}

	fmt.Printf("Najlepsze ratio (%v, %v): %.2f, dla poziomu zalesienia: %.2f\n", m, n, maxRatio, bestPercentage)

	if toPlot {
		plotData(percentages, ratioData)
	}

	return maxRatio, bestPercentage
}

func plotData(percentages []float64, ratios []float64) error {
	series := chart.ContinuousSeries{
		Name:    "Ratio vs Zalesienie",
		XValues: percentages,
		YValues: ratios,
	}

	graph := chart.Chart{
		Series: []chart.Series{series},
	}

	graph.XAxis = chart.XAxis{
		Name:      "Zalesienie",
		NameStyle: chart.StyleShow(),
		Style:     chart.StyleShow(),
		Range:     &chart.ContinuousRange{Min: 0, Max: 1},
		ValueFormatter: func(v interface{}) string {
			if vf, isFloat := v.(float64); isFloat {
				return fmt.Sprintf("%.2f", vf)
			}
			return ""
		},
	}

	graph.YAxis = chart.YAxis{
		Name:      "Ratio",
		NameStyle: chart.StyleShow(),
		Style:     chart.StyleShow(),
		Range:     &chart.ContinuousRange{Min: 0, Max: 0.5},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	file, err := os.Create("visuals/plot.png")
	if err != nil {
		fmt.Println(err, "1")
		return err
	}
	defer file.Close()

	err = graph.Render(chart.PNG, file)
	if err != nil {
		fmt.Println(err, "2")
		return err
	}
	return nil
}

func main() {
	var rows int = 40
	var columns int = 40
	var percentage float64 = 0.5
	runSimulationWithVisualization(rows, columns, percentage)
	findOptimum(20, 20, true)
	// findOptimum(16, 25, false)
	// findOptimum(10, 40, false)
	// findOptimum(8, 50, false)
	// findOptimum(5, 80, false)
	// findOptimum(4, 100, false)
	// findOptimum(2, 200, false)

}
