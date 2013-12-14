package interval

import (
	"math/rand"
)

const (
	NUM_APPS = 1000
)

//Interval represents a part of the real number line [Left,Right). Intervals
//typicall maintain other state (such as a name) about the interval in question.
type Interval interface {
	Left() float64
	Right() float64
}

//Distribution returns a collection of Intervals that cover the region [0.0, 1.0)
//in increasing order.
type Distribution interface {
	Len() int
	Ith(int) Interval
}

//ChooseRandomItem selects a random interval from the distribution provided.
//It uses a binary search to find the interval, once the "die is cast." This
//function panics if the interval does not completely cover [0,1). This
//function uses the supplied source of randomness and if that is nil uses
//rand.Float64()
func ChooseRandomItem(dist Distribution, r *rand.Rand) Interval {
	var coin float64
	if r == nil {
		coin = rand.Float64()
	} else {
		coin = r.Float64()
	}
	lower := 0
	upper := dist.Len()
	curr := (upper - lower) / 2
	for {
		interval := dist.Ith(curr)
		if interval.Left() <= coin && coin < interval.Right() {
			return interval
		}
		var dist int
		if coin < interval.Left() {
			dist = -((curr - lower) / 2)
			upper = curr
		} else {
			dist = (upper - curr) / 2
			lower = curr
		}
		//tricky: never want to move zero places, no matter how small the delta
		//so we make 0 be one for the purpose of moving... but we need to know
		//which *way* to move
		if dist == 0 {
			if upper == curr { //testing the above assignment
				dist = -1
			} else {
				dist = 1
			}
		}
		curr += dist
	}
	panic("distribution does not cover [0,1)")
}

//SimpleEqualDistribution is a distribution that divides the region [0, 1)
//evenly amount its keys.
type SimpleEqualDistribution struct {
	keys     []string
	interval [][2]float64
}

//SimpleInterval is an interval represented by two float64s. The left
//value is first and is inclusive, and the right value is second and exclusive.
//It maintains its own copy of it's interval position as wel.
type SimpleInterval struct {
	I  int
	LR [2]float64
}

//Left returns the first of the two values that make up this interval.
func (self SimpleInterval) Left() float64 {
	return self.LR[0]
}

//Right returns the second of the two values that make up this interval.
func (self SimpleInterval) Right() float64 {
	return self.LR[1]
}

//NewEqualDistribution creates a SimpleEqualDistribution based on the provided
//set of choices.
func NewEqualDistribution(keys []string) *SimpleEqualDistribution {
	ints := make([][2]float64, len(keys))
	k := make([]string, len(keys))
	prev := 0.0
	for i, key := range keys {
		k[i] = key
		if i != len(keys)-1 {
			right := float64(i+1) / float64(len(keys))
			ints[i] = [2]float64{prev, right}
			prev = right
		} else {
			ints[i] = [2]float64{prev, 1.0}
		}
	}
	return &SimpleEqualDistribution{k, ints}
}

//Len returns the number of keys that were used to construct this distribution.
func (self *SimpleEqualDistribution) Len() int {
	return len(self.keys)
}

//Ith returns the ith interval in this distribution.
func (self *SimpleEqualDistribution) Ith(i int) Interval {
	return SimpleInterval{i, self.interval[i]}
}
