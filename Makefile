GOSRC:=$(shell find ./ -iname *.go)
TARGET:=bin/plot_exporter

$(TARGET): $(GOSRC)
	go build -o $(TARGET) ./cmd/plot_exporter
