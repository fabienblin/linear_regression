package main

import (
	"log"
	"os"

	"github.com/go-gota/gota/series"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"github.com/go-gota/gota/dataframe"
)

var learningRate float64 = .01
var iterations int = 150000
var maxKM float64 = 0
var maxPrice float64 = 0
var minKM float64 = 0
var minPrice float64 = 0
var ratioKM float64 = 0
var ratioPrice float64 = 0

func check(e error) {
	if e != nil {
		log.Fatalln("Error")
		os.Exit(1)
	}
}

// Adjusted function (ax + b) of linear regression
func estimatePrice(milage float64, theta0 float64, theta1 float64) float64 {
	return theta1*milage + theta0
}

// Adjust b
func tmpTheta0(data dataframe.DataFrame, theta0 float64, theta1 float64) float64 {
	nrow := data.Nrow()
	var cost float64 = 0

	for i := 0; i < nrow; i++ {
		km := data.Elem(i, 0).Float()
		price := data.Elem(i, 1).Float()
		cost += estimatePrice(km, theta0, theta1) - price
	}
	var tmp = cost * learningRate / float64(nrow)
	return theta0 - learningRate*tmp
}

// Adjust a
func tmpTheta1(data dataframe.DataFrame, theta0 float64, theta1 float64) float64 {
	nrow := data.Nrow()
	var cost float64 = 0

	for i := 0; i < nrow; i++ {
		km := data.Elem(i, 0).Float()
		price := data.Elem(i, 1).Float()
		cost += (estimatePrice(km-price, theta0, theta1) - price) * km
	}
	var tmp = cost * learningRate / float64(nrow)
	return theta1 - learningRate*tmp
}

// Build plotter and add data points
func plotData(data dataframe.DataFrame) *plot.Plot {
	plot, _ := plot.New()
	plot.Title.Text = "Linear Regression"
	plot.X.Label.Text = "km"
	plot.Y.Label.Text = "price"
	pts := make(plotter.XYs, data.Nrow())
	nrow := data.Nrow()

	for i := 0; i < nrow; i++ {
		km := data.Elem(i, 0).Float()
		price := data.Elem(i, 1).Float()
		pts[i].X = km
		pts[i].Y = price
	}
	s, _ := plotter.NewScatter(pts)
	plot.Add(plotter.NewGrid())
	plot.Add(s)

	return plot
}

// Add line estimatePrice() to plotter
func plotLine(plot *plot.Plot, data dataframe.DataFrame, theta0 float64, theta1 float64) {
	xy := make(plotter.XYs, 2)

	xy[0].X = 0
	xy[0].Y = estimatePrice(0, theta0, theta1)

	xy[1].X = data.Col("km").Max()
	xy[1].Y = estimatePrice(data.Col("km").Max(), theta0, theta1)

	line, _ := plotter.NewLine(xy)
	plot.Add(line)
}

// Data matrix modificator
func subtractMin(s series.Series) series.Series {
	if s.Name == "km" {
		for i := 0; i < s.Len(); i++ {
			s.Set(i, series.Floats(s.Elem(i).Float()-minKM))
		}
	} else if s.Name == "price" {
		for i := 0; i < s.Len(); i++ {
			s.Set(i, series.Floats(s.Elem(i).Float()-minPrice))
		}
	}
	return series.Floats(s)
}

// Data matrix modificator
func divideRatio(s series.Series) series.Series {
	if s.Name == "km" {
		for i := 0; i < s.Len(); i++ {
			s.Set(i, series.Floats(s.Elem(i).Float()/ratioKM))
		}
	} else if s.Name == "price" {
		for i := 0; i < s.Len(); i++ {
			s.Set(i, series.Floats(s.Elem(i).Float()/ratioPrice))
		}
	}
	return series.Floats(s)
}

// Data matrix modificator
func multiplyRatio(s series.Series) series.Series {
	if s.Name == "km" {
		for i := 0; i < s.Len(); i++ {
			s.Set(i, series.Floats(s.Elem(i).Float()*ratioKM))
		}
	} else if s.Name == "price" {
		for i := 0; i < s.Len(); i++ {
			s.Set(i, series.Floats(s.Elem(i).Float()*ratioPrice))
		}
	}
	return series.Floats(s)
}

// Data matrix modificator
func addMin(s series.Series) series.Series {
	if s.Name == "km" {
		for i := 0; i < s.Len(); i++ {
			s.Set(i, series.Floats(s.Elem(i).Float()+minKM))
		}
	} else if s.Name == "price" {
		for i := 0; i < s.Len(); i++ {
			s.Set(i, series.Floats(s.Elem(i).Float()+minPrice))
		}
	}
	return series.Floats(s)
}

// Normalize dataset to be defined between 0 and 1
func normalizeData(data dataframe.DataFrame) dataframe.DataFrame {

	// define minimums
	minKM = data.Col("km").Min()
	minPrice = data.Col("price").Min()

	// define maximums
	maxKM = data.Col("km").Max()
	maxPrice = data.Col("price").Max()

	// define ratios
	ratioKM = maxKM - minKM
	ratioPrice = maxPrice - minPrice

	// check ratios are not 0 for later division
	if ratioKM == 0 || ratioPrice == 0 {
		check(nil)
	}

	// substract min from data
	data = data.Capply(subtractMin)

	// divide data by max
	data = data.Capply(divideRatio)

	return data
}

// Denormalize dataset to it's original state
func denormalizeData(data dataframe.DataFrame) dataframe.DataFrame {
	data = data.Capply(multiplyRatio)
	data = data.Capply(addMin)
	return data
}

func main() {
	// Load data from CSV source
	dataFile, errData := os.Open("data.csv")
	check(errData)
	defer dataFile.Close()

	// Prepare model output on json file
	modelFile, errModel := os.Create("model.json")
	check(errModel)
	defer modelFile.Close()

	// Normalize data
	data := dataframe.ReadCSV(dataFile)
	if data.Err != nil {
		os.Exit(1)
	}
	data = normalizeData(data)
	if data.Err != nil {
		os.Exit(1)
	}

	// Train model
	var tmptheta0 float64 = 0
	var tmptheta1 float64 = 0
	var theta0 float64 = 0
	var theta1 float64 = 0

	for i := 0; i < iterations; i++ {
		tmptheta0 = tmpTheta0(data, theta0, theta1)
		tmptheta1 = tmpTheta1(data, theta0, theta1)
		theta0 = tmptheta0
		theta1 = tmptheta1
	}

	// Denormalize line and data
	theta1 = theta1 * (ratioPrice / ratioKM)
	theta0 = theta0*ratioPrice + minPrice - theta1*minKM
	data = denormalizeData(data)

	// Output trained model on json file and png image
	model := dataframe.LoadMaps([]map[string]interface{}{
		map[string]interface{}{"theta0": theta0, "theta1": theta1},
	})

	plot := plotData(data)
	plotLine(plot, data, model.Elem(0, 0).Float(), model.Elem(0, 1).Float())
	plot.Save(4*vg.Inch, 4*vg.Inch, "ml.png")

	model.WriteJSON(modelFile)
}
