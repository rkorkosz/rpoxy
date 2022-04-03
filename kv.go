package main

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type KV struct {
	client   *http.Client
	cache    map[string]URL
	endpoint string
}

func NewKV(endpoint string, client *http.Client) (*KV, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if client == nil {
		transport := httpTransport()
		client = &http.Client{Transport: transport}
	}
	return &KV{client: client, endpoint: u.String()}, nil
}

func (kv *KV) GetHost(ctx context.Context, host string) (URL, error) {
	u, ok := kv.cache[host]
	if ok {
		return u, nil
	}
	return kv.getHost(ctx, host)
}

func (kv *KV) getHost(ctx context.Context, host string) (URL, error) {
	u := kv.endpoint + "/hosts/" + host
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return URL{}, err
	}
	req.WithContext(ctx)
	resp, err := kv.client.Do(req)
	if err != nil {
		return URL{}, err
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return URL{}, err
	}
	parsed, err := url.Parse(string(out))
	if err != nil {
		return URL{}, err
	}
	return URL(*parsed), nil
}
