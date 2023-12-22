BIN := flow

build:
	@echo "building..."
	cd src/run && go build -o $(BIN) .
