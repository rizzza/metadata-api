package graphapi_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/permissions-api/pkg/permissions/mockpermissions"
	"go.infratographer.com/x/gidx"

	ent "go.infratographer.com/metadata-api/internal/ent/generated"
	"go.infratographer.com/metadata-api/internal/testclient"
)

func TestLoadingMetadataNodeMetadata(t *testing.T) {
	ctx := context.Background()

	perms := new(mockpermissions.MockPermissions)
	ctx = perms.ContextWithHandler(ctx)

	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	meta := MetadataBuilder{}.MustNew(ctx)
	ant := AnnotationBuilder{Metadata: meta}.MustNew(ctx)

	testCases := []struct {
		TestName           string
		OrderBy            *testclient.AnnotationNamespaceOrder
		NodeID             gidx.PrefixedID
		ResponseAnnotation *ent.Annotation
		errorMsg           string
	}{
		{
			TestName:           "returns metadata for a given node",
			NodeID:             meta.NodeID,
			ResponseAnnotation: ant,
		},
		{
			TestName: "null is returned when there is no metadata for a node id",
			NodeID:   gidx.MustNewID("testing"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().GetNodeMetadata(ctx, tt.NodeID)

			if tt.errorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)

				return
			}

			require.NoError(t, err)
			require.Len(t, resp.Entities, 1)
			if tt.ResponseAnnotation == nil {
				assert.Nil(t, resp.Entities[0].Metadata)

				return
			}

			require.NotNil(t, resp.Entities[0].Metadata)
			require.Len(t, resp.Entities[0].Metadata.Annotations.Edges, 1)
			expectNS, err := tt.ResponseAnnotation.Namespace(ctx)
			require.NoError(t, err)
			respAnt := resp.Entities[0].Metadata.Annotations.Edges[0].Node

			assert.EqualValues(t, expectNS.Name, respAnt.Namespace.Name)
		})
	}
}
