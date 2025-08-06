package types

import "errors"

type PropertyType int

const (
	// Basic Types
	Invalid PropertyType = iota
	Bool
	Int64
	Uint64
	Float64
	String
	// Special Types
	JSON
	XML
	Date
	Text
	// Arrays
	TextArray
	NumberArray
	BooleanArray
	JSONArray
	XMLArray
)

func (p PropertyType) String() string {
	switch p {
	case Invalid:
		return "Invalid"
	case Bool:
		return "Bool"
	case Int64:
		return "Int64"
	case Uint64:
		return "Uint64"
	case Float64:
		return "Float64"
	case String:
		return "String"
	case JSON:
		return "JSON"
	case XML:
		return "XML"
	case Date:
		return "Date"
	case Text:
		return "Text"
	case TextArray:
		return "TextArray"
	case NumberArray:
		return "NumberArray"
	case BooleanArray:
		return "BooleanArray"
	case JSONArray:
		return "JSONArray"
	case XMLArray:
		return "XMLArray"
	default:
		return "Unknown"
	}
}

type NodeProperty struct {
	// Type of the property (string, number, boolean, array, json, xml)
	Type PropertyType `json:"type"`
	// Name of the property
	Name string `json:"name"`
	// Description of the property
	Description string `json:"description"`
	// Whether the property is required
	Optional bool `json:"optional"`
	// Key of the property
	Key string `json:"key"`
	// Value of the property
	Value string `json:"value"`
	// Options of the property
	Options []string `json:"options"`
	// not editable
	ReadOnly bool `json:"readOnly"`
}

func (p *NodeProperty) Validate() error {
	if p.Type == Invalid {
		return errors.New("invalid property type")
	}

	switch p.Type {
	case String:
		if p.Options != nil {
			return errors.New("string property cannot have options")
		}
	case Int64:
		if p.Options != nil {
			return errors.New("int64 property cannot have options")
		}
	case Uint64:
		if p.Options != nil {
			return errors.New("uint64 property cannot have options")
		}
	}
	return nil
}
