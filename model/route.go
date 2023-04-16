package model

import "gorm.io/gorm"

type Route struct {
	gorm.Model
	Name          string
	Description   string
	Method        string
	RequestPath   string
	TargetService string
	TargetPath    string

	// middleware
	NeedAuth bool
}
