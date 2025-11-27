package big_int_generate

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"os"
	"time"
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)
	two  = big.NewInt(2)
)

func InitFile() error {
	nums := []string{"0", "0", "0", "0", "0", "0", "0", "0", "0", "0"}
	numsCh, err := GetDataFromFile(os.Getenv("FILE_WITH_BIG_NUMS"))
	if err == nil {
		nums = numsCh
	}
	return SaveDataToFile(os.Getenv("FILE_WITH_BIG_NUMS"), nums)
}

func GenerateRandomPrimeNums(size int) error {
	nums, err := GetDataFromFile(os.Getenv("FILE_WITH_BIG_NUMS"))
	if err != nil {
		log.Fatal(err)
		return err
	}
	startTime := time.Now()
	for i := 0; ; {
		num := GetRandomPrime(size)
		for !uniqueInNums(num.String(), nums) {
			num = GetRandomPrime(size)
		}
		log.Printf("Generated random prime number by %s", time.Now().Sub(startTime).String())
		startTime = time.Now()
		nums[i] = num.String()
		err = SaveDataToFile(os.Getenv("FILE_WITH_BIG_NUMS"), nums)
		if err != nil {
			log.Fatal(err)
		}
		i = (i + 1) % 10
	}
}

func uniqueInNums(num string, nums []string) bool {
	for _, n := range nums {
		if n == num {
			return false
		}
	}
	return true
}

func GetDataFromFile(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	var data []string
	err = decoder.Decode(&data)
	return data, err
}

func SaveDataToFile(filename string, nums []string) error {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("SaveDataToFile: %v", err)
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", " ")
	return encoder.Encode(nums)
}

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
