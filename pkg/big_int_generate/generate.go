package big_int_generate

import (
	"crypto/rand"
	"math/big"
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)
	two  = big.NewInt(2)
)

func PowModulo(b *big.Int, e *big.Int, m *big.Int) *big.Int {
	result := new(big.Int)
	return result.Exp(b, e, m)
}

func GetRandom(size int) *big.Int {
	n := big.NewInt(0)
	if size == 0 {
		return n
	}
	for n.BitLen() != size {
		b := make([]byte, size/8)
		rand.Read(b)
		n.SetBytes(b)
	}
	return n
}

func GetRandomPrime(size int) *big.Int {
	n := GetRandom(size)
	p := GetNextPrime(n)
	return p
}

func GetNextPrime(n *big.Int) *big.Int {
	prime := false
	if n.Bit(0) == 0 {
		n.Add(n, one)
	}
	for !prime {
		n.Add(n, two)
		prime = IsPrime(n)
	}
	return n
}

func IsPrime(n *big.Int) bool {
	if n.Cmp(zero) == 0 {
		return false
	}
	if n.Cmp(one) == 0 {
		return false
	}
	if n.Bit(0) == 0 {
		return false
	}
	radixList := []int64{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97} //, 2591, 7879, 9041, 10663, 11587, 12637};
	for _, rr := range radixList {
		nb := big.NewInt(0)
		nb.Add(zero, n)
		radix := big.NewInt(rr)
		if rr <= 7 {
			if nb.Mod(n, radix).Cmp(zero) == 0 {
				return false
			}
		}
		nb.Add(zero, n)
		if !isPrimeForRadix(nb, radix) {
			return false
		}
	}
	return true
}

func isPrimeForRadix(n *big.Int, radix *big.Int) bool {
	//Test Miller-Rabin
	if n.Cmp(radix) == 0 {
		return true
	}
	dd := big.NewInt(0)
	dd = dd.Sub(n, one)
	hh := 0
	for dd.Bit(0) == 0 {
		hh++
		dd.Div(dd, two)
	}
	nn1 := big.NewInt(0)
	nn1 = nn1.Sub(n, one)
	xx := big.NewInt(0)
	xx = PowModulo(radix, dd, n)
	if xx.Cmp(one) == 0 {
		return true
	}
	if xx.Cmp(nn1) == 0 {
		return true
	}
	ii := 1
	for ii != hh {
		xx = PowModulo(xx, two, n)
		if xx.Cmp(one) == 0 {
			return false
		}
		if xx.Cmp(nn1) == 0 {
			return true
		}
		ii++
	}
	return false
}
