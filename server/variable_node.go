package server

import (
	"context"
	"sync"

	"github.com/awcullen/opcua"
)

type VariableNode struct {
	sync.RWMutex
	nodeId                  opcua.NodeID
	nodeClass               opcua.NodeClass
	browseName              opcua.QualifiedName
	displayName             opcua.LocalizedText
	description             opcua.LocalizedText
	rolePermissions         []opcua.RolePermissionType
	accessRestrictions      uint16
	references              []opcua.Reference
	value                   opcua.DataValue
	dataType                opcua.NodeID
	valueRank               int32
	arrayDimensions         []uint32
	accessLevel             byte
	minimumSamplingInterval float64
	historizing             bool
	readValueHandler        func(context.Context, opcua.ReadValueID) opcua.DataValue
	writeValueHandler       func(context.Context, opcua.WriteValue) opcua.StatusCode
}

var _ Node = (*VariableNode)(nil)

func NewVariableNode(nodeID opcua.NodeID, browseName opcua.QualifiedName, displayName opcua.LocalizedText, description opcua.LocalizedText, rolePermissions []opcua.RolePermissionType, references []opcua.Reference, value opcua.DataValue, dataType opcua.NodeID, valueRank int32, arrayDimensions []uint32, accessLevel byte, minimumSamplingInterval float64, historizing bool) *VariableNode {
	return &VariableNode{
		nodeId:                  nodeID,
		nodeClass:               opcua.NodeClassVariable,
		browseName:              browseName,
		displayName:             displayName,
		description:             description,
		rolePermissions:         rolePermissions,
		accessRestrictions:      0,
		references:              references,
		value:                   value,
		dataType:                dataType,
		valueRank:               valueRank,
		arrayDimensions:         arrayDimensions,
		accessLevel:             accessLevel,
		minimumSamplingInterval: minimumSamplingInterval,
		historizing:             historizing,
	}
}

// NodeID returns the NodeID attribute of this node.
func (n *VariableNode) NodeID() opcua.NodeID {
	return n.nodeId
}

// NodeClass returns the NodeClass attribute of this node.
func (n *VariableNode) NodeClass() opcua.NodeClass {
	return n.nodeClass
}

// BrowseName returns the BrowseName attribute of this node.
func (n *VariableNode) BrowseName() opcua.QualifiedName {
	return n.browseName
}

// DisplayName returns the DisplayName attribute of this node.
func (n *VariableNode) DisplayName() opcua.LocalizedText {
	return n.displayName
}

// Description returns the Description attribute of this node.
func (n *VariableNode) Description() opcua.LocalizedText {
	return n.description
}

// RolePermissions returns the RolePermissions attribute of this node.
func (n *VariableNode) RolePermissions() []opcua.RolePermissionType {
	return n.rolePermissions
}

// UserRolePermissions returns the RolePermissions attribute of this node for the current user.
func (n *VariableNode) UserRolePermissions(ctx context.Context) []opcua.RolePermissionType {
	filteredPermissions := []opcua.RolePermissionType{}
	session, ok := ctx.Value(SessionKey).(*Session)
	if !ok {
		return filteredPermissions
	}
	roles := session.UserRoles()
	rolePermissions := n.RolePermissions()
	if rolePermissions == nil {
		rolePermissions = session.Server().RolePermissions()
	}
	for _, role := range roles {
		for _, rp := range rolePermissions {
			if rp.RoleID == role {
				filteredPermissions = append(filteredPermissions, rp)
			}
		}
	}
	return filteredPermissions
}

// References returns the References of this node.
func (n *VariableNode) References() []opcua.Reference {
	n.RLock()
	res := n.references
	n.RUnlock()
	return res
}

// SetReferences sets the References of the Variable.
func (n *VariableNode) SetReferences(value []opcua.Reference) {
	n.Lock()
	n.references = value
	n.Unlock()
}

// Value returns the value of the Variable.
func (n *VariableNode) Value() opcua.DataValue {
	n.RLock()
	res := n.value
	n.RUnlock()
	return res
}

