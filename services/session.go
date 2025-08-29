package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ironpark/teatime/internal/node"
	rc "github.com/ironpark/teatime/internal/recipe"
	"github.com/samber/lo"
)

type EditSession struct {
	ID           string     `json:"id"`
	LoadedAt     time.Time  `json:"loaded_at"`
	LastModified time.Time  `json:"last_modified"`
	NeedsSave    bool       `json:"needs_save"`
	Recipe       *rc.Recipe `json:"recipe"`
	lock         sync.RWMutex
	nodeCounter  map[string]int
}

func newSession(recipe *rc.Recipe) *EditSession {
	session := &EditSession{
		ID:           uuid.New().String(),
		LoadedAt:     time.Now(),
		LastModified: time.Now(),
		Recipe:       recipe,
		nodeCounter:  make(map[string]int),
	}
	for _, n := range recipe.Nodes {
		session.nodeCounter[n.NodeData.Ref]++
	}
	return session
}

func (s *EditSession) GetNodeCount(ref string) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.nodeCounter[ref]
}

func (s *EditSession) CreateNode(ref string, x, y int) (rc.Node, error) {
	createdNode, err := node.GetNodeByRef(ref)
	if err != nil {
		return rc.Node{}, err
	}
	lowerName := strings.ReplaceAll(strings.TrimSpace(strings.ToLower(createdNode.Name())), " ", "-")
	s.lock.Lock()
	defer s.lock.Unlock()
	nodeId := fmt.Sprintf("%s%d", lowerName, s.nodeCounter[ref])
	node := rc.Node{
		Id:       nodeId,
		Position: rc.Position{x, y},
		Type:     string(createdNode.Type()),
		NodeData: rc.NodeData{
			Ref:           ref,
			Icon:          createdNode.Icon(),
			Label:         createdNode.Name(),
			Name:          createdNode.Name(),
			NodeType:      string(createdNode.Type()),
			Description:   createdNode.Description(),
			Properties:    createdNode.GetProperties(node.PropertyContext{}),
			Outputs:       createdNode.GetOutput(node.PropertyContext{}),
			OutputHandles: createdNode.GetOutputHandles(node.PropertyContext{}),
		},
	}
	s.nodeCounter[ref]++
	s.Recipe.Nodes = append(s.Recipe.Nodes, node)
	s.NeedsSave = true
	s.LastModified = time.Now()
	return node, nil
}

func (s *EditSession) UpdateRecipe(recipe *rc.Recipe) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Recipe = recipe
	s.NeedsSave = true
	s.LastModified = time.Now()
	return nil
}

func (s *EditSession) UpdateNode(nodeId string, x, y int, label string, properties map[string]any) (rc.Node, error) {
	foundNode, err := s.GetNode(nodeId)
	if err != nil {
		return rc.Node{}, err
	}
	createdNode, err := node.GetNodeByRef(foundNode.NodeData.Ref)
	if err != nil {
		return rc.Node{}, err
	}
	foundNode.Position = rc.Position{x, y}
	newProps := createdNode.GetProperties(properties)
	for i, prop := range newProps {
		val, ok := properties[prop.Key]
		if ok {
			newProps[i].Value = val
		}
	}
	foundNode.NodeData = rc.NodeData{
		Ref:           foundNode.NodeData.Ref,
		Icon:          foundNode.NodeData.Icon,
		Name:          foundNode.NodeData.Name,
		Label:         label,
		NodeType:      foundNode.NodeData.NodeType,
		Description:   foundNode.NodeData.Description,
		Properties:    newProps,
		Outputs:       createdNode.GetOutput(properties),
		OutputHandles: createdNode.GetOutputHandles(properties),
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.NeedsSave = true
	s.LastModified = time.Now()
	for i, n := range s.Recipe.Nodes {
		if n.Id == nodeId {
			s.Recipe.Nodes[i] = foundNode
			break
		}
	}
	return foundNode, nil
}

func (s *EditSession) UpdateNodePosition(nodeId string, x, y int) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	for i, n := range s.Recipe.Nodes {
		if n.Id == nodeId {
			s.Recipe.Nodes[i].Position = rc.Position{x, y}
			s.NeedsSave = true
			break
		}
	}
	s.LastModified = time.Now()
	return nil
}

func (s *EditSession) GetNode(nodeId string) (rc.Node, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	foundNode, ok := lo.Find(s.Recipe.Nodes, func(n rc.Node) bool {
		return n.Id == nodeId
	})
	if !ok {
		return rc.Node{}, fmt.Errorf("node not found")
	}
	return foundNode, nil
}

func (s *EditSession) DeleteNode(nodeId string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, i, ok := lo.FindIndexOf(s.Recipe.Nodes, func(n rc.Node) bool {
		return n.Id == nodeId
	})
	if !ok {
		return fmt.Errorf("node not found")
	}
	s.Recipe.Nodes = append(s.Recipe.Nodes[:i], s.Recipe.Nodes[i+1:]...)
	s.NeedsSave = true
	s.LastModified = time.Now()
	return nil
}

func (s *EditSession) Save() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.NeedsSave {
		if err := s.Recipe.Save(); err != nil {
			return err
		}
		s.NeedsSave = false
	}
	return nil
}
