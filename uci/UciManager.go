package uci

import "fmt"

type _UciManager struct {
}

func (m *_UciManager) GetUci(engine string) (*UciProcess, error) {
	if Verbose {
		fmt.Printf("Get UCI: %s\n", engine)
	}
	u := NewUci(engine)
	err := u.Start()
	if err != nil {
		return nil, err
	}
	err = u.SendUciNewGame()
	return u, err
}

func (m *_UciManager) Return(uci *UciProcess) {
	uci.Terminate()
	if Verbose {
		fmt.Printf("End UCI\n")
	}
}

var _manager *_UciManager = &_UciManager{}

func UciManager() *_UciManager {
	return _manager
}
