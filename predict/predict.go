package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
)

func estimatePrice(milage float64, theta0 float64, theta1 float64) float64 {
	return theta1*milage + theta0
}

func main() {
	f, err := os.Open("model.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	model := dataframe.ReadJSON(f)

	fmt.Printf("Please input a mileage : ")
	var input float64
	fmt.Scanf("%f", &input)
	y := estimatePrice(input, model.Elem(0, 0).Float(), model.Elem(0, 1).Float())
	fmt.Printf("Prediction : %f\n", y)
}
