package dao

import (
	"database/sql"
	"github.com/pkg/errors"
)

func QueryUser(id int) (interface{}, error) {
	r, err := NoRows()
	//如果没有查询记录，就把error记录到wrap中
	if err == sql.ErrNoRows {
		return nil, errors.Wrapf(errors.New("该条结果未找到"), "user id: %d", id)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "查询异常,id为：%d", id)
	}
	return r, nil
}

func NoRows() (interface{}, error) {
	return nil, sql.ErrNoRows
}
