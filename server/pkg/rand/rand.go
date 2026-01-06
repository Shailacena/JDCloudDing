package rand

import (
	"math/rand"
	"time"
)

var (
	Random *rand.Rand
)

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	Random = rand.New(source)
}
