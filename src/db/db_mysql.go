package db

import (
	"bytes"
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//DbSource 全局mysql访问对象
var DbSource *sql.DB

func init() {
	var err error
	DbSource, err = sql.Open("mysql", "truck:jjb910123@tcp(120.26.10.20:3306)/truck30?parseTime=true&charset=utf8")
	if err != nil {
		log.Fatal("Database connect error: " + err.Error())
	}
	DbSource.SetMaxOpenConns(10)
	DbSource.SetMaxIdleConns(5)
	DbSource.Ping()
	log.Println("Mysql DB connected")
}

//BuildInsert SQL Insert语句生成, eg. INSERT INTO xxx (`a`) VALUES (1), (2)
func BuildInsert(table string, cols []string, len int) string {
	buf := bytes.Buffer{}
	buf.WriteString("INSERT INTO ")
	buf.WriteString(QuoteIdent(table))

	placeholderBuf := new(bytes.Buffer)
	placeholderBuf.WriteString("(")
	buf.WriteString(" (")
	for i, col := range cols {
		if i > 0 {
			buf.WriteString(",")
			placeholderBuf.WriteString(",")
		}
		buf.WriteString(QuoteIdent(col))
		placeholderBuf.WriteString(Placeholder())
	}
	buf.WriteString(") VALUES ")
	placeholderBuf.WriteString(")")
	buf.WriteString(placeholderBuf.String())

	placeholderStr := placeholderBuf.String()
	for i := 1; i < len; i++ {
		buf.WriteString(",")
		buf.WriteString(placeholderStr)
	}
	return buf.String()
}

//SaveOne 单行数据插入
func SaveOne(sql string, args []interface{}) (int64, error) {
	stmt, _ := DbSource.Prepare(sql)
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// TxExec 事务执行方法
func TxExec(f func(tx *sql.Tx) error) (err error) {
	//开启事务
	tx, err := DbSource.Begin()
	if err != nil {
		return err
	}

	err = f(tx)

	defer tx.Rollback()
	tx.Commit()
	return
}

// SaveTx 带事务的数据插入
func SaveTx(tx *sql.Tx, sql string, args []interface{}) (id int64, err error) {
	_, id, err = execTx(tx, sql, args...)
	return
}

//Update 单行数据更新
func Update(sql string, args ...interface{}) (int64, error) {
	stmt, err := DbSource.Prepare(sql)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return num, nil
}

//UpdateTx 单行数据更新
func UpdateTx(tx *sql.Tx, sql string, args ...interface{}) (eff int64, err error) {
	eff, _, err = execTx(tx, sql, args...)
	return
}

func execTx(tx *sql.Tx, sql string, args ...interface{}) (int64, int64, error) {
	stmt, err := tx.Prepare(sql)
	// defer stmt.Close()
	if err != nil {
		return 0, 0, err
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, 0, err
	}

	eff, err := res.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	return eff, lastID, nil
}

//GetOne 获取单行数据，返回所有数据类型为string
func GetOne(rows *sql.Rows) []string {
	if rows == nil {
		return nil
	}
	cols, _ := rows.Columns()
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))
	dest := make([]interface{}, len(cols))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}
	if rows.Next() {
		err := rows.Scan(dest...)
		if err == nil {
			for i, raw := range rawResult {
				if raw == nil {
					result[i] = ""
				} else {
					result[i] = string(raw)
				}
			}
		}
	} else {
		return nil
	}
	return result
}

//GetOneMap 获取一行返回map
func GetOneMap(rows *sql.Rows) map[string]string {
	if rows == nil {
		return nil
	}
	cols, _ := rows.Columns()
	rawResult := make([][]byte, len(cols))
	result := make(map[string]string, len(cols))
	dest := make([]interface{}, len(cols))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}
	if rows.Next() {
		err := rows.Scan(dest...)
		if err == nil {
			for i, raw := range rawResult {
				if raw == nil {
					result[cols[i]] = ""
				} else {
					result[cols[i]] = string(raw)
				}
			}
		}
	} else {
		return nil
	}
	return result
}

//QuoteIdent sql语句字段引号 eg.`a`
func QuoteIdent(s string) string {
	return quoteIdent(s, "`")
}

//Placeholder sql语句预处理值占位符
func Placeholder() string {
	return "?"
}

func quoteIdent(s, quote string) string {
	part := strings.SplitN(s, ".", 2)
	if len(part) == 2 {
		return quoteIdent(part[0], quote) + "." + quoteIdent(part[1], quote)
	}
	return quote + s + quote
}

// Scanner 获取字段方法
type Scanner func(rs *sql.Rows) error

// QueryMore 查询多行
func QueryMore(query string, f Scanner, args ...interface{}) error {
	rs, err := DbSource.Query(query, args...)
	if err != nil {
		return err
	}
	defer rs.Close()

	return f(rs)
}
