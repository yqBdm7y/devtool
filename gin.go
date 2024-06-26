package d

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Gin struct{}

// Returns a successful response in gin format
func (g Gin) Success(c *gin.Context, a InterfaceApi) {
	// If gin.Context is nil
	if c == nil {
		return
	}
	c.JSON(http.StatusOK, a.Success())
}

// Returns an error response in gin format
func (g Gin) Error(c *gin.Context, a InterfaceApi) {
	// If gin.Context is nil
	if c == nil {
		return
	}
	c.JSON(http.StatusOK, a.Error())
}

// Returns a pagination response in gin format
func (g Gin) Pagination(c *gin.Context, a InterfaceApi, p InterfacePagination) {
	// If gin.Context is nil
	if c == nil {
		return
	}
	c.JSON(http.StatusOK, a.Pagination(p))
}

// Returns data or error response in gin format
func (g Gin) DataOrError(c *gin.Context, a InterfaceApi) {
	// If gin.Context is nil
	if c == nil {
		return
	}
	if a.IsErrorResponse() {
		g.Error(c, a)
	} else {
		g.Success(c, a)
	}
}

// Generate lazy query parameters based on parameters and value
// Example : GenerateFuzzyQuery(GORM_DB_QUERY, []string{"name", "sex"})
func (g Gin) GenerateFuzzyQuery(c *gin.Context, tx *gorm.DB, fields []string) (*gorm.DB, error) {
	// If gin.Context is nil
	if c == nil {
		return nil, errors.New("gin.Context is nil")
	}
	// If fields is nil,no error will be reported and the original value will be returned
	if fields == nil {
		return tx, nil
	}

	var m = make(map[string]string)
	for _, v := range fields {
		m[v] = c.Query(v)
	}

	var gorm LibraryGorm
	return gorm.GenerateFuzzyQueries(tx, m)
}

// Get list with fuzzy query
// Example:
// var query = database.Database{}.Get().Model(&database.Supplier{}).Preload("Products").Order("created_at desc")
// dg := d.Gin{Context: c}
// var data []database.Supplier
// p, err := dg.GetListWithFuzzyQuery(query, nil, &data)
func (g Gin) GetListWithFuzzyQuery(c *gin.Context, query *gorm.DB, fuzzy_query_field_name []string, data_list_pointer interface{}) (p InterfacePagination, err error) {
	// If gin.Context is nil
	if c == nil {
		return p, errors.New("gin.Context is nil")
	}

	tx, err := g.GenerateFuzzyQuery(c, query, fuzzy_query_field_name)
	if err != nil {
		return p, err
	}

	var total int64
	tx.Count(&total)
	// Generate paginated data
	f := LibraryGorm{}.Paginate(c.Request)
	result := f(tx).Find(data_list_pointer)
	if result.Error != nil {
		return p, result.Error
	}

	page, err := strconv.Atoi(c.Query(FieldNamePaginationPage))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query(FieldNamePaginationPageSize))
	if err != nil {
		pageSize = 20
	}

	p = Pagination[InterfacePagination]{}.Get().Set(page, pageSize, int(total), nil)

	return p, nil
}

// API request interceptor in GIN, modify the returned fields
func (g Gin) ModifyApiFieldName() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get custom field map
		m := Config[InterfaceConfig]{}.Get().GetStringMap(ConfigPathApiField)
		if len(m) == 0 {
			c.Next()
			return
		}

		// Edit Query Keys
		g.editQueryKeys(c, m)

		if c.Request.Body == nil {
			c.Next()
			return
		}

		var bodyBytes []byte
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil || len(bodyBytes) == 0 {
			c.Next()
			return
		}

		var bodyMap map[string]interface{}
		err = json.Unmarshal(bodyBytes, &bodyMap)
		if err != nil {
			c.Next()
			return
		}

		g.recursionReplaceRequestKeys(bodyMap, m)

		modifiedBodyBytes, err := json.Marshal(bodyMap)
		if err != nil {
			c.Next()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(modifiedBodyBytes))
		c.Request.ContentLength = int64(len(modifiedBodyBytes))

		c.Next()
	}
}

func (g Gin) recursionReplaceRequestKeys(bodyMap, fieldMap map[string]interface{}) {
	for key, value := range bodyMap {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			g.recursionReplaceRequestKeys(nestedMap, fieldMap)
		} else {
			for k, v := range fieldMap {
				if v.(string) == key {
					bodyMap[k] = value
					delete(bodyMap, key)
				}
			}
		}
	}
}

// Edit Query Keys
func (g Gin) editQueryKeys(c *gin.Context, fieldMap map[string]interface{}) {
	for k, v := range fieldMap {
		if subValue, ok := c.GetQuery(v.(string)); ok {
			urlValue := c.Request.URL.Query()
			urlValue.Set(k, subValue)
			urlValue.Del(v.(string))
			c.Request.URL.RawQuery = urlValue.Encode()
		}
	}
}
