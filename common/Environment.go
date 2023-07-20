package common

type environment struct {
	EnginePath string // the path to the engine location
}

var Environment = &environment{
	EnginePath: "../engines/mac/",
}
