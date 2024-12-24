package drawing

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"go-wechat/config"
)

const (
	// 请求地址
	Addr = "https://visual.volcengineapi.com"
	Path = "/" // 路径，不包含 Query

	// 请求接口信息
	Service = "cv"
	Region  = "cn-north-1"
	Action  = "CVProcess"
	Version = "2022-08-31"
)

func hmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func getSignedKey(secretKey, date, region, service string) []byte {
	kDate := hmacSHA256([]byte(secretKey), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "request")

	return kSigning
}

func hashSHA256(data []byte) []byte {
	hash := sha256.New()
	if _, err := hash.Write(data); err != nil {
		log.Printf("input hash err:%s", err.Error())
	}

	return hash.Sum(nil)
}

func DoRequest(method string, queries url.Values, body []byte) (*http.Response, error) {
	// 1. 构建请求
	queries.Set("Action", Action)
	queries.Set("Version", Version)
	requestAddr := fmt.Sprintf("%s%s?%s", Addr, Path, queries.Encode())
	log.Printf("request addr: %s\n", requestAddr)

	request, err := http.NewRequest(method, requestAddr, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("bad request: %w", err)
	}

	// 2. 构建签名材料
	now := time.Now()
	date := now.UTC().Format("20060102T150405Z")
	authDate := date[:8]
	request.Header.Set("X-Date", date)

	payload := hex.EncodeToString(hashSHA256(body))
	request.Header.Set("X-Content-Sha256", payload)
	request.Header.Set("Content-Type", "application/json")

	queryString := strings.Replace(queries.Encode(), "+", "%20", -1)
	signedHeaders := []string{"host", "x-date", "x-content-sha256", "content-type"}
	var headerList []string
	for _, header := range signedHeaders {
		if header == "host" {
			headerList = append(headerList, header+":"+request.Host)
		} else {
			v := request.Header.Get(header)
			headerList = append(headerList, header+":"+strings.TrimSpace(v))
		}
	}
	headerString := strings.Join(headerList, "\n")

	canonicalString := strings.Join([]string{
		method,
		Path,
		queryString,
		headerString + "\n",
		strings.Join(signedHeaders, ";"),
		payload,
	}, "\n")

	hashedCanonicalString := hex.EncodeToString(hashSHA256([]byte(canonicalString)))

	credentialScope := authDate + "/" + Region + "/" + Service + "/request"
	signString := strings.Join([]string{
		"HMAC-SHA256",
		date,
		credentialScope,
		hashedCanonicalString,
	}, "\n")

	// 3. 构建认证请求头
	signedKey := getSignedKey(config.Conf.Ai.DrawApiSecret, authDate, Region, Service)
	signature := hex.EncodeToString(hmacSHA256(signedKey, signString))

	authorization := "HMAC-SHA256" +
		" Credential=" + config.Conf.Ai.DrawApiKey + "/" + credentialScope +
		", SignedHeaders=" + strings.Join(signedHeaders, ";") +
		", Signature=" + signature
	request.Header.Set("Authorization", authorization)

	// 4. 打印请求，发起请求
	requestRaw, err := httputil.DumpRequest(request, true)
	if err != nil {
		return nil, fmt.Errorf("dump request err: %w", err)
	}

	log.Printf("request:\n%s\n", string(requestRaw))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do request err: %w", err)
	}

	return response, nil
}
