package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKVGetHost(t *testing.T) {
	expected, err := url.Parse("http://localhost:8000")
	require.NoError(t, err)
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(expected.String()))
	}))
	ts.EnableHTTP2 = true
	ts.StartTLS()
	defer ts.Close()
	kv := NewKV(ts.URL, ts.Client())
	ctx := context.Background()
	actual, err := kv.GetHost(ctx, "test")
	require.NoError(t, err)
	assert.Equal(t, URL(*expected), actual)
}

var result URL

func BenchmarkKVGetHost(b *testing.B) {
	expected, err := url.Parse("http://localhost:8000")
	require.NoError(b, err)
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(expected.String()))
	}))
	ts.EnableHTTP2 = true
	ts.StartTLS()
	defer ts.Close()
	kv := NewKV(ts.URL, ts.Client())
	ctx := context.Background()
	var r URL
	for n := 0; n < b.N; n++ {
		r, _ = kv.GetHost(ctx, "test")
	}
	result = r
}
