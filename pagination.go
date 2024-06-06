package d

// Pagination interface, implement at least the following methods to facilitate internal calls in the devtool library
type InterfacePagination interface {
	Init()
	Set(page, page_size, total int, datalist interface{}) InterfacePagination
	ToMap() map[string]interface{}
}

var (
	FieldNamePaginationPage     = "page"
	FieldNamePaginationPageSize = "page_size"
	FieldNamePaginationTotal    = "total"
	FieldNamePaginationList     = "list"
)

var (
	pagination InterfacePagination // Global variable, stores the initialized interface, if not initialized, it is nil
)

// Pagination library unified access entry
type Pagination[T InterfacePagination] struct {
}

// Initialization
func (p Pagination[T]) Init(conf T) {
	pagination = conf
}

// Get the initialized interface. If it is not initialized, Pagination library is used by default.
func (p Pagination[T]) Get() T {
	if pagination == nil {
		LibraryPagination{}.Init()
	}
	// Value copy, changing the value will not affect the original value
	v, _ := pagination.(T)
	return v
}

// Pagination library
type LibraryPagination struct {
	Page     int
	PageSize int
	Total    int
	DataList interface{}
}

// Initialization
func (l LibraryPagination) Init() {
	Pagination[LibraryPagination]{}.Init(LibraryPagination{})
}

func (l LibraryPagination) Set(page, page_size, total int, datalist interface{}) InterfacePagination {
	return LibraryPagination{
		Page:     page,
		PageSize: page_size,
		Total:    total,
		DataList: datalist,
	}
}

// Pagination to map
func (l LibraryPagination) ToMap() map[string]interface{} {
	return map[string]interface{}{
		FieldNamePaginationPage:     l.Page,
		FieldNamePaginationPageSize: l.PageSize,
		FieldNamePaginationTotal:    l.Total,
		FieldNamePaginationList:     l.DataList,
	}
}
