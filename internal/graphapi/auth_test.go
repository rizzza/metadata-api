package graphapi_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/permissions-api/pkg/permissions/mockpermissions"
	"go.infratographer.com/x/echojwtx"
	"go.infratographer.com/x/testing/auth"

	"go.infratographer.com/metadata-api/internal/testclient"
)

func TestJWTAnnotationNSCreateWithAuthClient(t *testing.T) {
	oauthCLI, issuer, oAuthClose := auth.OAuthTestClient("urn:test:status", "")
	defer oAuthClose()

	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	srv, err := newTestServer(
		withAuthConfig(
			&echojwtx.AuthConfig{
				Issuer: issuer,
			},
		),
		withPermissions(
			permissions.WithDefaultChecker(permissions.DefaultAllowChecker),
		),
	)

	require.NoError(t, err)
	require.NotNil(t, srv)

	defer srv.Close()

	resp, err := graphTestClient(
		withGraphClientHTTPClient(oauthCLI),
		withGraphClientServerURL(srv.URL+"/query"),
	).AnnotationNamespaceCreate(ctx, testclient.CreateAnnotationNamespaceInput{Name: "test", OwnerID: "tnttent-test"})

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "test", resp.AnnotationNamespaceCreate.AnnotationNamespace.Name)
	assert.Equal(t, "metamns", resp.AnnotationNamespaceCreate.AnnotationNamespace.ID.Prefix())
	assert.False(t, resp.AnnotationNamespaceCreate.AnnotationNamespace.Private)
}

func TestJWTAnnotationNSGetWithDefaultClient(t *testing.T) {
	_, issuer, oAuthClose := auth.OAuthTestClient("urn:test:loadbalancer", "")
	defer oAuthClose()

	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	srv, err := newTestServer(
		withAuthConfig(
			&echojwtx.AuthConfig{
				Issuer: issuer,
			},
		),
		withPermissions(
			permissions.WithDefaultChecker(permissions.DefaultAllowChecker),
		),
	)

	require.NoError(t, err)
	require.NotNil(t, srv)

	defer srv.Close()

	resp, err := graphTestClient(
		withGraphClientHTTPClient(http.DefaultClient),
		withGraphClientServerURL(srv.URL+"/query"),
	).AnnotationNamespaceCreate(ctx, testclient.CreateAnnotationNamespaceInput{Name: "test", OwnerID: "tnttent-test"})

	require.Error(t, err, "Expected an authorization error")
	require.Nil(t, resp)
	assert.ErrorContains(t, err, `{"networkErrors":{"code":401`)
}
