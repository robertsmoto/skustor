/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
    "log"
    "os"
    "path"

    //"github.com/robertsmoto/skustor/internal/configs"
    "github.com/robertsmoto/skustor/internal/models"
    "github.com/robertsmoto/skustor/internal/postgres"
	"github.com/spf13/cobra"
)

// pgbuilderCmd represents the pgbuilder command
var pgfixturesCmd = &cobra.Command{
	Use:   "pgfixtures",
	Short: "Builds the postgres fixtures.",
	Long: `
        Builds the postgres fixtures for groups, items, people, orders.
        `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("building fixtures ...")


        //conf := configs.Config{}
        //err := configs.Load(&conf)

        aid := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"
        fixturesFile, err := os.ReadFile(
            path.Join(os.Getenv("ROOTDR"), "cmd/data/fixtures.json"))
        if err != nil {
            log.Println("pgfixturesCmd %s", err)
        }

        // open the db
        devPostgres := postgres.PostgresDb{}
        pgDb, err := postgres.Open(&devPostgres)



        // instantiate the structs
        collectionNodes := models.CollectionNodes{}
        itemNodes := models.ItemNodes{}
        contentNodes := models.ContentNodes{}

        // loader validator
        lvNodes := []models.LoaderValidator{
            &collectionNodes,
            &itemNodes,
            &contentNodes,
        }
        for _, node := range lvNodes {
            err = models.LoadValidateHandler(node, &fixturesFile)
            if err != nil {
                log.Println("pgfixtures 01", err)
            }
        }

        // upsert
        upsertNodes := []models.Upserter{
            &collectionNodes,
            &itemNodes,
            &contentNodes,
        }
        for _, node := range upsertNodes {
            err = models.UpsertHandler(node, aid, pgDb)
            if err != nil {
                log.Println("pgfixtures 02", err)
            }
        }

        // foreign key update
        fkNodes := []models.ForeignKeyUpdater{
            &collectionNodes,
            &itemNodes,
            &contentNodes,
        }
        for _, node := range fkNodes {
            err = models.ForeignKeyUpdateHandler(node, pgDb)
            if err != nil {
                log.Println("pgfixtures 03", err)
            }
        }

        // related table upsert
        rtNodes := []models.RelatedTableUpserter{
            &collectionNodes,
            &itemNodes,
            &contentNodes,
        }
        for _, node := range rtNodes {
            err = models.RelatedTableUpsertHandler(node, aid, pgDb)
            if err != nil {
                log.Println("pgfixtures 04", err)
            }
        }



        pgDb.Close()
	},
}

func init() {
	rootCmd.AddCommand(pgfixturesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pgbuilderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pgbuilderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
