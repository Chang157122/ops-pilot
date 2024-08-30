package models

import "opsPilot/internal/pkg/log"

func CreateTable() {
	var models []struct{ LoginModel }
	for _, v := range models {
		if !db.HasTable(&v) {
			table := db.CreateTable(&v)
			if table.Error != nil {
				log.Logger.Errorf("create table failed err: %v", table.Error)
			}
		}
	}
}
