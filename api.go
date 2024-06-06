package d

import (
	"encoding/json"
)

// API interface, implement at least the following methods to facilitate internal calls in the devtool library
type InterfaceApi interface {
	Init()
	Success() interface{}
	Error() interface{}
	Pagination(InterfacePagination) interface{}
	IsErrorResponse() bool
}

const (
	ConfigPathApiField = "api.field"
)

var (
	api InterfaceApi // Global variable, stores the initialized interface, if not initialized, it is nil
)

// Api library unified access entry
type Api[T InterfaceApi] struct{}

// Initialization
func (a Api[T]) Init(conf T) {
	api = conf
}

// Get the initialized interface. If it is not initialized, Turnstile library is used by default.
// Example:
// api := d.Api[d.LibraryApi]{}.Get()
// api.Response.Data = map[string]string{ "token":  token }
func (a Api[T]) Get() T {
	if api == nil {
		LibraryApi{}.Init()
	}
	// Value copy, changing the value will not affect the original value
	v, _ := api.(T)
	return v
}

// Api library
type LibraryApi struct {
	Response library_api_response
}

type library_api_response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

// Initialization
func (l LibraryApi) Init() {
	Api[LibraryApi]{}.Init(LibraryApi{})
}

// Returns the structure of a successful response
func (l LibraryApi) Success() interface{} {
	l.Response.Success = true
	if len(l.Response.Message) == 0 {
		l.Response.Message = "Success"
	}
	data, _ := l.ModifyApiFieldName(l.Response)
	return data
}

// Returns the structure of the error response
func (l LibraryApi) Error() interface{} {
	l.Response.Success = false
	if len(l.Response.Message) == 0 {
		l.Response.Message = "Error"
	}
	data, _ := l.ModifyApiFieldName(l.Response)
	return data
}

// Returns the structure of the pagination response
func (l LibraryApi) Pagination(p InterfacePagination) interface{} {
	l.Response.Success = true
	if len(l.Response.Message) == 0 {
		l.Response.Message = "Success"
	}
	l.Response.Data = p.ToMap()
	data, _ := l.ModifyApiFieldName(l.Response)
	return data
}

// Determine whether the current response is an error
func (l LibraryApi) IsErrorResponse() bool {
	if l.Response.Error != nil {
		return true
	}
	return false
}

// API interceptor, modify the returned fields
func (l LibraryApi) ModifyApiFieldName(data interface{}) (interface{}, error) {
	fieldMap := Config[InterfaceConfig]{}.Get().GetStringMap(ConfigPathApiField)
	if len(fieldMap) == 0 {
		return data, nil
	}

	b, err := json.Marshal(data)
	if err != nil {
		return data, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return data, err
	}

	l.recursionReplaceResponseKeys(m, fieldMap)

	return m, nil
}

// Recursive method to determine whether the structure is a slice or an object and modify it
func (l LibraryApi) recursionReplaceResponseKeys(jsonMap interface{}, fieldMap map[string]interface{}) {
	switch jv := jsonMap.(type) {
	case map[string]interface{}:
		for k, v := range jv {
			switch jv2 := v.(type) {
			case map[string]interface{}:
				l.recursionReplaceResponseKeys(jv2, fieldMap)
			case []interface{}:
				l.recursionReplaceResponseKeys(jv2, fieldMap)
			default:
				if fv, ok := fieldMap[k]; ok {
					jv[fv.(string)] = v
					delete(jv, k)
				}
			}
		}
	case []interface{}:
		for k, _ := range jv {
			l.recursionReplaceResponseKeys(jv[k], fieldMap)
		}
	}
}
