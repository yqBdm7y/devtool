package d

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
