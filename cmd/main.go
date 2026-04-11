package main

import (
	"github.com/ebube-ugwu/snape/internal/data"
)

func main() {
	data.MigrateDB(data.DbURI)

	// db := data.OpenDB(data.DbURI)

}
