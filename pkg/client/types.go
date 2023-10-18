package client

import (
	"encoding/json"
)

type StatusUpdate struct {
	StatusUpdate StatusUpdateResponse `graphql:"statusUpdate(input: $input)"`
}

type StatusUpdateInput struct {
	// The node ID for this status.
	NodeID string `graphql:"nodeID" json:"nodeID"`
	// The namespace ID for this status.
	NamespaceID string `graphql:"namespaceID" json:"namespaceID"`
	// The source for this status.
	Source string `graphql:"source" json:"source"`
	// The data to save in this status.
	Data json.RawMessage `graphql:"data" json:"data"`
}

type StatusUpdateResponse struct {
	Status struct {
		ID                string          `graphql:"id"`
		Data              json.RawMessage `graphql:"data"`
		Source            string          `graphql:"source"`
		StatusNamespaceID string          `graphql:"statusNamespaceID"`

		Metadata struct {
			ID     string `graphql:"id"`
			NodeID string `graphql:"nodeID"`
		} `graphql:"metadata"`
	} `graphql:"status"`
}
