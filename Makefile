all:
	mkdir -p bin
	go build -o bin/predict predict/predict.go
	go build -o bin/train train/train.go

clean:
	rm -f model.json
	rm -f ml.png
	rm -rf ./bin

re : clean all

.PHONY: predict train