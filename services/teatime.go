package services

import (
	"github.com/ironpark/teatime/stores"
)

type TeaTime struct {
	store *stores.Store
}

func NewTeaTime(store *stores.Store) *TeaTime {
	return &TeaTime{store: store}
}

func (t *TeaTime) GetTeaTime() string {
	return "Tea Time!"
}

func (t *TeaTime) GetNodeInfos() []stores.NodeInfo {
	return t.store.GetNodeInfos()
}

func (t *TeaTime) GetNodeInfosByType(nodeType string) []stores.NodeInfo {
	return t.store.GetNodeInfosByType(nodeType)
}

func (t *TeaTime) GetNodeInfo(id string) stores.NodeInfo {
	return t.store.GetNodeInfo(id)
}

func (t *TeaTime) CreateNode(nodeId string) stores.Node {
	return t.store.CreateNode(nodeId)
}
