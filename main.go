package main

import (
	"github.com/dawsonliu/godbr/core"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func main() {
	//LoadApis()

	// core.LoadCsis()
	core.Start()
}
