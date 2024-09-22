package interfaces

type Backend interface {
	Number() int
}

type FilterAPI struct {
	sys Backend
}

func NewFilterAPI(sys Backend) *FilterAPI {
	return &FilterAPI{sys: sys}
}

func (f *FilterAPI) GetNumber() int {
	return f.sys.Number()
}
