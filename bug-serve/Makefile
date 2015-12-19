all: Go React serve

Go:
	go build .

js: React

jsx: React

React:
	babel --out-dir=js jsx

serve:
	./bug-serve
