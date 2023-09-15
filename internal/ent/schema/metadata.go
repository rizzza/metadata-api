package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"go.infratographer.com/x/entx"
	"go.infratographer.com/x/gidx"
)

// Metadata holds the schema definition for the Metadata entity.
type Metadata struct {
	ent.Schema
}

// Mixin of the Metadata
func (Metadata) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NewTimestampMixin(),
	}
}

// Fields of the Metadata.
func (Metadata) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("ID for the metadata.").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(MetadataPrefix) }).
			Unique().
			Immutable(),
		field.String("node_id").
			Comment("ID of the node for this metadata").
			GoType(gidx.PrefixedID("")).
			NotEmpty().
			Unique().
			Immutable().
			Annotations(
				// TODO: implement a way to add annotation to specify additional subjects
				// entx.EventsHookAdditionalSubject(),
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput),
			),
	}
}

// Indexes of the Metadata
func (Metadata) Indexes() []ent.Index {
	return []ent.Index{}
}

// Edges of the Metadata
func (Metadata) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("annotations", Annotation.Type).
			Ref("metadata").
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				entgql.RelayConnection(),
			),

		edge.From("statuses", Status.Type).
			Ref("metadata").
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				entgql.RelayConnection(),
			),
	}
}

// Annotations for the Metadata
func (Metadata) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		entx.GraphKeyDirective("nodeID"),
		entx.EventsHookSubjectName("metadata"),
		prefixIDDirective(MetadataPrefix),
		entgql.RelayConnection(),
	}
}
