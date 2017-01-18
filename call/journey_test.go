package call

import (
	"testing"
)

func TestCall(t *testing.T) {
	datas := Call()

	t.Log(string(datas))
}
