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

func TestTenantAnnotationNamespaces(t *testing.T) {
	ctx := context.Background()

	perms := new(mockpermissions.MockPermissions)
	ctx = perms.ContextWithHandler(ctx)

	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ownerID := gidx.MustNewID("testing")
	antColors := AnnotationNamespaceBuilder{OwnerID: ownerID, Name: "nicole.dev/colors"}.MustNew(ctx)
	antPeople := AnnotationNamespaceBuilder{OwnerID: ownerID, Name: "nicole.dev/people"}.MustNew(ctx)
	antPlaces := AnnotationNamespaceBuilder{OwnerID: ownerID, Name: "nicole.dev/places"}.MustNew(ctx)
	AnnotationNamespaceBuilder{}.MustNew(ctx)
	// Update antColors so it's updated at is most recent
	antColors.Update().SaveX(ctx)

	testCases := []struct {
		TestName        string
		OrderBy         *testclient.AnnotationNamespaceOrder
		ResourceOwnerID gidx.PrefixedID
		ResponseOrder   []*ent.AnnotationNamespace
		errorMsg        string
	}{
		{
			TestName:        "Ordered By NAME ASC",
			OrderBy:         &testclient.AnnotationNamespaceOrder{Field: "NAME", Direction: "ASC"},
			ResourceOwnerID: ownerID,
			ResponseOrder:   []*ent.AnnotationNamespace{antColors, antPeople, antPlaces},
		},
		{
			TestName:        "Ordered By NAME DESC",
			OrderBy:         &testclient.AnnotationNamespaceOrder{Field: "NAME", Direction: "DESC"},
			ResourceOwnerID: ownerID,
			ResponseOrder:   []*ent.AnnotationNamespace{antPlaces, antPeople, antColors},
		},
		{
			TestName:        "Ordered By CREATED_AT ASC",
			OrderBy:         &testclient.AnnotationNamespaceOrder{Field: "CREATED_AT", Direction: "ASC"},
			ResourceOwnerID: ownerID,
			ResponseOrder:   []*ent.AnnotationNamespace{antColors, antPeople, antPlaces},
		},
		{
			TestName:        "Ordered By CREATED_AT DESC",
			OrderBy:         &testclient.AnnotationNamespaceOrder{Field: "CREATED_AT", Direction: "DESC"},
			ResourceOwnerID: ownerID,
			ResponseOrder:   []*ent.AnnotationNamespace{antPlaces, antPeople, antColors},
		},
		{
			TestName:        "Ordered By UPDATED_AT ASC",
			OrderBy:         &testclient.AnnotationNamespaceOrder{Field: "UPDATED_AT", Direction: "ASC"},
			ResourceOwnerID: ownerID,
			ResponseOrder:   []*ent.AnnotationNamespace{antPeople, antPlaces, antColors},
		},
		{
			TestName:        "Ordered By UPDATED_AT DESC",
			OrderBy:         &testclient.AnnotationNamespaceOrder{Field: "UPDATED_AT", Direction: "DESC"},
			ResourceOwnerID: ownerID,
			ResponseOrder:   []*ent.AnnotationNamespace{antColors, antPlaces, antPeople},
		},
		{
			TestName:        "No Annotation Namespaces for Tenant",
			ResourceOwnerID: gidx.MustNewID("testing"),
			ResponseOrder:   []*ent.AnnotationNamespace{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().GetResourceOwnerAnnotationNamespaces(ctx, tt.ResourceOwnerID, tt.OrderBy)

			if tt.errorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)

				return
			}

			require.Len(t, resp.Entities[0].AnnotationNamespaces.Edges, len(tt.ResponseOrder))
			for i, lb := range tt.ResponseOrder {
				respNS := resp.Entities[0].AnnotationNamespaces.Edges[i].Node
				assert.Equal(t, lb.ID, respNS.ID)
				assert.Equal(t, lb.Name, respNS.Name)
			}
		})
	}
}
