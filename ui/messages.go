package ui

import "gorm.io/gorm"

type errorMsg error

type dbMsg struct {
	DB *gorm.DB
}
