package d

import "errors"

type MySQL struct{}

// Generate lazy query parameters based on parameters and value
// Example : GenerateFuzzyQueries(map[string]string{"name": "John", "sex": "female"})
func (m MySQL) GenerateFuzzyQueries(fields map[string]string) (whereClause string, args []interface{}, err error) {
	// If map is nil
	if fields == nil {
		return "", nil, errors.New("map is nil")
	}

	for k, v := range fields {
		// If there is no value, skip the current field
		if len(v) == 0 {
			continue
		}
		// If whereClause is empty, it is the first traversal, otherwise just addthe OR clause
		if len(whereClause) == 0 {
			whereClause = k + " LIKE ?"
		} else {
			whereClause += " AND " + k + " LIKE ?"
		}
		args = append(args, "%"+v+"%")
	}

	return whereClause, args, nil
}
