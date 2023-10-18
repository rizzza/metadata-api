package graphapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"encoding/json"

	"go.infratographer.com/x/gidx"

	"go.infratographer.com/permissions-api/pkg/permissions"

	"go.infratographer.com/metadata-api/internal/ent/generated"
	"go.infratographer.com/metadata-api/internal/ent/generated/metadata"
	"go.infratographer.com/metadata-api/internal/ent/generated/status"
)

// StatusUpdate is the resolver for the statusUpdate field.
func (r *mutationResolver) StatusUpdate(ctx context.Context, input StatusUpdateInput) (*StatusUpdateResponse, error) {
	logger := r.logger.With("nodeID", input.NodeID, "namespaceID", input.NamespaceID, "source", input.Source)

	if input.NamespaceID == "" {
		return nil, &ErrInvalidField{field: "NamespaceID", err: ErrFieldEmpty}
	}

	if input.NodeID == "" {
		return nil, &ErrInvalidField{field: "NodeID", err: ErrFieldEmpty}
	}

	if input.Source == "" {
		return nil, &ErrInvalidField{field: "Source", err: ErrFieldEmpty}
	}

	if _, err := gidx.Parse(input.NodeID.String()); err != nil {
		return nil, &ErrInvalidField{field: "NodeID", err: err}
	}

	if _, err := gidx.Parse(input.NamespaceID.String()); err != nil {
		return nil, &ErrInvalidField{field: "NamespaceID", err: err}
	}

	if !json.Valid(input.Data) {
		return nil, &ErrInvalidField{field: "Data", err: ErrInvalidJSON}
	}

	if err := permissions.CheckAccess(ctx, input.NamespaceID, actionMetadataStatusNamespaceUpdate); err != nil {
		return nil, err
	}

	if _, err := r.client.StatusNamespace.Get(ctx, input.NamespaceID); err != nil {
		if generated.IsNotFound(err) {
			return nil, err
		}

		if generated.IsValidationError(err) {
			return nil, err
		}

		logger.Errorw("failed to get status namespace", "error", err)
		return nil, ErrInternalServerError
	}

	status, err := r.client.Status.Query().Where(
		status.HasMetadataWith(metadata.NodeID(input.NodeID)),
		status.StatusNamespaceID(input.NamespaceID),
		status.Source(input.Source),
	).First(ctx)
	if err != nil {
		md, err := r.client.Metadata.Query().Where(metadata.NodeID(input.NodeID)).First(ctx)
		if err != nil {
			if generated.IsNotFound(err) {
				md, err = r.client.Metadata.Create().SetNodeID(input.NodeID).Save(ctx)
				if err != nil {
					if generated.IsValidationError(err) {
						return nil, err
					}

					logger.Errorw("failed to create metadata", "error", err)
					return nil, ErrInternalServerError
				}
			} else {
				logger.Errorw("failed to get metadata", "error", err)
				return nil, ErrInternalServerError
			}
		}

		status, err = r.client.Status.Create().SetInput(generated.CreateStatusInput{
			MetadataID:  md.ID,
			NamespaceID: input.NamespaceID,
			Source:      input.Source,
			Data:        input.Data,
		}).Save(ctx)
		if err != nil {
			logger.Errorw("failed to create status", "error", err)
			return nil, ErrInternalServerError
		}

		return &StatusUpdateResponse{Status: status}, nil
	}

	status, err = status.Update().SetData(input.Data).Save(ctx)
	if err != nil {
		logger.Errorw("failed to update status", "error", err)
		return nil, ErrInternalServerError
	}

	return &StatusUpdateResponse{Status: status}, nil
}

// StatusDelete is the resolver for the statusDelete field.
func (r *mutationResolver) StatusDelete(ctx context.Context, input StatusDeleteInput) (*StatusDeleteResponse, error) {
	logger := r.logger.With("nodeID", input.NodeID, "namespaceID", input.NamespaceID, "source", input.Source)

	if input.NamespaceID == "" {
		return nil, &ErrInvalidField{field: "NamespaceID", err: ErrFieldEmpty}
	}

	if input.NodeID == "" {
		return nil, &ErrInvalidField{field: "NodeID", err: ErrFieldEmpty}
	}

	if input.Source == "" {
		return nil, &ErrInvalidField{field: "Source", err: ErrFieldEmpty}
	}

	if _, err := gidx.Parse(input.NodeID.String()); err != nil {
		return nil, &ErrInvalidField{field: "NodeID", err: err}
	}

	if _, err := gidx.Parse(input.NamespaceID.String()); err != nil {
		return nil, &ErrInvalidField{field: "NamespaceID", err: err}
	}

	if err := permissions.CheckAccess(ctx, input.NamespaceID, actionMetadataStatusNamespaceUpdate); err != nil {
		return nil, err
	}

	st, err := r.client.Status.Query().Where(
		status.HasMetadataWith(metadata.NodeID(input.NodeID)),
		status.StatusNamespaceID(input.NamespaceID),
		status.Source(input.Source),
	).First(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, err
		}

		logger.Errorw("failed to get status", "error", err)
		return nil, ErrInternalServerError
	}

	if err := r.client.Status.DeleteOne(st).Exec(ctx); err != nil {
		logger.Errorw("failed to delete status", "error", err)
		return nil, ErrInternalServerError
	}

	return &StatusDeleteResponse{DeletedID: st.ID}, nil
}
