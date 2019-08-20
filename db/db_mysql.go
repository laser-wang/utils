package db

import (
	"database/sql"
	"fmt"
	"time"

	"caton/xw.wang/utils"

	clog "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

func OpenDB(url string, maxLT int, maxOC int, maxIC int) (DBmysql *sql.DB) {
	var err error

	DBmysql, err = sql.Open("mysql", url)
	if err != nil {
		clog.Critical(err.Error())
		panic(err)
	}
	utils.CheckErr(err, utils.CHECK_FLAG_EXIT)

	DBmysql.SetConnMaxLifetime(time.Duration(maxLT) * time.Second)
	DBmysql.SetMaxOpenConns(maxOC)
	DBmysql.SetMaxIdleConns(maxIC)

	clog.Info("[db opened]")

	//	clog.Infof("[db opened] url:'%s'", url)
	//	clog.Infof("[db opened] max_life_time:'%d'", maxLT)
	//	clog.Infof("[db opened] max_open_conns:'%d'", maxOC)
	//	clog.Infof("[db opened] max_idle_conns:'%d'", maxIC)

	return DBmysql

}

func CloseDB(DBmysql *sql.DB) {
	DBmysql.Close()
	clog.Info("[db closed] mysql")
}

//func DoQuery(DBmysql *sql.DB, sql string, args ...interface{}) (results map[int]map[int]string, err error) {

//	clog.Trace("[sql]: ", sql+" args:"+utils.Args2Str(args...))

//	rows, err := DBmysql.Query(sql, args...)
//	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
//	if err != nil {
//		return nil, err
//	}

//	cols, _ := rows.Columns()
//	values := make([][]byte, len(cols))
//	scans := make([]interface{}, len(cols))
//	for i := range values {
//		scans[i] = &values[i]
//	}
//	results = make(map[int]map[int]string)

//	i := 0
//	for rows.Next() {

//		err = rows.Scan(scans...)
//		utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
//		row := make(map[int]string) //每行数据
//		for k, v := range values {  //每行数据是放在values里面，现在把它挪到row里
//			row[k] = string(v)
//		}
//		results[i] = row //装入结果集中
//		i++
//	}
//	rows.Close()

//	return results, nil

//}

func DoQuery(DBmysql *sql.DB, sql string, args ...interface{}) (results [][]string, err error) {

	clog.Trace("[sql]: ", sql+" args:"+utils.Args2Str(args...))

	rows, err := DBmysql.Query(sql, args...)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return nil, err
	}

	cols, _ := rows.Columns()
	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	results = make([][]string, 0)

	i := 0
	for rows.Next() {

		err = rows.Scan(scans...)
		utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
		row := make([]string, 0)
		for _, v := range values { //每行数据是放在values里面，现在把它挪到row里
			row = append(row, string(v))
		}
		results = append(results, row) //装入结果集中
		i++
	}
	rows.Close()

	return results, nil

}

// DoTransQuery 支持事务的query
func DoTransQuery(DBmysql *sql.DB, tx *sql.Tx, sql string, args ...interface{}) (results [][]string, err error) {

	clog.Trace("[sql]: ", sql+" args:"+utils.Args2Str(args...))

	rows, err := tx.Query(sql, args...)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return nil, err
	}

	cols, _ := rows.Columns()
	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	results = make([][]string, 0)

	i := 0
	for rows.Next() {

		err = rows.Scan(scans...)
		utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
		row := make([]string, 0)
		for _, v := range values { //每行数据是放在values里面，现在把它挪到row里
			row = append(row, string(v))
		}
		results = append(results, row) //装入结果集中
		i++
	}
	rows.Close()

	return results, nil

}

func DoExec(DBmysql *sql.DB, sql string, args ...interface{}) (bool, error) {

	clog.Trace("[sql]: ", sql+" args:"+utils.Args2Str(args...))

	_, err := DBmysql.Exec(sql, args...)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err == nil {
		return true, err
	} else {
		return false, err
	}
}

// DoTransExec 支持事务
func DoTransExec(DBmysql *sql.DB, tx *sql.Tx, sql string, args ...interface{}) (bool, error) {

	clog.Trace("[sql]: ", sql+" args:"+utils.Args2Str(args...))

	_, err := tx.Exec(sql, args...)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err == nil {
		return true, err
	} else {
		return false, err
	}
}

// DoExecBatch 开启事务，执行批处理
func DoExecBatch(DBmysql *sql.DB, sqls []string, args [][]interface{}) (bool, error) {
	tx, errBegin := DBmysql.Begin()
	utils.CheckErr(errBegin, utils.CHECK_FLAG_LOGONLY)
	if errBegin != nil {
		return false, errBegin
	}
	var errExec error
	for idx, sql := range sqls {
		clog.Trace("[sql]: ", sql+" args:"+utils.Args2Str(args[idx]...))
		_, errExec = tx.Exec(sql, args[idx]...)
		utils.CheckErr(errExec, utils.CHECK_FLAG_LOGONLY)
		if errExec != nil {
			errRollback := tx.Rollback()
			utils.CheckErr(errRollback, utils.CHECK_FLAG_LOGONLY)
			clog.Error("[sql]:" + sql)
			return false, errExec
		}
	}
	errCommit := tx.Commit()
	utils.CheckErr(errCommit, utils.CHECK_FLAG_LOGONLY)
	if errCommit != nil {
		errRollback := tx.Rollback()
		utils.CheckErr(errRollback, utils.CHECK_FLAG_LOGONLY)
		return false, errCommit
	}

	return true, nil
}

// SqlColon 返回 ?,?,?......
func SqlColon(cnt int) string {
	if cnt == 0 {
		return ""
	}
	sql := ""
	for i := 0; i < cnt; i++ {
		sql = sql + `?,`
	}
	if utils.Substring(sql, len(sql)-1, len(sql)) == "," {
		sql = utils.Substring(sql, 0, len(sql)-1)
	}
	fmt.Println("SqlColon:" + sql)
	return sql
}
