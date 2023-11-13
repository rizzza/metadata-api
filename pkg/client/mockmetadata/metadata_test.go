package mockmetadata_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.infratographer.com/metadata-api/pkg/client/mockmetadata"

	metadata "go.infratographer.com/metadata-api/pkg/client"
)

func TestMetadata(t *testing.T) {
	t.Run("update status", func(t *testing.T) {
		mockMeta := new(mockmetadata.MockMetadata)

		mockMeta.On("StatusUpdate", context.Background(), &metadata.StatusUpdateInput{}).Return(&metadata.StatusUpdate{}, nil)

		resp, err := mockMeta.StatusUpdate(context.Background(), &metadata.StatusUpdateInput{})
		require.NoError(t, err)
		require.NotNil(t, resp)

		mockMeta.AssertExpectations(t)
	})
}
