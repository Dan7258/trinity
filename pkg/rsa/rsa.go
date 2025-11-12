package rsa

import (
	"bufio"
	"errors"
	"math/big"
	"os"
)

type Encrypt struct {
	e *big.Int
	n *big.Int
}

type Decrypt struct {
	d *big.Int
	n *big.Int
}

type LargeNums struct {
	p *big.Int
	q *big.Int
}


type EncryptedData struct {
	EncryptedMessage string `json:"encrypted_message"`
	D                string `json:"d"`
	N                string `json:"n"`
}

func EncodeData(text string) (*EncryptedData, error) {
	ed := new(EncryptedData)
	largeNums := &LargeNums{
		new(big.Int),
		new(big.Int),
	}
	encrypt := &Encrypt{
		new(big.Int),
		new(big.Int),
	}
	encrypt.e = big.NewInt(65537)

	err := largeNums.generateLargeNums()
	if err != nil {
		return nil, err
	}

	encrypt.n = new(big.Int).Mul(largeNums.p, largeNums.q)
	ed.N = encrypt.n.String()

	phi := euler(largeNums)
	d, err := generateD(encrypt, phi)
	if err != nil {
		return nil, err
	}
	ed.D = d.String()
	c, err := encode(encrypt, text)
	if err != nil {
		panic(err)
	}
	ed.EncryptedMessage = c.String()
	return ed, nil
}

func DecodeData(ed *EncryptedData) (*string, error) {
	decrypt := new(Decrypt)
	decrypt.d = new(big.Int)
	decrypt.n = new(big.Int)
	var ok bool
	_, ok = decrypt.d.SetString(ed.D, 10)
	if !ok {
		return nil, errors.New("failed to decode D")
	}

	_, ok = decrypt.n.SetString(ed.N, 10)
	if !ok {
		return nil, errors.New("failed to decode N")
	}
	em := new(big.Int)
	_, ok = em.SetString(ed.EncryptedMessage, 10)
	if !ok {
		return nil, errors.New("failed to decode Message")
	}
	message, err := decode(decrypt, em)
	return &message, err
}

func (ln *LargeNums) generateLargeNums() error {
	//cmd := exec.Command("make")
	//err := cmd.Run()
	//if err != nil {
	//	return err
	//}
	file, err := os.Open("./pkg/rsa/numbers.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	ln.p.SetString(scanner.Text(), 10)
	scanner.Scan()
	ln.q.SetString(scanner.Text(), 10)
	return nil
}

func euler(ln *LargeNums) *big.Int {
	one := big.NewInt(1)
	pmo := new(big.Int).Sub(ln.p, one)
	qmo := new(big.Int).Sub(ln.q, one)
	phi := new(big.Int).Mul(pmo, qmo)
	return phi
}

func generateD(encrypt *Encrypt, phi *big.Int) (*big.Int, error) {
	d := new(big.Int).ModInverse(encrypt.e, phi)
	if d == nil {
		return nil, errors.New("e и phi не взаимно просты, нужно менять p и q")
	}
	return d, nil
}

func encode(en *Encrypt, text string) (*big.Int, error) {
	m := new(big.Int).SetBytes([]byte(text))
	if m.Cmp(en.n) >= 0 {
		return nil, errors.New("слишком большое сообщение")
	}
	c := new(big.Int).Exp(m, en.e, en.n)
	return c, nil
}

func decode(decrypt *Decrypt, c *big.Int) (string, error) {

	dec := new(big.Int).Exp(c, decrypt.d, decrypt.n)
	return string(dec.Bytes()), nil
}
