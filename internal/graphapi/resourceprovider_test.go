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
	testclient "go.infratographer.com/metadata-api/internal/testclient"
)

func TestResourceProviderStatusNamespaces(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)

	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	rpID := gidx.MustNewID("testing")
	stColors := StatusNamespaceBuilder{ResourceProviderID: rpID, Name: "instance.infratographer.com/colors"}.MustNew(ctx)
	stPeople := StatusNamespaceBuilder{ResourceProviderID: rpID, Name: "instance.infratographer.com/people"}.MustNew(ctx)
	stPlaces := StatusNamespaceBuilder{ResourceProviderID: rpID, Name: "instance.infratographer.com/places"}.MustNew(ctx)
	StatusNamespaceBuilder{}.MustNew(ctx)
	// Update stColors so it's updated at is most recent
	stColors.Update().SaveX(ctx)

	testCases := []struct {
		TestName           string
		OrderBy            *testclient.StatusNamespaceOrder
		ResourceProviderID gidx.PrefixedID
		ResponseOrder      []*ent.StatusNamespace
		errorMsg           string
	}{
		{
			TestName:           "Ordered By NAME ASC",
			OrderBy:            &testclient.StatusNamespaceOrder{Field: "NAME", Direction: "ASC"},
			ResourceProviderID: rpID,
			ResponseOrder:      []*ent.StatusNamespace{stColors, stPeople, stPlaces},
		},
		{
			TestName:           "Ordered By NAME DESC",
			OrderBy:            &testclient.StatusNamespaceOrder{Field: "NAME", Direction: "DESC"},
			ResourceProviderID: rpID,
			ResponseOrder:      []*ent.StatusNamespace{stPlaces, stPeople, stColors},
		},
		{
			TestName:           "Ordered By CREATED_AT ASC",
			OrderBy:            &testclient.StatusNamespaceOrder{Field: "CREATED_AT", Direction: "ASC"},
			ResourceProviderID: rpID,
			ResponseOrder:      []*ent.StatusNamespace{stColors, stPeople, stPlaces},
		},
		{
			TestName:           "Ordered By CREATED_AT DESC",
			OrderBy:            &testclient.StatusNamespaceOrder{Field: "CREATED_AT", Direction: "DESC"},
			ResourceProviderID: rpID,
			ResponseOrder:      []*ent.StatusNamespace{stPlaces, stPeople, stColors},
		},
		{
			TestName:           "Ordered By UPDATED_AT ASC",
			OrderBy:            &testclient.StatusNamespaceOrder{Field: "UPDATED_AT", Direction: "ASC"},
			ResourceProviderID: rpID,
			ResponseOrder:      []*ent.StatusNamespace{stPeople, stPlaces, stColors},
		},
		{
			TestName:           "Ordered By UPDATED_AT DESC",
			OrderBy:            &testclient.StatusNamespaceOrder{Field: "UPDATED_AT", Direction: "DESC"},
			ResourceProviderID: rpID,
			ResponseOrder:      []*ent.StatusNamespace{stColors, stPlaces, stPeople},
		},
		{
			TestName:           "No Annotation Namespaces for Tenant",
			ResourceProviderID: gidx.MustNewID("testing"),
			ResponseOrder:      []*ent.StatusNamespace{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().GetResourceProviderStatusNamespaces(ctx, tt.ResourceProviderID, tt.OrderBy)

			if tt.errorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)

				return
			}

			require.NoError(t, err)
			require.Len(t, resp.Entities[0].StatusNamespaces.Edges, len(tt.ResponseOrder))
			for i, lb := range tt.ResponseOrder {
				respNS := resp.Entities[0].StatusNamespaces.Edges[i].Node
				assert.Equal(t, lb.ID, respNS.ID)
				assert.Equal(t, lb.Name, respNS.Name)
			}
		})
	}
}
