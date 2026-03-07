#!make
.DEFAULT_GOAL := run

build:
	go build -o newsbot .

run:
	go run .

run_metadata:
	go run . -metadata

run_metadata_key:
	go run . -metadata-key=$(KEY)

run_single:
	go run . -single
