package interval

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestBasic(T *testing.T) {

	//maybe should seed here from dev random?
	source := rand.NewSource(time.Now().UTC().UnixNano())
	randomness := rand.New(source)

	RGB := []string{"red", "green", "blue"}

	distribution := NewEqualDistribution(RGB)

	if distribution.Len() != 3 {
		T.Fatalf("incorrect number of intervals, expect 3 but got %d", distribution.Len())
	}
	r := distribution.Ith(0)
	if r.Left() != 0.0 {
		T.Fatalf("red interval doesn't cover zero: %0.5f", r.Left())
	}
	if r.Right()-0.333 > 0.001 {
		T.Fatalf("red interval right should be 1/3 but is %0.5f", r.Right())
	}
	b := distribution.Ith(2)
	if r.Left()-0.666 > 0.001 {
		T.Fatalf("blue interval should start at 2/3 but is %0.5f", b.Left())
	}
	if b.Right() != 1.0 {
		T.Fatalf("blue interval right should be 1.0 but is %0.5f", b.Right())
	}

	count := make(map[string]int)

	ITERS := 1000000
	ONEPERCENT := 10000.0

	for i := 0; i < ITERS; i++ {
		interval := ChooseRandomItem(distribution, randomness)
		color := RGB[interval.(SimpleInterval).I]
		count[color] = count[color] + 1
	}

	for _, color := range RGB {
		//10K is 1% of 1M
		if math.Abs(float64(ITERS/len(RGB))-float64(count[color])) > ONEPERCENT {
			T.Fatalf("%s count is off by more than acceptable amount: %d", color, count[color])
		}
	}
}
