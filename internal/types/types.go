package types

// TypeMapping defines the mapping between YAML types and their implementations
type TypeMapping struct {
	GoType         string
	DefaultGoValue string
	PostgresType   string
}

// SupportedTypes contains all supported type mappings
var SupportedTypes = map[string]TypeMapping{
	"uuid": {
		GoType:         "uuid.UUID",
		DefaultGoValue: "uuid.New()",
		PostgresType:   "uuid",
	},
	"string": {
		GoType:         "string",
		DefaultGoValue: "`test_string`",
		PostgresType:   "text",
	},
	"time": {
		GoType:         "time.Time",
		DefaultGoValue: "time.Now()",
		PostgresType:   "timestamptz",
	},
	"int": {
		GoType:         "int",
		DefaultGoValue: "0",
		PostgresType:   "integer",
	},
	"bool": {
		GoType:         "bool",
		DefaultGoValue: "false",
		PostgresType:   "boolean",
	},
}

// Field represents an entity field with its properties
type Field struct {
	Name           string `yaml:"name"`
	Type           string `yaml:"type"`
	Primary        bool   `yaml:"primary"`
	Owner          bool   `yaml:"owner"`
	GoType         string
	DefaultGoValue string
	PostgresType   string
}

// Entity represents the main entity structure
type Entity struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`
}

// PrimaryField returns the primary key field
func (e *Entity) PrimaryField() *Field {
	for _, field := range e.Fields {
		if field.Primary {
			return &field
		}
	}
	return nil
}

// OwnerField returns the owner field
func (e *Entity) OwnerField() *Field {
	for _, field := range e.Fields {
		if field.Owner {
			return &field
		}
	}
	return nil
}

// CommonFields returns all non-primary and non-owner fields
func (e *Entity) CommonFields() []*Field {
	common := make([]*Field, 0, len(e.Fields))
	for _, field := range e.Fields {
		if field.Primary || field.Owner {
			continue
		}
		common = append(common, &field)
	}
	return common
}
