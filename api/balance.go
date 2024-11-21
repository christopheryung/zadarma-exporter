package api

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
  apiURL = "https://api.zadarma.com/v1/info/balance?format=json"
  userKey = ""
  secretKey = ""
)
type Response struct {
  Status string `json:"status"`
  Balance float64 `json:"balance"`
  Currency string `json:"currency"`
}

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

func GetBalance() (float64, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
    log.Println(err)
		return 0, err
	}

  authorizationHeader := getAuthorizationHeader()
	req.Header.Set("Authorization", authorizationHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
    log.Println(err)
		return 0, err 
	}
	defer resp.Body.Close()

  bodyBytes, _ := io.ReadAll(resp.Body)
  var response Response
  err = json.Unmarshal(bodyBytes, &response)
  if err != nil {
    log.Println(err)
  }

  return response.Balance, err
}
