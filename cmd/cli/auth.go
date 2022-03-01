package main

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

func doAuth() error {
	dbType := cel.DB.DatabaseType
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := cel.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := cel.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".sql", upFile)
	if err != nil {
		return err
	}

	err = copyDataToFile([]byte("drop table if exists users cascade; drop table if exists tokens cascade; drop table if exists remember_tokens;"), downFile)
	if err != nil {
		return err
	}

	err = doMigrate("up", "")
	if err != nil {
		return err
	}

	err = copyFileFromTemplate("templates/data/user.go.txt", cel.RootPath+"/data/user.go")
	if err != nil {
		return err
	}

	err = copyFileFromTemplate("templates/data/token.go.txt", cel.RootPath+"/data/token.go")
	if err != nil {
		return err
	}

	err = copyFileFromTemplate("templates/middleware/auth.go.txt", cel.RootPath+"/middleware/auth.go")
	if err != nil {
		return err
	}

	err = copyFileFromTemplate("templates/middleware/auth-token.go.txt", cel.RootPath+"/middleware/auth-token.go")
	if err != nil {
		return err
	}

	color.Yellow("    - users, tokens and remember_tokens migrations created and executed")
	color.Yellow("    - users and tokens models created")
	color.Yellow("    - auth middleware created")
	color.Yellow("")
	color.Yellow("    - Don't forget to add user and token models in data/models.go and to add appropiate middleware to your routes!")

	return nil
}
