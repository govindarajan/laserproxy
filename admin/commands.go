package admin

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/govindarajan/laserproxy/store"
	"github.com/govindarajan/laserproxy/worker"
)

func handleSaveCmd(qs []string, db *sql.DB) error {
	if len(qs) < 2 {
		return fmt.Errorf("Unknown command. You mean 'save config'?")
	}
	switch qs[1] {
	case "config":
		return store.MoveConfigToDisk(db)
	default:
		return fmt.Errorf("Unknown command. 'save config' alone supported as of now")
	}
	return nil
}

func handleFrontendCmd(qs []string, db *sql.DB) error {
	if len(qs) < 2 || strings.ToLower(qs[1]) != "reload" {
		return fmt.Errorf("Not supported now")
	}
	worker.RefreshFrontends(db)
	return nil
}

func handleLoadCmd(qs []string, db *sql.DB) error {
	if len(qs) < 2 {
		return fmt.Errorf("Unknown command. You mean 'load config'?")
	}
	switch qs[1] {
	case "config":
		store.LoadConfigFromDisk(db)
		worker.RefreshFrontends(db)
	default:
		return fmt.Errorf("Unknown command. 'save config' alone supported as of now")
	}
	return nil
}
