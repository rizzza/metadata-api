package graphapi_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/metadata-api/internal/ent/generated/annotation"
	"go.infratographer.com/metadata-api/internal/ent/generated/metadata"
	"go.infratographer.com/metadata-api/internal/testclient"
)

func TestAnnotationUpdate(t *testing.T) {
	ctx := context.Background()
	meta1 := MetadataBuilder{}.MustNew(ctx)
	ant1 := AnnotationBuilder{Metadata: meta1}.MustNew(ctx)

	testCases := []struct {
		TestName    string
		NodeID      gidx.PrefixedID
		NamespaceID gidx.PrefixedID
		ErrorMsg    string
	}{
		{
			TestName:    "Will create annotation for a node we don't have metadata for",
			NodeID:      gidx.MustNewID("testing"),
			NamespaceID: AnnotationNamespaceBuilder{}.MustNew(ctx).ID,
		},
		{
			TestName:    "Will create annotation for a node that has other metadata",
			NodeID:      meta1.NodeID,
			NamespaceID: AnnotationNamespaceBuilder{}.MustNew(ctx).ID,
		},
		{
			TestName:    "Will update annotation when annotation already exists",
			NodeID:      meta1.NodeID,
			NamespaceID: ant1.AnnotationNamespaceID,
		},
		{
			TestName:    "Fails when namespace doesn't exists",
			NodeID:      gidx.MustNewID("testing"),
			NamespaceID: gidx.MustNewID("notreal"),
			ErrorMsg:    "constraint failed", // TODO: This should have a better error message
		},
		{
			TestName:    "Fails when namespace is empty",
			NodeID:      gidx.MustNewID("testing"),
			NamespaceID: "",
			ErrorMsg:    "value is less than the required length",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			jsonData, err := gofakeit.JSON(nil)
			require.NoError(t, err)

			resp, err := graphTestClient().AnnotationUpdate(ctx, testclient.AnnotationUpdateInput{NodeID: tt.NodeID, NamespaceID: tt.NamespaceID, Data: json.RawMessage(jsonData)})

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp.AnnotationUpdate.Annotation)
			assert.JSONEq(t, string(jsonData), string(resp.AnnotationUpdate.Annotation.Data))

			antCount := EntClient.Annotation.Query().Where(annotation.AnnotationNamespaceID(tt.NamespaceID), annotation.HasMetadataWith(metadata.NodeID(tt.NodeID))).CountX(ctx)
			assert.Equal(t, 1, antCount)
		})
	}
}

func TestAnnotationDelete(t *testing.T) {
	ctx := context.Background()
	meta1 := MetadataBuilder{}.MustNew(ctx)

	testCases := []struct {
		TestName    string
		NodeID      gidx.PrefixedID
		NamespaceID gidx.PrefixedID
		ErrorMsg    string
	}{
		{
			TestName:    "Will delete annotation when found",
			NodeID:      meta1.NodeID,
			NamespaceID: AnnotationBuilder{Metadata: meta1}.MustNew(ctx).AnnotationNamespaceID,
		},
		{
			TestName:    "Will return an error if the annotation doesn't exists",
			NodeID:      meta1.NodeID,
			NamespaceID: AnnotationNamespaceBuilder{}.MustNew(ctx).ID,
			ErrorMsg:    "annotation not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().AnnotationDelete(ctx, testclient.AnnotationDeleteInput{NodeID: tt.NodeID, NamespaceID: tt.NamespaceID})

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp.AnnotationDelete)
			assert.NotNil(t, resp.AnnotationDelete.DeletedID)

			antCount := EntClient.Annotation.Query().Where(annotation.AnnotationNamespaceID(tt.NamespaceID), annotation.HasMetadataWith(metadata.NodeID(tt.NodeID))).CountX(ctx)
			assert.Equal(t, 0, antCount)
		})
	}
}
