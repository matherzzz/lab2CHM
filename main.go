package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"math"
	"os"
)

const (
	step = 101
	tau  = 1 / float64(step-1)
	y1   = 1
	y2   = 2
)

func generateLineItems(u [step][2]float64, a int) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(u); i++ {
		items = append(items, opts.LineData{Value: u[i][a]})
	}
	return items
}

func createLineChart(u [step][2]float64, s string) {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Lab2 in Go",
			Subtitle: "1 - red, 2 - blue",
		}))
	// Put data into instance
	var t [step]float64
	for i := 0; i < step; i++ {
		t[i] = math.Floor(float64(i)*tau*100) / 100
	}
	line.SetXAxis(t).
		AddSeries("Category A", generateLineItems(u, 0)).
		AddSeries("Category B", generateLineItems(u, 1)).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create(s)
	_ = line.Render(f)
}

func solveStraight(u [step][2]float64) [step][2]float64 {
	x := [step][2]float64{}
	x[0][0], x[0][1] = 0, 0
	for i := 1; i < step; i++ {
		t := float64(i) * tau
		x[i][0] = x[i-1][0] + tau*(math.Cos(t)*x[i-1][0]+t*x[i-1][1]+u[i][0]*(1-math.Exp(-t)))
		x[i][1] = x[i-1][1] + tau*(math.Sin(t)*x[i-1][1]+x[i-1][0]/(t+1)+u[i][1]*(1+math.Sin(2*t)))
	}
	return x
}

func solveReverse(x1, x2 float64) [step][2]float64 {
	psi := [step][2]float64{}
	psi[step-1][0], psi[step-1][1] = 2*(x1-y1), 2*(x2-y2)
	for i := step - 2; i >= 0; i-- {
		t := float64(i) * tau
		psi[i][0] = psi[i+1][0] - tau*(math.Cos(t)*psi[i+1][0]+psi[i+1][1]/(t+1))
		psi[i][1] = psi[i+1][1] - tau*(math.Sin(t)*psi[i+1][1]+t*psi[i+1][0])
	}
	return psi
}

func main() {
	var JDiff, psi, x, xJ, uJ, u [step][2]float64
	for i := 0; i < step; i++ {
		t := float64(i) * tau
		u[i][0] = 10.0 * t
		u[i][1] = 20.0 * t
	}
	for true {
		x = solveStraight(u)
		nrm := math.Sqrt(math.Pow(x[step-1][0]-y1, 2) + math.Pow(x[step-1][1]-y2, 2))
		if nrm < 0.001 {
			//fmt.Println(x[step-1][0], x[step-1][1])
			createLineChart(u, "U.html")
			createLineChart(x, "X.html")
			break
		}
		psi = solveReverse(x[step-1][0], x[step-1][1])
		for i := 0; i < step; i++ {
			t := float64(i) * tau
			JDiff[i][0] = psi[i][0] * (1 - math.Exp(-t))
			JDiff[i][1] = psi[i][1] * (1 + math.Sin(2*t))
			uJ[i][0] = u[i][0] - JDiff[i][0]
			uJ[i][1] = u[i][1] - JDiff[i][1]
		}
		xJ = solveStraight(uJ)
		int1 := math.Pow(JDiff[step-1][0], 2) + math.Pow(JDiff[step-1][1], 2) + math.Pow(JDiff[0][0], 2) + math.Pow(JDiff[0][1], 2)
		for i := 1; i < step-1; i++ {
			if i%2 == 1 {
				int1 += 4 * (math.Pow(JDiff[i][0], 2) + math.Pow(JDiff[i][1], 2))
			} else {
				int1 += 2 * (math.Pow(JDiff[i][0], 2) + math.Pow(JDiff[i][1], 2))
			}
		}
		int1 *= tau / 3
		norma := math.Pow(xJ[step-1][0]-x[step-1][0], 2) + math.Pow(xJ[step-1][1]-x[step-1][1], 2)
		alpha := 0.5 * int1 / norma
		for i := 0; i < step; i++ {
			u[i][0] -= alpha * JDiff[i][0]
			u[i][1] -= alpha * JDiff[i][1]
		}
	}
}
