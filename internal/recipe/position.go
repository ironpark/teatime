package recipe

import (
	"encoding/json"
	"fmt"
)

type Position [2]int

func (p Position) MarshalYAML() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d,%d]", p[0], p[1])), nil
}

func (p Position) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"x":%d,"y":%d}`, p[0], p[1])), nil
}

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
