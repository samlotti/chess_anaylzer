package uci

type _UciManager struct {
}

func (m *_UciManager) GetUci(engine string) (*UciProcess, error) {
	u := NewUci(engine)
	err := u.Start()
	return u, err
}

func (m *_UciManager) Return(uci *UciProcess) {
	uci.Terminate()
}

var _manager *_UciManager = &_UciManager{}

func UciManager() *_UciManager {
	return _manager
}
