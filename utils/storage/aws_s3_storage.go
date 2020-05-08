package storage

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/duoke/base/hash"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	AwsRequestTimeout    = 30 // seconds
	AwsUnsignedPayload   = "UNSIGNED-PAYLOAD"
	AwsSigningAlgorithm  = "AWS4-HMAC-SHA256"
	AwsISO8601DateFormat = "20060102T150405Z"
	AwsScopeDateFormat   = "20060102"
	AwsS3RegionName      = "cn-north-1"
	AwsS3ServiceName     = "s3"
	AwsS3ServiceRequest  = "aws4_request"
)

type AwsS3Config struct {
	Bucket    string `yaml:"bucket"`
	Domain    string `yaml:"domain"`
	KeyPrefix string `yaml:"keyPrefix"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Timeout   int    `yaml:"timeout"`
}

// AwsS3SignRequest implements the aws signing algorithm version 4.
// This implementation if part of the algorithm only cater to our needs here.
func AwsS3SignRequest(storageCfg *AwsS3Config, req *http.Request) {
	// canonical headers
	reqDate := time.Now().UTC()
	reqDateStr := reqDate.Format(AwsISO8601DateFormat)
	canonicalHeaders := map[string]string{
		"x-amz-content-sha256": AwsUnsignedPayload,
		"x-amz-date":           reqDateStr,
	}

	// add to the req
	for key, value := range canonicalHeaders {
		req.Header.Add(key, value)
	}

	if contentType := req.Header.Get("Content-Type"); contentType != "" {
		canonicalHeaders["content-type"] = contentType
	}

	canonicalHeaders["host"] = req.Host

	// get headers keys
	canonicalHeadersKeys := make([]string, 0, len(canonicalHeaders))
	for key, _ := range canonicalHeaders {
		canonicalHeadersKeys = append(canonicalHeadersKeys, key)
	}

	// signed headers keys
	sort.Strings(canonicalHeadersKeys)
	signedHeadersStr := strings.Join(canonicalHeadersKeys, ";")

	// pack headers whose keys are sorted
	canonicalHeadersItems := make([]string, 0, len(canonicalHeaders))
	for _, key := range canonicalHeadersKeys {
		canonicalHeadersItems = append(canonicalHeadersItems, fmt.Sprintf("%s:%s", key, canonicalHeaders[key]))
	}

	// pack request items
	canonicalRequestItems := []string{
		req.Method,
		req.URL.Path,
		strings.Replace(req.URL.Query().Encode(), "+", "%20", -1),
		strings.Join(canonicalHeadersItems, "\n") + "\n",
		signedHeadersStr,
		AwsUnsignedPayload,
	}

	// prepare the str to sign
	canonicalRequestData := []byte(strings.Join(canonicalRequestItems, "\n"))
	credentialScope := strings.Join([]string{reqDate.Format(AwsScopeDateFormat), AwsS3RegionName, AwsS3ServiceName, AwsS3ServiceRequest}, "/")
	strToSign := strings.Join([]string{AwsSigningAlgorithm, reqDateStr, credentialScope, hash.Sha256HexString(canonicalRequestData)}, "\n")

	// calc the signing key
	dateKey := hash.HmacSha256([]byte(reqDate.Format(AwsScopeDateFormat)), []byte(fmt.Sprintf("AWS4%s", storageCfg.SecretKey)))
	dateRegionKey := hash.HmacSha256([]byte(AwsS3RegionName), dateKey)
	dateRegionServiceKey := hash.HmacSha256([]byte(AwsS3ServiceName), dateRegionKey)
	signingKey := hash.HmacSha256([]byte(AwsS3ServiceRequest), dateRegionServiceKey)

	awsSignature := hex.EncodeToString(hash.HmacSha256([]byte(strToSign), signingKey))
	awsCredential := fmt.Sprintf("%s/%s", storageCfg.AccessKey, credentialScope)

	authToken := fmt.Sprintf("%s Credential=%s,SignedHeaders=%s,Signature=%s", AwsSigningAlgorithm, awsCredential, signedHeadersStr, awsSignature)
	req.Header.Add("Authorization", authToken)
	req.Header.Add("Date", reqDate.Format(http.TimeFormat))
}

func AwsS3PutObject(storageCfg *AwsS3Config, fileKey string, fileData []byte) (err error) {
	endPoint := fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(storageCfg.Domain, "/"), storageCfg.Bucket, fileKey)
	req, newErr := http.NewRequest(http.MethodPut, endPoint, bytes.NewReader(fileData))
	if newErr != nil {
		err = newErr
		return
	}

	// set content-type
	req.Header.Add("Content-Type", "text/plain")

	// sign the request
	AwsS3SignRequest(storageCfg, req)

	// new http client
	client := http.DefaultClient
	if storageCfg.Timeout > 0 {
		client.Timeout = time.Second * time.Duration(storageCfg.Timeout)
	} else {
		client.Timeout = time.Second * time.Duration(AwsRequestTimeout)
	}

	// fire request
	resp, respErr := client.Do(req)
	if respErr != nil {
		err = respErr
		return
	}

	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusOK {
		// create error from the body
		errData, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			err = readErr
			return
		}
		err = errors.New(string(errData))
		return
	}

	return
}

func AwsS3GetObject(storageCfg *AwsS3Config, fileKey string) (fileBody []byte, err error) {
	endPoint := fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(storageCfg.Domain, "/"), storageCfg.Bucket, fileKey)
	req, newErr := http.NewRequest(http.MethodGet, endPoint, nil)
	if newErr != nil {
		err = newErr
		return
	}

	// sign the request
	AwsS3SignRequest(storageCfg, req)

	// new http client
	client := http.DefaultClient
	if storageCfg.Timeout > 0 {
		client.Timeout = time.Second * time.Duration(storageCfg.Timeout)
	} else {
		client.Timeout = time.Second * time.Duration(AwsRequestTimeout)
	}

	// fire request
	resp, respErr := client.Do(req)
	if respErr != nil {
		err = respErr
		return
	}

	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			err = errors.New(fmt.Sprintf("%s", resp.Status))
		} else {
			// create error from the body
			errData, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				err = readErr
				return
			}
			err = errors.New(string(errData))
		}
		return
	}

	// read file body
	fileBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}
