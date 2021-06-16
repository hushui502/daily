// https://github.com/Shelnutt2/db2struct
package main

import (
	db2struct "da2struct"
	"fmt"
	"os"
	"strconv"

	goopt "github.com/droundy/goopt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/howeyc/gopass"
)

var (
	mariadbHost       = os.Getenv("MYSQL_HOST")
	mariadbHostPassed = goopt.String([]string{"-H", "--host"}, "", "Host to check mariadb status of")
	mariadbPort       = goopt.Int([]string{"--mysql_port"}, 3306, "Specify a port to connect to")
	mariadbTable      = goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")
	mariadbDatabase   = goopt.String([]string{"-d", "--database"}, "nil", "Database to for connection")
	mariadbPassword   *string
	mariadbUser       = goopt.String([]string{"-u", "--user"}, "user", "user to connect to database")
	verbose           = goopt.Flag([]string{"-v", "--verbose"}, []string{}, "Enable verbose output", "")
	packageName       = goopt.String([]string{"--package"}, "", "name to set for package")
	structName        = goopt.String([]string{"--struct"}, "", "name to set for struct")
	jsonAnnotation    = goopt.Flag([]string{"--json"}, []string{"--no-json"}, "Add json annotations (default)", "Disable json annotations")
	gormAnnotation    = goopt.Flag([]string{"--gorm"}, []string{}, "Add gorm annotations (tags)", "")
	gureguTypes       = goopt.Flag([]string{"--guregu"}, []string{}, "Add guregu null types", "")
	targetFile        = goopt.String([]string{"--target"}, "", "Save file path")
)

func init() {
	goopt.OptArg([]string{"-p", "--password"}, "", "Mysql password", getMariadbPassword)
	// Setup goopts
	goopt.Description = func() string {
		return "Mariadb http Check"
	}
	goopt.Version = "0.0.2"
	goopt.Summary = "db2struct [-H] [-p] [-v] --package pkgName --struct structName --database databaseName --table tableName"

	//Parse options
	goopt.Parse(nil)
}

func main() {
	// Username is required
	if mariadbUser == nil || *mariadbUser == "user" {
		fmt.Println("Username is required! Add it with --user=name")
		return
	}

	// If a mariadb host is passed use it
	if mariadbHostPassed != nil && *mariadbHostPassed != "" {
		mariadbHost = *mariadbHostPassed
	}

	if mariadbPassword != nil && *mariadbPassword == "" {
		fmt.Print("Password: ")
		password, err := gopass.GetPasswd()
		stringPassword := string(password)
		mariadbPassword = &stringPassword
		if err != nil {
			fmt.Println("Error reading password: " + err.Error())
			return
		}
	} else if mariadbPassword == nil {
		p := ""
		mariadbPassword = &p
	}

	if *verbose {
		fmt.Println("Connecting to mysql server " + mariadbHost + ":" + strconv.Itoa(*mariadbPort))
	}

	if mariadbDatabase == nil || *mariadbDatabase == "" {
		fmt.Println("Database can not be null")
		return
	}

	if mariadbTable == nil || *mariadbTable == "" {
		fmt.Println("Table can not be null")
		return
	}

	columnDataTypes, columnsSorted, err := db2struct.GetColumnsFromMysqlTable(*mariadbUser, *mariadbPassword, mariadbHost, *mariadbPort, *mariadbDatabase, *mariadbTable)
	if err != nil {
		fmt.Println("Error in selecting column data information from mysql information schema")
		return
	}

	// If structName is not set we need to default it
	if structName == nil || *structName == "" {
		*structName = "newstruct"
	}
	// If packageName is not set we need to default it
	if packageName == nil || *packageName == "" {
		*packageName = "newpackage"
	}
	// Generate struct string based on columnDataTypes
	structure, err := db2struct.Generate(*columnDataTypes, columnsSorted, *mariadbTable, *structName, *packageName, *jsonAnnotation, *gormAnnotation, *gureguTypes)
	if err != nil {
		fmt.Println("Error in creating struct from json: " + err.Error())
		return
	}

	if targetFile != nil && *targetFile != "" {
		file, err := os.OpenFile(*targetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Open File fail: " + err.Error())
			return
		}
		length, err := file.WriteString(string(structure))
		if err != nil {
			fmt.Println("Save File fail: " + err.Error())
			return
		}
		fmt.Printf("wrote %d bytes\n", length)
	} else {
		fmt.Println(string(structure))
	}
}

func getMariadbPassword(password string) error {
	mariadbPassword = new(string)
	*mariadbPassword = password
	return nil
}