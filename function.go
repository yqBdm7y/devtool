package d

import (
	"crypto/rand"
	"encoding/base64"
)

// Example: u := d.ChangeType(usr, func(du database.User, mu *User) { mu.User = du })
func ChangeType[D, M any](database_value D, to_model_func func(D, *M)) M {
	var m M
	to_model_func(database_value, &m)
	return m
}

func ChangeTypeList[D, M any](database_list []D, to_model_func func(D, *M)) []M {
	var arr = []M{}
	for _, v := range database_list {
		arr = append(arr, ChangeType[D, M](v, to_model_func))
	}
	return arr
}

// Tree
type TreeInterface[T any] interface {
	IsTop() bool
	GetId() int
	GetParentId() int
	AppendChildren(T)
}

type tree_node[T any] struct {
	Data     T
	Children []*tree_node[T]
}

func GenerateTree[T any, PT interface {
	*T
	TreeInterface[T]
}](data []T) (treeData []T) {
	treeMap := make(map[int]*tree_node[T], len(data))
	for _, v := range data {
		treeMap[PT.GetId(&v)] = &tree_node[T]{Data: v}
	}

	var treeNodeList []*tree_node[T]
	for _, v := range treeMap {
		if PT.IsTop(&v.Data) {
			treeNodeList = append(treeNodeList, v)
		} else {
			if p, ok := treeMap[PT.GetParentId(&v.Data)]; ok {
				p.Children = append(p.Children, v)
			}
		}
	}
	treeData = recursion_tree_append_children[T, PT](treeNodeList)
	return treeData
}

func recursion_tree_append_children[T any, PT interface {
	*T
	TreeInterface[T]
}](list []*tree_node[T]) []T {
	var treeData []T
	for _, v := range list {
		if len(v.Children) != 0 {
			for _, v2 := range recursion_tree_append_children[T, PT](v.Children) {
				PT.AppendChildren(&v.Data, v2)
			}
		}
		treeData = append(treeData, v.Data)
	}
	return treeData
}

// Filter all the bottom-level leaf nodes
func FilterLeafNode[T any](slice []T, get_id func(T) int, get_parent_id func(T) int) []T {
	// Create a map to store whether each menu's ID is the parentID of other menus
	parentMap := make(map[int]bool)

	// Traverse slices and mark parentID
	for _, s := range slice {
		parentMap[get_parent_id(s)] = true
	}

	// Filter out the bottom slice
	var leafNodes []T
	for _, s := range slice {
		if !parentMap[get_id(s)] {
			leafNodes = append(leafNodes, s)
		}
	}

	return leafNodes
}

// Used to generate random strings
// Example: randomKey, err := GenerateRandomString(32)
// Return value example: tFredJ-Ii5Eh0hQAHaJXSSz8Ffd7S6xTY2s-ZMxOLCM=
func GenerateRandomString(length int) (string, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}
