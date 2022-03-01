package main

import (
	"errors"
	"fmt"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"strings"
	"time"
)

func doMake(arg2, arg3 string) error {
	switch arg2 {
	case "migration":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the migration a name"))
		}
		dbType := cel.DB.DatabaseType
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)

		upFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		err := copyFileFromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		if err != nil {
			exitGracefully(err)
		}
		err = copyFileFromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			exitGracefully(err)
		}
	case "auth":
		return doAuth()
	case "handler":
		if arg3 == "" {
			return errors.New("you must give the handler a name")
		}
		filename := cel.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(filename) {
			return errors.New(filename + " already exists!")
		}

		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			return err
		}
		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

		err = ioutil.WriteFile(filename, []byte(handler), 0644)
		if err != nil {
			return err
		}
	case "model":
		if arg3 == "" {
			return errors.New("you must give the model a name")
		}
		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			return err
		}
		model := string(data)
		plur := pluralize.NewClient()
		var modelName = arg3
		var tableName = arg3

		if plur.IsPlural(arg3) {
			modelName = plur.Singular(arg3)
			tableName = strings.ToLower(tableName)
		} else {
			tableName = strings.ToLower(plur.Plural(arg3))
		}
		fileName := cel.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
		if fileExists(fileName) {
			return errors.New(fileName + " already exists!")
		}
		model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)
		err = copyDataToFile([]byte(model), fileName)
		if err != nil {
			return err
		}
	case "session":
		return doSessionTable()
	}

	return nil
}
