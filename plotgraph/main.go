package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"golang.org/x/tools/benchmark/parse"
)

const (
	BENCH_FILE = "./plotgraph/files/expIncreaseGeneralCloudFlare2.txt"
)

func main() {
	wordsPath, err := filepath.Abs(BENCH_FILE)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(wordsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	p := plot.New()

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	points := make(map[string]plotter.XYs)
	for scanner.Scan() {
		line := scanner.Text()
		b, err := parse.ParseLine(line)
		if err != nil {
			continue
		}
		strings.Split(b.Name, "")

		group := strings.Split(b.Name, "/")[1]
		groupSplit := strings.Split(group, "_")
		num := strings.Split(strings.Split(groupSplit[len(groupSplit)-1], "-")[0], "#")[0]

		name := strings.Join(groupSplit[:len(groupSplit)-1], "_")
		flNum, _ := strconv.ParseFloat(num, 64)

		if flNum >= 10 && flNum <= 10000 {
			points[name] = append(points[name], plotter.XY{X: flNum, Y: b.NsPerOp})
		}

	}
	fmt.Println(points)

	var args []interface{}
	for name, numVal := range points {
		args = append(args, name)
		args = append(args, numVal)
	}
	plotutil.AddLinePoints(p, args...)
	p.Legend.Top = true
	p.Legend.Left = true

	if err := p.Save(8*vg.Inch, 8*vg.Inch, "/home/pedroegs/expIncreaseGeneralCloudFlareGraph2_full.png"); err != nil {
		panic(err)
	}

}
