package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"booking.com/pkg/constants"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateTlsConfig(certFilePath, keyFilePath,
	caCertFilePath string) (*tls.Config, error) {
	tlsCfg := &tls.Config{}

	certificate, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		return nil, err
	}
	tlsCfg.Certificates = []tls.Certificate{certificate}

	if caCertFilePath != "" {
		certPool, err := CreateCertPool(caCertFilePath)
		if err != nil {
			return nil, err
		}
		tlsCfg.RootCAs = certPool
	}
	tlsCfg.MinVersion = tls.VersionTLS12
	return tlsCfg, nil
}

func CreateCertPool(caCertFilePath string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(caCertFilePath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	return certPool, nil
}
func ReadRequestBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(reqData))
	return reqData, nil
}
func ParseHttpRequest(r *http.Request, obj any) error {
	reqData, err := ReadRequestBody(r)
	if err != nil {
		return err
	}
	log.Print(string(reqData))
	return json.Unmarshal(reqData, obj)
}

func BuildHttpRequest(method, path string,
	queryParams map[string]string, payload any) (*http.Request, error) {
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	if queryParams != nil {
		query := url.Query()
		for key, val := range queryParams {
			query.Add(key, val)
		}
		url.RawQuery = query.Encode()
	}
	if payload != nil {
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		return http.NewRequest(method, url.String(), bytes.NewBuffer(bodyBytes))
	}
	return http.NewRequest(method, url.String(), nil)
}

func ReadHttpResponse(rsp *http.Response) ([]byte, error) {
	rspBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	rsp.Body = io.NopCloser(bytes.NewBuffer(rspBytes))
	return rspBytes, nil
}

func ParseHttpResponse(rsp *http.Response, response any) error {
	respBytes, err := ReadHttpResponse(rsp)
	if err != nil {
		return err
	}
	cType := rsp.Header.Get(constants.ContentType)
	if strings.Contains(cType, constants.ContentTypeTextPlain) {
		if strPtr, ok := response.(*string); ok {
			*strPtr = string(respBytes)
			return nil
		} else {
			return errors.New("invalid rsp type")
		}
	}
	return json.Unmarshal(respBytes, response)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func GetUUID() string {
	return uuid.New().String()
}

func GetExpTime(exp int64) int64 {
	return time.Now().Add(time.Duration(exp) * time.Minute).Unix()
}
func IsExpired(exp int64) bool {
	return exp < time.Now().Unix()
}
