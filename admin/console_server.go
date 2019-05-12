package admin

import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/govindarajan/laserproxy/logger"
	"github.com/govindarajan/laserproxy/store"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
)

func StartAdminServer() {
	l, _ := net.Listen("tcp", "127.0.0.1:4000")
	// TODO: Take it from config
	adminUser := "someuser"
	adminPwd := "1"
	for {
		// Wait for new connection
		c, _ := l.Accept()

		// Process the connection in thread.
		go func() {
			conn, err := server.NewConn(c, adminUser, adminPwd, &SQLiteHandler{})
			if err != nil || conn == nil {
				// Wrong pwd??
				logger.LogWarn(err.Error())
				return
			}
			for {
				err := conn.HandleCommand()
				if err != nil {
					logger.LogInfo(err.Error())
					break
				}
			}
		}()
	}
}

type SQLiteHandler struct {
	db *sql.DB
}

func (h *SQLiteHandler) UseDB(dbName string) (e error) {
	// TODO : Change the DB according to db name
	h.db, e = store.GetConnection()
	return e
}

func (h *SQLiteHandler) HandleQuery(query string) (*mysql.Result, error) {
	if h.db == nil {
		return nil, fmt.Errorf("no database selected")
	}
	query = strings.TrimRight(query, ";")
	ss := strings.Split(query, " ")

	res := &mysql.Result{}
	switch strings.ToLower(ss[0]) {
	case "select":
		rows, err := h.db.Query(query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		cl, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		var sr [][]interface{}
		for rows.Next() {
			rowVals := make([]interface{}, len(cl))
			rowValPtrs := make([]interface{}, len(cl))

			for i, _ := range cl {
				rowValPtrs[i] = &rowVals[i]
			}
			if err = rows.Scan(rowValPtrs...); err != nil {
				return nil, err
			}
			sr = append(sr, rowVals)
		}

		res.Resultset, err = mysql.BuildSimpleResultset(cl, sr, false)
		if err != nil {
			return nil, err
		}

	case "drop":
		return nil, fmt.Errorf("not supported now")

	case "frontend":
		if err := handleFrontendCmd(ss, h.db); err != nil {
			return nil, err
		}

	case "save":
		if err := handleSaveCmd(ss, h.db); err != nil {
			return nil, err
		}

	case "load":
		if err := handleLoadCmd(ss, h.db); err != nil {
			return nil, err
		}

	default:
		r, err := h.db.Exec(query)
		if err != nil {
			return nil, err
		}
		ra, err := r.RowsAffected()
		if err != nil {
			return nil, err
		}
		res.AffectedRows = uint64(ra)
	}

	return res, nil
}

func (h *SQLiteHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	return nil, fmt.Errorf("not supported now")
}

func (h *SQLiteHandler) HandleStmtPrepare(query string) (params int, columns int, context interface{}, err error) {

	return 1, 1, nil, nil
}

func (h *SQLiteHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	return nil, nil
}

func (h *SQLiteHandler) HandleStmtClose(context interface{}) error {
	return fmt.Errorf("not supported now")
}

func (h *SQLiteHandler) HandleOtherCommand(cmd byte, data []byte) error {
	return fmt.Errorf("not supported now")
}
