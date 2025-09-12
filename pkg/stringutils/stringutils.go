package stringutils

import (
	crypto_rand "crypto/rand"
	"encoding/base64"
	math_rand "math/rand"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type RandomStringOptions string

const (
	RandomStringOptionsAlphanumeric RandomStringOptions = "alphanumeric"
	RandomStringOptionsAlphabetic   RandomStringOptions = "alphabetic"
	RandomStringOptionsNumeric      RandomStringOptions = "numeric"
)

func GenerateRandomString(length int, options RandomStringOptions) string {
	math_rand.New(math_rand.NewSource(time.Now().UnixNano()))

	var chars []rune

	switch options {
	case RandomStringOptionsAlphanumeric:
		chars = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case RandomStringOptionsAlphabetic:
		chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case RandomStringOptionsNumeric:
		chars = []rune("0123456789")
	default:
		chars = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}

	b := make([]rune, length)

	for i := range b {
		b[i] = chars[math_rand.Intn(len(chars))]
	}

	return string(b)
}

func GenerateUniqueString(length int) string {
	b := make([]byte, length)
	_, err := crypto_rand.Read(b)
	if err != nil {
		panic(err)
	}

	s := base64.RawURLEncoding.EncodeToString(b)

	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, s)[:length]
}

func CamelCaseToPascalCase(s string) string {
	return cases.Title(language.English).String(s)
}
