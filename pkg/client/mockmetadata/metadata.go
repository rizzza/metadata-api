// Package mockmetadata
// Simplifying testing metadata node status updates in applications.
package mockmetadata

import (
	"context"

	"github.com/stretchr/testify/mock"

	metadata "go.infratographer.com/metadata-api/pkg/client"
)

// MockMetadata implements permissions.AuthRelationshipRequestHandler.
type MockMetadata struct {
	mock.Mock
}

// CreateAuthRelationships implements permissions.AuthRelationshipRequestHandler.
func (m *MockMetadata) StatusUpdate(ctx context.Context, input *metadata.StatusUpdateInput) (*metadata.StatusUpdate, error) {
	calledArgs := []interface{}{ctx, input}

	args := m.Called(calledArgs...)

	return args.Get(0).(*metadata.StatusUpdate), args.Error(1)
}
