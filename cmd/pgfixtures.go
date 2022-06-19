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

    "github.com/robertsmoto/skustor/models"
    "github.com/robertsmoto/skustor/configs"
    "github.com/robertsmoto/skustor/tools"
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


        conf := configs.Config{}
        err := configs.Load(&conf)

        fixturesFile, err := os.ReadFile(path.Join(conf.RootDir, "cmd/data/fixtures.json"))
        if err != nil {
            log.Print("pgfixturesCmd ", err)
        }

        // instantiate the model structs
        users := models.SvUserNodes{}
        locations := models.LocationNodes{}
        priceClasses := models.PriceClassNodes{}
        units := models.UnitNodes{}
        clusters := models.ClusterNodes{}
        items := models.ItemNodes{}

        // load the model structs
        models.LoaderHandler(&users, fixturesFile)
        models.LoaderHandler(&locations, fixturesFile)
        models.LoaderHandler(&priceClasses, fixturesFile)
        models.LoaderHandler(&units, fixturesFile)
        models.LoaderHandler(&clusters, fixturesFile)
        models.LoaderHandler(&items, fixturesFile)

        // open the db
        devPostgres := tools.PostgresDev{}
        devDb, err := tools.Open(&devPostgres)

        // upsert the structs
        userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

        // create only top-level data
        for _, user := range users.Nodes {
            models.UpsertHandler(&user, devDb)
        }
        for _, location := range locations.Nodes {
            location.SvUserId = userId
            models.UpsertHandler(&location, devDb)
            // uploading addresses
            for _, addressNode := range location.AddressNodes.Nodes {
                addressNode.SvUserId = location.SvUserId // make sure to add ids
                addressNode.LocationId = location.Id    // make sure to add ids
                models.UpsertHandler(&addressNode, devDb)
            }
        }
        for _, pc := range priceClasses.Nodes {
            pc.SvUserId = userId
            models.UpsertHandler(&pc, devDb)
        }
        for _, unit := range units.Nodes {
            unit.SvUserId = userId
            models.UpsertHandler(&unit, devDb)
        }

        for _, cluster := range clusters.Nodes {
            cluster.SvUserId = userId
            models.UpsertHandler(&cluster, devDb)

            // uploading images
            for _, imgNode := range cluster.ImageNodes.Nodes {
                imgNode.SvUserId = cluster.SvUserId // make sure to add ids
                imgNode.ClusterId = cluster.Id    // make sure to add ids
                models.ImgHandler(&imgNode, devDb)
            }
        }

        for _, item := range items.Nodes {
            item.SvUserId = userId
            models.UpsertHandler(&item, devDb)

            // uploading images
            for _, imgNode := range item.ImageNodes.Nodes {
                imgNode.SvUserId = item.SvUserId // make sure to add ids
                imgNode.ItemId = item.Id    // make sure to add ids
                models.ImgHandler(&imgNode, devDb)
            }
        }
        devDb.Close()
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
