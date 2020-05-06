all:	dependencies
	mkdir -p bin
	go build -o bin/predict predict/predict.go
	go build -o bin/train train/train.go

clean:
	rm -f model.json
	rm -f ml.png
	rm -rf ./bin

re : clean all

dependencies:
	go get github.com/go-gota/gota/dataframes
	go get github.com/go-gota/gota/series
	go get gonum.org/v1/plot

.PHONY: predict train
