package helpers

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	uuid "github.com/satori/go.uuid"
)

var namespace = Must(FromString("0f9da948-a6fb-4c45-9edc-4685c3f3317d"))

func getUserId(address string) string {
	return uuid.NewV5(namespace, address).String()
}

func GetAccountId(address string) string {
	return uuid.NewV5(namespace, getUserId(strings.ToLower(address))+strconv.Itoa(0)).String()
}

func FromString(input string) (u uuid.UUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

func Must(u uuid.UUID, err error) uuid.UUID {
	if err != nil {
		panic(err)
	}
	return u
}

func RandomClientId() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%016d", rand.Intn(10000000000000000))
}

func GenerateQueryPath(url string, params url.Values) string {
	if len(params) == 0 {
		return url
	}
	return fmt.Sprintf("%s?%s", url, params.Encode())
}

func CreateTypedSignature(signature string, sigType int) string {
	return fmt.Sprintf("%s0%s", fixRawSignature(signature), strconv.Itoa(sigType))
}

func fixRawSignature(signature string) string {
	stripped := strings.TrimPrefix(signature, "0x")
	if len(stripped) != 130 {
		panic(fmt.Sprintf("Invalid raw signature: %s", signature))
	}
	rs := stripped[:128]
	v := stripped[128:130]
	if v == "00" {
		return "0x" + rs + "1b"
	}
	if v == "01" {
		return "0x" + rs + "1c"
	}
	if v == "1b" || v == "1c" {
		return "0x" + stripped
	}
	panic(fmt.Sprintf("Invalid v value: %s", v))
}

func HashString(input string) string {
	return hexutil.Encode(solsha3.SoliditySHA3([]string{"string"}, input))
}

func ExpireAfter(duration time.Duration) string {
	return time.Now().Add(duration).UTC().Format("2006-01-02T15:04:05.999Z")
}
