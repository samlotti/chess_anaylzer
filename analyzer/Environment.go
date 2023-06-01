package analyzer

type environment struct {
	EnginePath string // the path to the engine location
}

var Environment = &environment{
	EnginePath: "../engines/",
}
