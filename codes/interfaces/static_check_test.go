package interfaces

import "testing"

type TestInterface interface {
	Name() string
}

type Instance struct{}

func (i *Instance) Name() string {
	return "impl TestInterface"
}

type WrongInstance struct{}

var _ TestInterface = (*Instance)(nil)

func TestStaticCheck(t *testing.T) {
	ins := Instance{}
	ins.Name()

	// wins := WrongInstance{}
	// wins.Name()
}
