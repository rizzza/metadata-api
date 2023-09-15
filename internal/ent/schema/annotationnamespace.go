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

// AnnotationNamespace holds the schema definition for the AnnotationNamespace entity.
type AnnotationNamespace struct {
	ent.Schema
}

// Mixin of the Provider
func (AnnotationNamespace) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NewTimestampMixin(),
	}
}

// Fields of the Provider.
func (AnnotationNamespace) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(AnnotationNamespacePrefix) }).
			Unique().
			Immutable().
			Comment("The ID for the annotation namespace.").
			Annotations(
				entgql.OrderField("ID"),
			),
		field.String("name").
			NotEmpty().
			Comment("The name of the annotation namespace.").
			Annotations(
				entgql.OrderField("NAME"),
			),
		field.String("owner_id").
			GoType(gidx.PrefixedID("")).
			NotEmpty().
			Immutable().
			Comment("The ID for the owner for this annotation namespace.").
			Annotations(
				entgql.QueryField(),
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput, entgql.SkipType),
				entgql.OrderField("OWNER"),
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
func (AnnotationNamespace) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("owner_id"),
		index.Fields("owner_id", "name").Unique(),
	}
}

// Edges of the Provider
func (AnnotationNamespace) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("annotations", Annotation.Type).
			Ref("namespace").
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
			),
	}
}

// Annotations for the Provider
func (AnnotationNamespace) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		prefixIDDirective(AnnotationNamespacePrefix),
		schema.Comment("Representation of an annotation namespace. Annotation namespaces are used group annotation data that is provided by the same source and uses the same schema."),
		entgql.RelayConnection(),
		entgql.Mutations(
			entgql.MutationCreate().Description("Input information to create an annotation namespace."),
			entgql.MutationUpdate().Description("Input information to update an annotation namespace."),
		),
		entx.EventsHookSubjectName("annotation-namespace"),
	}
}