// SetValue sets the value of the Variable.
func (n *VariableNode) SetValue(value opcua.DataValue) {
	n.Lock()
	n.value = value
	n.Unlock()
}

// DataType returns the DataType attribute of this node.
func (n *VariableNode) DataType() opcua.NodeID {
	return n.dataType
}

// ValueRank returns the ValueRank attribute of this node.
func (n *VariableNode) ValueRank() int32 {
	return n.valueRank
}

// ArrayDimensions returns the ArrayDimensions attribute of this node.
func (n *VariableNode) ArrayDimensions() []uint32 {
	return n.arrayDimensions
}

// AccessLevel returns the AccessLevel attribute of this node.
func (n *VariableNode) AccessLevel() byte {
	return n.accessLevel
}

// UserAccessLevel returns the AccessLevel attribute of this node for this user.
func (n *VariableNode) UserAccessLevel(ctx context.Context) byte {
	accessLevel := n.accessLevel
	session, ok := ctx.Value(SessionKey).(*Session)
	if !ok {
		return 0
	}
	roles := session.UserRoles()
	rolePermissions := n.RolePermissions()
	if rolePermissions == nil {
		rolePermissions = session.Server().RolePermissions()
	}
	var currentRead, currentWrite, historyRead bool
	for _, role := range roles {
		for _, rp := range rolePermissions {
			if rp.RoleID == role {
				if rp.Permissions&opcua.PermissionTypeRead != 0 {
					currentRead = true
				}
				if rp.Permissions&opcua.PermissionTypeWrite != 0 {
					currentWrite = true
				}
				if rp.Permissions&opcua.PermissionTypeReadHistory != 0 {
					historyRead = true
				}
			}
		}
	}
	if !currentRead {
		accessLevel &^= opcua.AccessLevelsCurrentRead
	}
	if !currentWrite {
		accessLevel &^= opcua.AccessLevelsCurrentWrite
	}
	if !historyRead {
		accessLevel &^= opcua.AccessLevelsHistoryRead
	}
	return accessLevel
}

// MinimumSamplingInterval returns the MinimumSamplingInterval attribute of this node.
func (n *VariableNode) MinimumSamplingInterval() float64 {
	return n.minimumSamplingInterval
}

// Historizing returns the Historizing attribute of this node.
func (n *VariableNode) Historizing() bool {
	n.RLock()
	ret := n.historizing
	n.RUnlock()
	return ret
}

// SetHistorizing sets the Historizing attribute of this node.
func (n *VariableNode) SetHistorizing(historizing bool) {
	n.Lock()
	n.historizing = historizing
	n.Unlock()
}

// SetReadValueHandler sets the ReadValueHandler of this node.
func (n *VariableNode) SetReadValueHandler(value func(context.Context, opcua.ReadValueID) opcua.DataValue) {
	n.Lock()
	n.readValueHandler = value
	n.Unlock()
}

// SetWriteValueHandler sets the WriteValueHandler of this node.
func (n *VariableNode) SetWriteValueHandler(value func(context.Context, opcua.WriteValue) opcua.StatusCode) {
	n.Lock()
	n.writeValueHandler = value
	n.Unlock()
}

// IsAttributeIDValid returns true if attributeId is supported for the node.
func (n *VariableNode) IsAttributeIDValid(attributeID uint32) bool {
	switch attributeID {
	case opcua.AttributeIDNodeID, opcua.AttributeIDNodeClass, opcua.AttributeIDBrowseName,
		opcua.AttributeIDDisplayName, opcua.AttributeIDDescription, opcua.AttributeIDRolePermissions,
		opcua.AttributeIDUserRolePermissions, opcua.AttributeIDValue, opcua.AttributeIDDataType,
		opcua.AttributeIDValueRank, opcua.AttributeIDArrayDimensions, opcua.AttributeIDAccessLevel,
		opcua.AttributeIDUserAccessLevel, opcua.AttributeIDMinimumSamplingInterval, opcua.AttributeIDHistorizing:
		return true
	default:
		return false
	}
}
