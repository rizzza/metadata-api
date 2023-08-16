package graphapi_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/x/gidx"

	testclient "go.infratographer.com/metadata-api/internal/testclient"
)

func TestAnnotationNamespacesCreate(t *testing.T) {
	ctx := context.Background()
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
			ErrorMsg:                 "constraint failed", // TODO: This should have a better error message
		},
		{
			TestName:                 "Fails when owner is empty",
			AnnotationNamespaceInput: testclient.CreateAnnotationNamespaceInput{Name: ns1.Name, OwnerID: ""},
			ErrorMsg:                 "value is less than the required length", // TODO: This should have a better error message
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
			ErrorMsg: "constraint failed", // TODO: This should have a better error message
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
