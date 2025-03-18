package dogs

import (
	"os"
	"testing"

	"github.com/rhargreaves/dog-walking/test/e2e/common"
)

func TestMain(m *testing.M) {
	common.Authenticate(&testing.T{})
	code := m.Run()
	os.Exit(code)
}
