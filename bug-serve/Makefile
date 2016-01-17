all: Generate React Go

Generate:
	go generate
Go:
	go build .

js: React

jsx: React

React:
	babel --out-dir=js jsx

serve:
	./bug-serve
