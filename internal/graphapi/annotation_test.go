package graphapi_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/permissions-api/pkg/permissions/mockpermissions"
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/metadata-api/internal/ent/generated/annotation"
	"go.infratographer.com/metadata-api/internal/ent/generated/metadata"
	"go.infratographer.com/metadata-api/internal/testclient"
)

func TestAnnotationUpdate(t *testing.T) {
	ctx := context.Background()

	perms := new(mockpermissions.MockPermissions)
	ctx = perms.ContextWithHandler(ctx)

	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	meta1 := MetadataBuilder{}.MustNew(ctx)
	ant1 := AnnotationBuilder{Metadata: meta1}.MustNew(ctx)

	testCases := []struct {
		TestName    string
		NodeID      gidx.PrefixedID
		NamespaceID gidx.PrefixedID
		JSONData    json.RawMessage // optional, otherwise generated
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
			TestName:    "Fails when namespace doesn't exist",
			NodeID:      gidx.MustNewID("testing"),
			NamespaceID: gidx.MustNewID("notreal"),
			ErrorMsg:    "not found",
		},
		{
			TestName:    "Fails when namespace is empty",
			NodeID:      gidx.MustNewID("testing"),
			NamespaceID: "",
			ErrorMsg:    "must not be empty",
		},
		{
			TestName:    "Fails when namespace gidx is invalid",
			NodeID:      gidx.MustNewID("testing"),
			NamespaceID: "test-invalid-id",
			ErrorMsg:    "invalid id",
		},
		{
			TestName:    "Fails when node is empty",
			NodeID:      "",
			NamespaceID: gidx.MustNewID("testing"),
			ErrorMsg:    "must not be empty",
		},
		{
			TestName:    "Fails when node gidx is invalid",
			NodeID:      "test-invalid-id",
			NamespaceID: gidx.MustNewID("testing"),
			ErrorMsg:    "invalid id",
		},
		{
			TestName:    "Fails to update nodeID status with invalid json data",
			NodeID:      gidx.MustNewID("testing"),
			NamespaceID: AnnotationNamespaceBuilder{}.MustNew(ctx).ID,
			JSONData:    json.RawMessage(`{{}`),
			ErrorMsg:    "error calling MarshalJSON",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			if tt.JSONData == nil {
				jsonData, err := gofakeit.JSON(nil)
				tt.JSONData = json.RawMessage(jsonData)
				require.NoError(t, err)
			}

			resp, err := graphTestClient().AnnotationUpdate(ctx, testclient.AnnotationUpdateInput{NodeID: tt.NodeID, NamespaceID: tt.NamespaceID, Data: tt.JSONData})

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp.AnnotationUpdate.Annotation)
			assert.JSONEq(t, string(tt.JSONData), string(resp.AnnotationUpdate.Annotation.Data))

			antCount := EntClient.Annotation.Query().Where(annotation.AnnotationNamespaceID(tt.NamespaceID), annotation.HasMetadataWith(metadata.NodeID(tt.NodeID))).CountX(ctx)
			assert.Equal(t, 1, antCount)
		})
	}
}

func TestAnnotationDelete(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	ctx = perms.ContextWithHandler(ctx)

	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

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
			TestName:    "Fails when the annotation doesn't exist",
			NodeID:      meta1.NodeID,
			NamespaceID: AnnotationNamespaceBuilder{}.MustNew(ctx).ID,
			ErrorMsg:    "annotation not found",
		},
		{
			TestName:    "Fails when namespace is empty",
			NodeID:      gidx.MustNewID("testing"),
			NamespaceID: "",
			ErrorMsg:    "must not be empty",
		},
		{
			TestName:    "Fails when namespace gidx is invalid",
			NodeID:      gidx.MustNewID("testing"),
			NamespaceID: "test-invalid-id",
			ErrorMsg:    "invalid id",
		},
		{
			TestName:    "Fails when node is empty",
			NodeID:      "",
			NamespaceID: gidx.MustNewID("testing"),
			ErrorMsg:    "must not be empty",
		},
		{
			TestName:    "Fails when node gidx is invalid",
			NodeID:      "test-invalid-id",
			NamespaceID: gidx.MustNewID("testing"),
			ErrorMsg:    "invalid id",
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
