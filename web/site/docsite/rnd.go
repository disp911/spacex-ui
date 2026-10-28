package docsite

import (
	"hash/fnv"
	"math/bits"
	"math/rand/v2"
)

// rnd is the generator's own source of choices. PCG gives the numbers, which
// its algorithm fixes for good; the helpers here are ours, so a Go release
// that changes the helpers of math/rand never changes a site already shown.
type rnd struct{ src *rand.PCG }

// newRnd starts a stream of its own for one part of the site, so that a
// change in how one part is made leaves the choices of the others as they
// were.
func newRnd(seed uint64, part string) *rnd {
	h := fnv.New64a()
	h.Write([]byte(part))
	// Seeds that differ in a few bits must still make unrelated sites.
	s := mix(seed ^ h.Sum64())
	return &rnd{rand.NewPCG(s, mix(s))}
}

// mix is the finaliser of SplitMix64.
func mix(z uint64) uint64 {
	z += 0x9e3779b97f4a7c15
	z = (z ^ z>>30) * 0xbf58476d1ce4e5b9
	z = (z ^ z>>27) * 0x94d049bb133111eb
	return z ^ z>>31
}

// intn takes the high bits of a number, which are the well-mixed ones.
func (r *rnd) intn(n int) int {
	hi, _ := bits.Mul64(r.src.Uint64(), uint64(n))
	return int(hi)
}

// between is a number from lo to hi, both included.
func (r *rnd) between(lo, hi int) int {
	return lo + r.intn(hi-lo+1)
}

func (r *rnd) chance(p float64) bool {
	return float64(r.src.Uint64()>>11)/(1<<53) < p
}

func pick[T any](r *rnd, s []T) T {
	return s[r.intn(len(s))]
}

func shuffled[T any](r *rnd, s []T) []T {
	out := append([]T(nil), s...)
	for i := len(out) - 1; i > 0; i-- {
		j := r.intn(i + 1)
		out[i], out[j] = out[j], out[i]
	}
	return out
}
