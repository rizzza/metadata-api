package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateStatus(t *testing.T) {
	cli := Client{}
	ctx := context.Background()

	t.Run("unauthorized", func(t *testing.T) {
		respJSON := `{"message":"invalid or expired jwt"}`

		cli.gqlCli = mustNewGQLTestClient(respJSON, http.StatusUnauthorized)

		lb, err := cli.StatusUpdate(context.Background(), &StatusUpdateInput{})
		require.Nil(t, lb)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("permission denied", func(t *testing.T) {
		respJSON := `{"message":"subject doesn't have access"}`

		cli.gqlCli = mustNewGQLTestClient(respJSON, http.StatusForbidden)

		lb, err := cli.StatusUpdate(context.Background(), &StatusUpdateInput{})
		require.Nil(t, lb)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrPermissionDenied)
	})

	t.Run("fails to update with bad NodeID prefix", func(t *testing.T) {
		respJSON := `{
			"errors": [
				{
				"message": "invalid id: expected prefix length is 7, 'badprefix' is 9, field: NodeID",
				"path": [
					"statusUpdate"
				]
				}
			],
			"data": null
			}`

		cli.gqlCli = mustNewGQLTestClient(respJSON, http.StatusOK)

		status, err := cli.StatusUpdate(ctx, &StatusUpdateInput{
			NodeID:      "badprefix-testing",
			NamespaceID: "metasns-testing",
		})

		require.Nil(t, status)
		require.Error(t, err)
		assert.ErrorContains(t, err, "invalid id")
	})

	t.Run("fails to update with bad NamespaceID prefix", func(t *testing.T) {
		respJSON := `{
			"errors": [
				{
				"message": "invalid id: expected prefix length is 7, 'badprefix' is 9, field: NamespaceID",
				"path": [
					"statusUpdate"
				]
				}
			],
			"data": null
			}`

		cli.gqlCli = mustNewGQLTestClient(respJSON, http.StatusOK)

		status, err := cli.StatusUpdate(ctx, &StatusUpdateInput{
			NodeID:      "loadbal-testing",
			NamespaceID: "badprefix-testing",
		})

		require.Nil(t, status)
		require.Error(t, err)
		assert.ErrorContains(t, err, "invalid id")
	})

	t.Run("fails to update with unknown NamespaceID", func(t *testing.T) {
		respJSON := `{
			"errors": [
				{
					"message": "generated: status_namespace not found",
					"path": [
						"statusUpdate"
					]
				}
			],
			"data": null
			}`

		cli.gqlCli = mustNewGQLTestClient(respJSON, http.StatusOK)

		status, err := cli.StatusUpdate(ctx, &StatusUpdateInput{
			NodeID:      "loadbal-testing",
			NamespaceID: "metasns-does-not-exist",
		})

		require.Nil(t, status)
		require.Error(t, err)
		assert.ErrorContains(t, err, "not found")
	})
	t.Run("successfully updates a node status", func(t *testing.T) {
		respJSON := `{
			"data": {
				"statusUpdate": {
					"status": {
						"id": "metasts-testing",
						"data": {"state":"ACTIVE"},
						"source": "unit-test",
						"statusNamespaceID": "metasns-testing",
						"metadata": {
							"id": "metadat-testing",
							"nodeID": "loadbal-test"
						}
					}
				}
			}
		}`

		cli.gqlCli = mustNewGQLTestClient(respJSON, http.StatusOK)

		resp, err := cli.StatusUpdate(ctx, &StatusUpdateInput{
			NodeID:      "loadbal-testing",
			NamespaceID: "metasts-testing",
			Source:      "unit-test",
			Data:        json.RawMessage(`{"state":"ACTIVE"}`),
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "metasts-testing", resp.StatusUpdate.Status.ID)
		assert.JSONEq(t, `{"state":"ACTIVE"}`, string(resp.StatusUpdate.Status.Data))
		assert.Equal(t, "unit-test", resp.StatusUpdate.Status.Source)
		assert.Equal(t, "metasns-testing", resp.StatusUpdate.Status.StatusNamespaceID)
		assert.Equal(t, "metadat-testing", resp.StatusUpdate.Status.Metadata.ID)
		assert.Equal(t, "loadbal-test", resp.StatusUpdate.Status.Metadata.NodeID)
	})
}

func mustNewGQLTestClient(respJSON string, respCode int) *graphql.Client {
	mux := http.NewServeMux()
	mux.HandleFunc("/query", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(respCode)
		w.Header().Set("Content-Type", "application/json")
		_, err := io.WriteString(w, respJSON)
		if err != nil {
			panic(err)
		}
	})

	return graphql.NewClient("/query", &http.Client{Transport: localRoundTripper{handler: mux}})
}

type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)

	return w.Result(), nil
}
