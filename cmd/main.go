package main

import (
	"fmt"
	"os"

	"github.com/simplesvet/dbexport"
	"github.com/spf13/cobra"
)

func exportDbObjectsCmd(objectType string, objectName string) {
	fmt.Println("Exporting", objectType)

	var dbObjects []dbexport.DbObject

	if objectType == dbexport.PROCEDURES {
		dbObjects = dbexport.GetProcedures(objectName)
	} else if objectType == dbexport.FUNCTIONS {
		dbObjects = dbexport.GetFunctions(objectName)
	} else if objectType == dbexport.TRIGGERS {
		dbObjects = dbexport.GetTriggers(objectName)
	} else if objectType == dbexport.EVENTS {
		dbObjects = dbexport.GetEvents(objectName)
	} else if objectType == dbexport.VIEWS {
		dbObjects = dbexport.GetViews(objectName)
	} else if objectType == dbexport.TABLES {
		dbObjects = dbexport.GetTables(objectName)
	} else {
		dbObjects = dbexport.GetAll()
	}

	savedFiles := dbexport.SaveDbObjects(dbObjects)

	for _, file := range savedFiles {
		fmt.Println("File saved in", file)
	}
}

var objName string

var rootCmd = &cobra.Command{
	Use:   "dbexport all | dbexport [object_type] --name [object_name]",
	Short: "DBExport is a fast tool to sync databases objects with the file system",
	Long:  `DBExport is a fast tool to sync databases objects with the file system`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		firstArg := args[0]

		if firstArg == "help" {
			cmd.Help()
		} else if firstArg == "all" {
			exportDbObjectsCmd("", "")
		} else {
			exportDbObjectsCmd(firstArg, objName)
		}
	},
}

func main() {
	rootCmd.Flags().StringVarP(&objName, "name", "n", "", "nome do objeto no banco")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
