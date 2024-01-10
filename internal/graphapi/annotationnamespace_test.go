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

	ent "go.infratographer.com/metadata-api/internal/ent/generated"
	"go.infratographer.com/metadata-api/internal/testclient"
)

func TestAnnotationNamespacesCreate(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ns1 := AnnotationNamespaceBuilder{}.MustNew(ctx)

	testCases := []struct {
		TestName                 string
		AnnotationNamespaceInput testclient.CreateAnnotationNamespaceInput
		ErrorMsg                 string
	}{
		{
			TestName:                 "Successful path",
			AnnotationNamespaceInput: testclient.CreateAnnotationNamespaceInput{Name: gofakeit.DomainName(), OwnerID: gidx.MustNewID("testing")},
		},
		{
			TestName:                 "Successful even when name is in use by another owner",
			AnnotationNamespaceInput: testclient.CreateAnnotationNamespaceInput{Name: ns1.Name, OwnerID: gidx.MustNewID("tprefix")},
		},
		{
			TestName:                 "Failed when name is in use by same owner",
			AnnotationNamespaceInput: testclient.CreateAnnotationNamespaceInput{Name: ns1.Name, OwnerID: ns1.OwnerID},
			ErrorMsg:                 "must be unique",
		},
		{
			TestName:                 "Fails when owner is empty",
			AnnotationNamespaceInput: testclient.CreateAnnotationNamespaceInput{Name: ns1.Name, OwnerID: ""},
			ErrorMsg:                 "must not be empty",
		},
		{
			TestName:                 "Fails when name is empty",
			AnnotationNamespaceInput: testclient.CreateAnnotationNamespaceInput{Name: "", OwnerID: ns1.OwnerID},
			ErrorMsg:                 "must not be empty",
		},
		{
			TestName:                 "Fails when owner is an invalid gidx",
			AnnotationNamespaceInput: testclient.CreateAnnotationNamespaceInput{Name: ns1.Name, OwnerID: "test-invalid-id"},
			ErrorMsg:                 "invalid id",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().AnnotationNamespaceCreate(ctx, tt.AnnotationNamespaceInput)

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp.AnnotationNamespaceCreate.AnnotationNamespace)
			assert.Equal(t, tt.AnnotationNamespaceInput.Name, resp.AnnotationNamespaceCreate.AnnotationNamespace.Name)
		})
	}
}

func TestAnnotationNamespacesDelete(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)

	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	perms.On("DeleteAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ns1 := AnnotationNamespaceBuilder{}.MustNew(ctx)
	ns2 := AnnotationNamespaceBuilder{}.MustNew(ctx)
	ns3 := AnnotationNamespaceBuilder{}.MustNew(ctx)

	AnnotationBuilder{AnnotationNamespace: ns1}.MustNew(ctx)
	AnnotationBuilder{AnnotationNamespace: ns2}.MustNew(ctx)
	AnnotationBuilder{AnnotationNamespace: ns2}.MustNew(ctx)

	testCases := []struct {
		TestName               string
		AnnotationNamespaceID  gidx.PrefixedID
		Force                  bool
		AnnotationDeletedCount int64
		ErrorMsg               string
	}{
		{
			TestName:              "Fails when there are annotations using it",
			AnnotationNamespaceID: ns1.ID,
			ErrorMsg:              "namespace is in use and can't be deleted",
		},
		{
			TestName:              "Fails when name is empty",
			AnnotationNamespaceID: "",
			ErrorMsg:              "must not be empty",
		},
		{
			TestName:              "Fails when id is an invalid gidx",
			AnnotationNamespaceID: "test-invalid-id",
			ErrorMsg:              "invalid id",
		},
		{
			TestName:              "Fails when id is not found",
			AnnotationNamespaceID: gidx.MustNewID("testing"),
			ErrorMsg:              "not found",
		},
		{
			TestName:              "Successful when nothing is using it",
			AnnotationNamespaceID: ns3.ID,
		},
		{
			TestName:               "Successful even when it has annotations if you force it",
			AnnotationNamespaceID:  ns2.ID,
			Force:                  true,
			AnnotationDeletedCount: 2,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().AnnotationNamespaceDelete(ctx, tt.AnnotationNamespaceID, tt.Force)

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp.AnnotationNamespaceDelete)
			assert.Equal(t, tt.AnnotationNamespaceID, resp.AnnotationNamespaceDelete.DeletedID)
			assert.Equal(t, tt.AnnotationDeletedCount, resp.AnnotationNamespaceDelete.AnnotationDeletedCount)
		})
	}
}

func TestAnnotationNamespacesUpdate(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)

	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ns := AnnotationNamespaceBuilder{}.MustNew(ctx)
	ns2 := AnnotationNamespaceBuilder{OwnerID: ns.OwnerID}.MustNew(ctx)

	testCases := []struct {
		TestName string
		ID       gidx.PrefixedID
		NewName  string
		ErrorMsg string
	}{
		{
			TestName: "Successful path",
			ID:       AnnotationNamespaceBuilder{}.MustNew(ctx).ID,
			NewName:  gofakeit.DomainName(),
		},
		{
			TestName: "Successful even when name is in use by another tenant",
			ID:       AnnotationNamespaceBuilder{}.MustNew(ctx).ID,
			NewName:  ns.Name,
		},
		{
			TestName: "Failed when name is in use by same tenant",
			ID:       ns2.ID,
			NewName:  ns.Name,
			ErrorMsg: "must be unique",
		},
		{
			TestName: "Fails when id is empty",
			ID:       "",
			ErrorMsg: "must not be empty",
		},
		{
			TestName: "Fails when id is not found",
			ID:       gidx.MustNewID("testing"),
			NewName:  ns.Name,
			ErrorMsg: "not found",
		},
		{
			TestName: "Fails when id is an invalid gidx",
			ID:       "test-invalid-id",
			NewName:  ns.Name,
			ErrorMsg: "invalid id",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().AnnotationNamespaceUpdate(ctx, tt.ID, testclient.UpdateAnnotationNamespaceInput{Name: &tt.NewName})

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp.AnnotationNamespaceUpdate.AnnotationNamespace)
			assert.Equal(t, tt.NewName, resp.AnnotationNamespaceUpdate.AnnotationNamespace.Name)
		})
	}
}

func TestAnnotationNamespacesGet(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)

	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ns1 := AnnotationNamespaceBuilder{}.MustNew(ctx)

	testCases := []struct {
		TestName   string
		QueryID    gidx.PrefixedID
		ExpectedLB *ent.AnnotationNamespace
		ErrorMsg   string
	}{
		{
			TestName:   "Successful path",
			QueryID:    ns1.ID,
			ExpectedLB: ns1,
		},
		{
			TestName: "Fails when id is empty",
			QueryID:  "",
			ErrorMsg: "must not be empty",
		},
		{
			TestName: "Fails when id is an invalid gidx",
			QueryID:  "test-invalid-id",
			ErrorMsg: "invalid id",
		},
		{
			TestName: "Fails when id is not found",
			QueryID:  gidx.MustNewID("testing"),
			ErrorMsg: "not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().GetAnnotationNamespace(ctx, tt.QueryID)

			if tt.ErrorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.ErrorMsg)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.EqualValues(t, tt.ExpectedLB.Name, resp.AnnotationNamespace.Name)
		})
	}
}
