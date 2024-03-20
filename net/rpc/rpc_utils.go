package rpc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	APIAddress  string
	LoginUser   string
	AccessToken string
	BasicToken  string
	CustomToken string
	UserAgent   string
	Timeout     int
}

func CallAPI(cfg *Config, path, method string, query url.Values, body []byte, apiRet APIRet) (err error) {
	return CallAPIWithContext(context.Background(), cfg, path, method, query, body, apiRet)
}

func CallAPIWithContext(ctx context.Context, cfg *Config, path, method string, query url.Values, body []byte, apiRet APIRet) (err error) {
	reqURL := fmt.Sprintf("%s%s", strings.TrimSuffix(cfg.APIAddress, "/"), path)
	rpcClient := NewClientWithTimeout(time.Duration(cfg.Timeout) * time.Second)
	header := http.Header{}
	// check access token
	if cfg.AccessToken != "" {
		header.Set("Authorization", "Bearer "+cfg.AccessToken)
	}
	// check basic token
	if cfg.BasicToken != "" {
		header.Set("Authorization", "Basic "+cfg.BasicToken)
	}
	// check custom token
	if cfg.CustomToken != "" {
		header.Set("Authorization", cfg.CustomToken)
	}
	// check login user
	if cfg.LoginUser != "" {
		header.Set("X-Login-User", cfg.LoginUser)
	}
	if cfg.UserAgent != "" {
		header.Set("User-Agent", cfg.UserAgent)
	}
	// ends
	err = rpcClient.Call(ctx, reqURL, method, header, query, body, apiRet)
	return
}
