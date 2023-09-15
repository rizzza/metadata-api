package schema

import (
	"encoding/json"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"go.infratographer.com/x/entx"
	"go.infratographer.com/x/gidx"
)

// Status holds the schema definition for the Status entity.
type Status struct {
	ent.Schema
}

// Mixin of the Status
func (Status) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NewTimestampMixin(),
	}
}

// Fields of the Status.
func (Status) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(StatusPrefix) }).
			Unique().
			Immutable(),
		field.String("metadata_id").
			Comment("ID of the metadata of this status").
			GoType(gidx.PrefixedID("")).
			NotEmpty().
			Immutable().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput),
				// TODO: implement a way to add annotation to specify additional subjects
				// entx.EventsHookAdditionalSubject(),
			),
		field.String("status_namespace_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			NotEmpty().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput),
				// TODO: implement a way to add annotation to specify additional subjects
				// entx.EventsHookAdditionalSubject(),
			),
		field.String("source").
			Immutable().
			NotEmpty().
			Annotations(
				entgql.Skip(entgql.SkipMutationUpdateInput),
			),
		field.JSON("data", json.RawMessage{}).
			Comment("JSON formatted data of this annotation.").
			StorageKey("json_data").
			Annotations(
				entgql.Type("JSON"),
			),
	}
}

// Indexes of the Status
func (Status) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("metadata_id", "status_namespace_id"),
		index.Fields("metadata_id", "status_namespace_id", "source").Unique(),
		index.Fields("status_namespace_id", "data").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
		),
	}
}

// Edges of the Status
func (Status) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("namespace", StatusNamespace.Type).
			Unique().
			Required().
			Immutable().
			Field("status_namespace_id").
			Annotations(),
		edge.To("metadata", Metadata.Type).
			Unique().
			Required().
			Immutable().
			Field("metadata_id").
			Annotations(),
	}
}

// Annotations for the Status
func (Status) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		prefixIDDirective(StatusPrefix),
		entgql.RelayConnection(),
		entgql.Mutations(
			entgql.MutationCreate().Description("Input information to create a status namespace."),
			entgql.MutationUpdate().Description("Input information to update a status namespace."),
		),
		entx.EventsHookSubjectName("status"),
	}
}
