package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
)

const (
  apiURL = "https://api.zadarma.com/v1/info/balance?format=json"
  userKey = ""
  secretKey = ""
)

// Zadarma-specific algorithm to get Authorization header
func getAuthorizationHeader() string {
  md5Hash := md5.Sum([]byte("format=json"))
  signature := fmt.Sprintf("/v1/info/balanceformat=json%s", hex.EncodeToString(md5Hash[:]))
  encodedSignature := encodeSignature(signature, secretKey)
	return fmt.Sprintf("%s:%s", userKey, encodedSignature)
}

func encodeSignature(signatureString, secret string) string {
	hmacSha1 := hmac.New(sha1.New, []byte(secret))
	hmacSha1.Write([]byte(signatureString))
	hmacHash := hmacSha1.Sum(nil)
	return base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(hmacHash)))
}

func main() {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

  authorizationHeader := getAuthorizationHeader()
	req.Header.Set("Authorization", authorizationHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}
