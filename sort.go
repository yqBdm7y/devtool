package d

import (
	"errors"
	"sort"
)

type InterfaceSort interface {
	GetSort() int
	GetChildren() interface{}
}

func SortListWithChildrenBySortField[T InterfaceSort](list []T) error {
	sort.Slice(list, func(i, j int) bool {
		return list[i].GetSort() < list[j].GetSort()
	})

	for _, v := range list {
		children, ok := v.GetChildren().([]T)
		if !ok {
			return errors.New("incorrect parameter type")
		}
		if len(children) > 0 {
			SortListWithChildrenBySortField(children)
		}
	}
	return nil
}
