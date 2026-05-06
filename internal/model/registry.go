package model

// Registry holds all GORM models to be included in auto-migration.
var Registry []interface{}

// Register adds a model to the migration registry.
func Register(m interface{}) {
	Registry = append(Registry, m)
}
