module github.com/samlotti/chess_anaylzer/chessboard/main

go 1.18

require (
	github.com/samlotti/blip v0.8.10
	github.com/samlotti/chess_anaylzer/analyzer v0.0.0-00010101000000-000000000000
	github.com/samlotti/chess_anaylzer/httpservice v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.3 // indirect

)

replace (
	github.com/samlotti/chess_anaylzer/analyzer => ../analyzer
	github.com/samlotti/chess_anaylzer/httpservice => ../httpservice
)
