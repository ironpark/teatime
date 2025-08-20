package recipe

import (
	"encoding/json"
	"fmt"
)

// Position represents a 2D coordinate as [x, y] array.
// It supports both YAML array format [x,y] and JSON object format {"x":x,"y":y}.
type Position [2]int

// MarshalYAML serializes Position to YAML array format [x,y].
func (p Position) MarshalYAML() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d,%d]", p[0], p[1])), nil
}

// MarshalJSON serializes Position to JSON object format {"x":x,"y":y}.
func (p Position) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"x":%d,"y":%d}`, p[0], p[1])), nil
}

// UnmarshalJSON deserializes JSON object format {"x":x,"y":y} into Position.
func (p *Position) UnmarshalJSON(data []byte) error {
	var position struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	if err := json.Unmarshal(data, &position); err != nil {
		return err
	}
	p[0], p[1] = position.X, position.Y
	return nil
}
