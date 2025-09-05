package encodingutils

import (
	"errors"
	"math/big"
)

func GenerateEncodedBase62String(length *big.Int) string {
	if length.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}

	base62Chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	base := big.NewInt(62)
	result := ""
	zero := big.NewInt(0)
	rem := new(big.Int)

	for length.Cmp(zero) > 0 {
		length.DivMod(length, base, rem)
		result = string(base62Chars[rem.Int64()]) + result
	}

	return result
}

func DecodeEncodedBase62String(encoded string) (*big.Int, error) {
	base62Chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var base62index = func() map[rune]int {
		m := make(map[rune]int)
		for i, c := range base62Chars {
			m[c] = i
		}

		return m
	}()

	num := big.NewInt(0)
	base := big.NewInt(62)

	for _, c := range encoded {
		val, ok := base62index[c]
		if !ok {
			return nil, errors.New("invalid base62 character")
		}

		num.Mul(num, base)
		num.Add(num, big.NewInt(int64(val)))
	}

	return num, nil
}
