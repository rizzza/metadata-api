package graphapi_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/permissions-api/pkg/permissions/mockpermissions"
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/metadata-api/internal/testclient"
)

func TestStatusNamespacesCreate(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ns1 := StatusNamespaceBuilder{}.MustNew(ctx)

	testCases := []struct {
		TestName             string
		StatusNamespaceInput testclient.CreateStatusNamespaceInput
		ErrorMsg             string
	}{
		{
			TestName:             "Successful path",
			StatusNamespaceInput: testclient.CreateStatusNamespaceInput{Name: gofakeit.DomainName(), ResourceProviderID: gidx.MustNewID("testing")},
		},
		{
			TestName:             "Successful even when name is in use by another resource provider",
			StatusNamespaceInput: testclient.CreateStatusNamespaceInput{Name: ns1.Name, ResourceProviderID: gidx.MustNewID("tprefix")},
		},
		{
			TestName:             "Failed when name is in use by same resource provider",
			StatusNamespaceInput: testclient.CreateStatusNamespaceInput{Name: ns1.Name, ResourceProviderID: ns1.ResourceProviderID},
			ErrorMsg:             "constraint failed", // TODO: This should have a better error message
		},
		{
			TestName:             "Fails when resource provider is empty",
			StatusNamespaceInput: testclient.CreateStatusNamespaceInput{Name: ns1.Name, ResourceProviderID: ""},
			ErrorMsg:             "value is less than the required length", // TODO: This should have a better error message
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().StatusNamespaceCreate(ctx, tt.StatusNamespaceInput)

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp.StatusNamespaceCreate.StatusNamespace)
			assert.Equal(t, tt.StatusNamespaceInput.Name, resp.StatusNamespaceCreate.StatusNamespace.Name)
		})
	}
}

func TestStatusNamespacesDelete(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	perms.On("DeleteAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ns1 := StatusNamespaceBuilder{}.MustNew(ctx)
	ns2 := StatusNamespaceBuilder{}.MustNew(ctx)
	ns3 := StatusNamespaceBuilder{}.MustNew(ctx)

	StatusBuilder{StatusNamespace: ns1}.MustNew(ctx)
	StatusBuilder{StatusNamespace: ns2}.MustNew(ctx)
	StatusBuilder{StatusNamespace: ns2}.MustNew(ctx)

	testCases := []struct {
		TestName           string
		StatusNamespaceID  gidx.PrefixedID
		Force              bool
		StatusDeletedCount int64
		ErrorMsg           string
	}{
		{
			TestName:          "Fails when there are status' using it",
			StatusNamespaceID: ns1.ID,
			ErrorMsg:          "namespace is in use and can't be deleted",
		},
		{
			TestName:          "Successful when nothing is using it",
			StatusNamespaceID: ns3.ID,
		},
		{
			TestName:           "Successful even when it has status' if you force it",
			StatusNamespaceID:  ns2.ID,
			Force:              true,
			StatusDeletedCount: 2,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().StatusNamespaceDelete(ctx, tt.StatusNamespaceID, tt.Force)

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp.StatusNamespaceDelete)
			assert.Equal(t, tt.StatusNamespaceID, resp.StatusNamespaceDelete.DeletedID)
			assert.Equal(t, tt.StatusDeletedCount, resp.StatusNamespaceDelete.StatusDeletedCount)
		})
	}
}

func TestStatusNamespacesUpdate(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)

	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ns := StatusNamespaceBuilder{}.MustNew(ctx)
	ns2 := StatusNamespaceBuilder{ResourceProviderID: ns.ResourceProviderID}.MustNew(ctx)

	testCases := []struct {
		TestName string
		ID       gidx.PrefixedID
		NewName  string
		ErrorMsg string
	}{
		{
			TestName: "Successful path",
			ID:       StatusNamespaceBuilder{}.MustNew(ctx).ID,
			NewName:  gofakeit.DomainName(),
		},
		{
			TestName: "Successful even when name is in use by another tenant",
			ID:       StatusNamespaceBuilder{}.MustNew(ctx).ID,
			NewName:  ns.Name,
		},
		{
			TestName: "Failed when name is in use by same tenant",
			ID:       ns2.ID,
			NewName:  ns.Name,
			ErrorMsg: "constraint failed", // TODO: This should have a better error message
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().StatusNamespaceUpdate(ctx, tt.ID, testclient.UpdateStatusNamespaceInput{Name: &tt.NewName})

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp.StatusNamespaceUpdate.StatusNamespace)
			assert.Equal(t, tt.NewName, resp.StatusNamespaceUpdate.StatusNamespace.Name)
		})
	}
}
