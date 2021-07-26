package schema

import "entgo.io/ent"

// Group holds the schema definition for the Group entity.
type Group struct {
	ent.Schema
}

// Fields of the Group.
func (Group) Fields() []ent.Field {
	return nil
}

// Edges of the Group.
func (Group) Edges() []ent.Edge {
	return nil
}
