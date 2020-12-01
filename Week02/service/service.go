package service

import "Week02/dao"

func HandleQuery(id int) (interface{}, error) {
	return dao.QueryUser(id)
}
