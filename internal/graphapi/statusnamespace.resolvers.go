package graphapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"database/sql"
	"fmt"

	"go.infratographer.com/metadata-api/internal/ent/generated"
	"go.infratographer.com/metadata-api/internal/ent/generated/status"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/x/gidx"
)

// StatusNamespaceCreate is the resolver for the statusNamespaceCreate field.
func (r *mutationResolver) StatusNamespaceCreate(ctx context.Context, input generated.CreateStatusNamespaceInput) (*StatusNamespaceCreatePayload, error) {
	if err := permissions.CheckAccess(ctx, input.ResourceProviderID, actionMetadataStatusNamespaceUpdate); err != nil {
		return nil, err
	}

	ns, err := r.client.StatusNamespace.Create().SetInput(input).Save(ctx)
	if err != nil {
		return nil, err
	}

	return &StatusNamespaceCreatePayload{StatusNamespace: ns}, nil
}

// StatusNamespaceDelete is the resolver for the statusNamespaceDelete field.
func (r *mutationResolver) StatusNamespaceDelete(ctx context.Context, id gidx.PrefixedID, force bool) (*StatusNamespaceDeletePayload, error) {
	sns, err := r.client.StatusNamespace.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := permissions.CheckAccess(ctx, sns.ResourceProviderID, actionMetadataStatusNamespaceUpdate); err != nil {
		return nil, err
	}

	tx, err := r.client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	statuses, err := tx.Status.Query().Where(status.StatusNamespaceID(id)).All(ctx)
	if err != nil {
		return nil, err
	}

	statusCount := len(statuses)
	if statusCount != 0 {
		if force {
			statusCount = 0
			for _, status := range statuses {
				// TODO - :bug: - must delete one-by-one to ensure the deleted ID is available when the delete eventhook is triggered
				// statusCount, err = r.client.Status.Delete().Where(status.StatusNamespaceID(id)).Exec(ctx)
				if err := tx.Status.DeleteOneID(status.ID).Exec(ctx); err != nil {
					return nil, err
				}
				statusCount++
			}
		} else {
			return nil, fmt.Errorf("status namespace is in use and can't be deleted")
		}
	}

	if err := tx.StatusNamespace.DeleteOneID(id).Exec(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &StatusNamespaceDeletePayload{DeletedID: id, StatusDeletedCount: statusCount}, nil
}

// StatusNamespaceUpdate is the resolver for the statusNamespaceUpdate field.
func (r *mutationResolver) StatusNamespaceUpdate(ctx context.Context, id gidx.PrefixedID, input generated.UpdateStatusNamespaceInput) (*StatusNamespaceUpdatePayload, error) {
	sns, err := r.client.StatusNamespace.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := permissions.CheckAccess(ctx, sns.ResourceProviderID, actionMetadataStatusNamespaceUpdate); err != nil {
		return nil, err
	}

	ns, err := r.client.StatusNamespace.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	ns, err = ns.Update().SetInput(input).Save(ctx)
	if err != nil {
		return nil, err
	}

	return &StatusNamespaceUpdatePayload{StatusNamespace: ns}, nil
}
