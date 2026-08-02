package dev

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestDownloadPackageReportsProgress(t *testing.T) {
	packageData := []byte("a plugin package large enough to arrive in a few writes")
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "/plugins/plugin-instance/package", request.URL.Path)
		response.Header().Set("Content-Type", "application/zip")
		response.Header().Set("Content-Length", fmt.Sprint(len(packageData)))
		response.WriteHeader(http.StatusOK)
		_, err := response.Write(packageData)
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := openapi.NewClientWithResponses(server.URL)
	require.NoError(t, err)

	type progressEvent struct {
		downloaded int64
		total      int64
	}
	events := []progressEvent{}

	data, err := DownloadPackage(
		context.Background(),
		client,
		"plugin-instance",
		func(downloaded, total int64) {
			events = append(events, progressEvent{downloaded: downloaded, total: total})
		},
	)

	require.NoError(t, err)
	assert.Equal(t, packageData, data)
	require.NotEmpty(t, events)
	assert.Equal(t, progressEvent{downloaded: 0, total: int64(len(packageData))}, events[0])
	assert.Equal(t, progressEvent{downloaded: int64(len(packageData)), total: int64(len(packageData))}, events[len(events)-1])
}
