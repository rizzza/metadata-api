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
	"github.com/vektah/gqlparser/v2/ast"
	"go.infratographer.com/x/entx"
	"go.infratographer.com/x/gidx"
)

// Annotation holds the schema definition for the Annotation entity.
type Annotation struct {
	ent.Schema
}

// Mixin of the Annotation
func (Annotation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NewTimestampMixin(),
	}
}

// Fields of the Annotation.
func (Annotation) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("ID for the annotation.").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(AnnotationPrefix) }).
			Unique().
			Immutable(),
		field.String("metadata_id").
			Comment("ID of the metadata of this annotation").
			GoType(gidx.PrefixedID("")).
			Immutable().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput),
				entx.EventsHookAdditionalSubject(),
			),
		field.String("annotation_namespace_id").
			Comment("ID of the AnnotationNamespace of this annotation.").
			GoType(gidx.PrefixedID("")).
			NotEmpty().
			Immutable().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(^entgql.SkipMutationUpdateInput),
				entx.EventsHookAdditionalSubject(),
			),
		field.JSON("data", json.RawMessage{}).
			Comment("JSON formatted data of this annotation.").
			StorageKey("json_data").
			Annotations(
				entgql.Type("JSON"),
			),
	}
}

// Indexes of the Annotation
func (Annotation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("metadata_id", "annotation_namespace_id").Unique(),
		index.Fields("annotation_namespace_id", "data").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
		),
	}
}

// Edges of the Annotation
func (Annotation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("namespace", AnnotationNamespace.Type).
			Unique().
			Required().
			Immutable().
			Field("annotation_namespace_id").
			Annotations(),
		edge.To("metadata", Metadata.Type).
			Unique().
			Required().
			Immutable().
			Field("metadata_id").
			Annotations(),
	}
}

// Annotations for the Annotation
func (Annotation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		prefixIDDirective(AnnotationPrefix),
		entgql.RelayConnection(),
		entx.EventsHookSubjectName("annotation"),
	}
}

func prefixIDDirective(prefix string) entgql.Annotation {
	var args []*ast.Argument
	if prefix != "" {
		args = append(args, &ast.Argument{
			Name: "prefix",
			Value: &ast.Value{
				Raw:  prefix,
				Kind: ast.StringValue,
			},
		})
	}

	return entgql.Directives(entgql.NewDirective("prefixedID", args...))
}
