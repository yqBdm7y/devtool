package d

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
