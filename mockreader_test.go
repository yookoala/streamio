package streamio_test

import (
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/yookoala/streamio"
)

func init() {
	// seed random
	rand.Seed(int64(time.Now().Nanosecond()))
}

func TestMockReader(t *testing.T) {
	limit := uint32(rand.Int31n(2000))
	r := streamio.NewMockReader('a', 30, limit)
	content, _ := ioutil.ReadAll(r)

	if have, want := uint32(len(content)), limit; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
