package graphapi_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/metadata-api/internal/ent/generated/metadata"
	"go.infratographer.com/metadata-api/internal/ent/generated/status"
	"go.infratographer.com/metadata-api/internal/testclient"
)

func TestStatusUpdate(t *testing.T) {
	ctx := context.Background()
	meta1 := MetadataBuilder{}.MustNew(ctx)
	st1 := StatusBuilder{Metadata: meta1}.MustNew(ctx)

	testCases := []struct {
		TestName    string
		NodeID      gidx.PrefixedID
		NamespaceID gidx.PrefixedID
		Source      string
		ErrorMsg    string
	}{
		{
			TestName:    "Will create status for a node we don't have metadata for",
			NodeID:      gidx.MustNewID("testing"),
			NamespaceID: StatusNamespaceBuilder{}.MustNew(ctx).ID,
			Source:      "go-tests",
		},
		{
			TestName:    "Will create status for a node that has other metadata",
			NodeID:      meta1.NodeID,
			NamespaceID: StatusNamespaceBuilder{}.MustNew(ctx).ID,
			Source:      "go-tests",
		},
		{
			TestName:    "Will create status when status already exists from a different source",
			NodeID:      meta1.NodeID,
			NamespaceID: st1.StatusNamespaceID,
			Source:      "go-tests",
		},
		{
			TestName:    "Will update status when status already exists from the same source",
			NodeID:      meta1.NodeID,
			NamespaceID: st1.StatusNamespaceID,
			Source:      st1.Source,
		},
		{
			TestName:    "Fails when namespace doesn't exist",
			NodeID:      gidx.MustNewID("testing"),
			NamespaceID: gidx.MustNewID("notreal"),
			Source:      "go-tests",
			ErrorMsg:    "status_namespace not found",
		},
		{
			TestName:    "Fails when nodeID is empty",
			NodeID:      "",
			NamespaceID: st1.StatusNamespaceID,
			Source:      "go-tests",
			ErrorMsg:    "value is less than the required length",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			jsonData, err := gofakeit.JSON(nil)
			require.NoError(t, err)

			resp, err := graphTestClient().StatusUpdate(ctx, testclient.StatusUpdateInput{NodeID: tt.NodeID, NamespaceID: tt.NamespaceID, Source: tt.Source, Data: json.RawMessage(jsonData)})

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp.StatusUpdate.Status)
			assert.JSONEq(t, string(jsonData), string(resp.StatusUpdate.Status.Data))

			stCount := EntClient.Status.Query().Where(status.StatusNamespaceID(tt.NamespaceID), status.Source(tt.Source), status.HasMetadataWith(metadata.NodeID(tt.NodeID))).CountX(ctx)
			assert.Equal(t, 1, stCount)
		})
	}
}

func TestStatusDelete(t *testing.T) {
	ctx := context.Background()
	meta1 := MetadataBuilder{}.MustNew(ctx)
	st1 := StatusBuilder{Metadata: meta1}.MustNew(ctx)
	st2 := StatusBuilder{Metadata: meta1}.MustNew(ctx)

	testCases := []struct {
		TestName    string
		NodeID      gidx.PrefixedID
		NamespaceID gidx.PrefixedID
		Source      string
		ErrorMsg    string
	}{
		{
			TestName:    "Will delete status when found",
			NodeID:      meta1.NodeID,
			NamespaceID: st1.StatusNamespaceID,
			Source:      st1.Source,
		},
		{
			TestName:    "Will return an error if the status doesn't exists for the given source and namespace",
			NodeID:      meta1.NodeID,
			NamespaceID: st2.StatusNamespaceID,
			Source:      "this-is-not-source-you-are-looking-for",
			ErrorMsg:    "status not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().StatusDelete(ctx, testclient.StatusDeleteInput{NodeID: tt.NodeID, NamespaceID: tt.NamespaceID, Source: tt.Source})

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp.StatusDelete)
			assert.NotNil(t, resp.StatusDelete.DeletedID)

			count := EntClient.Status.Query().Where(status.Source(tt.Source), status.HasMetadataWith(metadata.NodeID(tt.NodeID))).CountX(ctx)
			assert.Equal(t, 0, count)
		})
	}
}
