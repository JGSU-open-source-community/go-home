package call

import (
	"testing"
)

func TestCall(t *testing.T) {
	datas := Call("G4775")

	t.Log(string(datas))
}
