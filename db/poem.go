package db

import (
	"encoding/json"
	"fmt"
)

type Poem struct {
	Id      int
	Title   string
	Author  string
	Dynasty string
	Content string
}

func (p *Poem) Insert() bool {
	stmtInsert, err := db.Prepare(" insert into poem (title,author,dynasty,content) values (?,?,?,?) ")
	if checkError(err) {
		return false
	}

	_, err = stmtInsert.Exec(p.Title, p.Author, p.Dynasty, p.Content)
	if checkError(err) {
		return false
	}
	return true
}

func (p *Poem) Save() {
	data, _ := json.Marshal(p)
	fmt.Println(string(data))

	res := p.Insert()
	fmt.Println("添加结果: ", res)
}

//根据条件获取诗
func QueryPoems(field string, value string) (poems []Poem, err error) {
	sqlStr := "select id,title,author,dynasty,content from poem where 1=1"
	sqlStr += fmt.Sprintf(" and %s = ?", field)
	fmt.Println(sqlStr, value)

	//预处理
	stmtOut, err := db.Prepare(sqlStr)
	if checkError(err) {
		return nil, err
	}

	//查询记录
	rows, err := stmtOut.Query(value)
	if checkError(err) {
		return nil, err
	}
	rows.Next()
	{
		p := Poem{}
		err = rows.Scan(&p.Id, &p.Title, &p.Author, &p.Dynasty, &p.Content)
		if checkError(err) {
			return nil, err
		}

		poems = append(poems, p)
	}
	return poems, nil
}
