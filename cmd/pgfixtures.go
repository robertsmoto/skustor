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

        fixturesFile, err := os.ReadFile(
            path.Join(os.Getenv("ROOTDR"), "cmd/data/fixtures.json"))
        if err != nil {
            log.Printf("pgfixturesCmd %s", err)
        }

        //instantiate the model structs
        //user := models.User{}
        //users := models.Users{}
        //place := models.Place{}
        //places := models.Places{}
        //priceClass := models.PriceClass{}
        //priceClasses := models.PriceClasses{}
        //unit := models.Unit{}
        //units := models.Units{}
        collectionNodes := models.CollectionNodes{}
        contentNodes := models.ContentNodes{}
        //item := models.Item{}
        //items := models.Items{}

        // loads and validates the nodes (singular versions of structs)
        loaderNodes := []models.LoaderProcesserUpserter {
            //&user,
            //&users,
            //&place,
            //&places,
            //&priceClass,
            //&priceClasses,
            //&unit,
            //&units,
            &collectionNodes,
            &contentNodes,
            //&item,
            //&items,
        }


        // open the db
        devPostgres := postgres.PostgresDb{}
        pgDb, err := postgres.Open(&devPostgres)
        userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

        for _, node := range loaderNodes {
            err = models.JsonLoaderUpserterHandler(
                node, userId, &fixturesFile, pgDb)
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
