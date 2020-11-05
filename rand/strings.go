package rand

import (
	"crypto/rand"
	"encoding/base64"
)

//RememberTokenBytes is used to simplify our life by not having to remember the number
const RememberTokenBytes = 32

//Read works by taking a byte slice and filling it with random values then returns two values an integer and an error. The integer will be equal to the length of a byte slice if there weren't any errors but it could be smaller if there was an error.If that case the amount of bytes written dictates when the error occured.
func Read(b []byte) (n int, err error) {
	return
}

//Bytes will generate n random bytes, or will
//return an error if there was one.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//String will generate a byte slice of size nBytes and then return a string that is based on the base64 URL encoded version of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

//RememberToken is a helper function used to generate remember tokens of a predeteremined bytesize
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
