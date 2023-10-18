package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"go.infratographer.com/x/entx"
	"go.infratographer.com/x/gidx"
)

// StatusNamespace holds the schema definition for the StatusNamespace entity.
type StatusNamespace struct {
	ent.Schema
}

// Mixin of the Provider
func (StatusNamespace) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NewTimestampMixin(),
	}
}

// Fields of the Provider.
func (StatusNamespace) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(StatusNamespacePrefix) }).
			Unique().
			Immutable().
			Comment("The ID for the status namespace.").
			Annotations(
				entgql.OrderField("ID"),
			),
		field.String("name").
			NotEmpty().
			Comment("The name of the status namespace.").
			Annotations(
				entgql.OrderField("NAME"),
			),
		field.String("resource_provider_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			NotEmpty().
			Comment("The ID of the resource provider for this status namespace.").
			Annotations(
				entgql.QueryField(),
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput, entgql.SkipType),
				entgql.OrderField("RESOURCEPROVIDER"),
				entx.EventsHookAdditionalSubject("owner"),
			),
		field.Bool("private").
			Default(false).
			Comment("Flag for if this namespace is private.").
			Annotations(
				entgql.QueryField(),
				entgql.Skip(entgql.SkipWhereInput),
				entgql.OrderField("PRIVATE"),
			),
	}
}

// Indexes of the Provider
func (StatusNamespace) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("resource_provider_id"),
		index.Fields("resource_provider_id", "name").Unique(),
	}
}

// Edges of the Provider
func (StatusNamespace) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("statuses", Status.Type).
			Ref("namespace").
			Annotations(
				entgql.Skip(entgql.SkipAll),
			),
	}
}

// Annotations for the Provider
func (StatusNamespace) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		prefixIDDirective(StatusNamespacePrefix),
		schema.Comment("Representation of a status namespace. Status namespaces are used group status data that is provided by a resource provider."),
		entgql.RelayConnection(),
		entgql.Mutations(
			entgql.MutationCreate().Description("Input information to create a status namespace."),
			entgql.MutationUpdate().Description("Input information to update a status namespace."),
		),
		entx.EventsHookSubjectName("status-namespace"),
	}
}
