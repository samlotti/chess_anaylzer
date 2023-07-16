module github.com/samlotti/chess_anaylzer/chessboard

go 1.18

replace github.com/samlotti/chess_anaylzer/minilex => ../minilex

require (
	github.com/samlotti/chess_anaylzer/minilex v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
