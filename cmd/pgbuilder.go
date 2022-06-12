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

    "github.com/robertsmoto/skustor/tools"
	"github.com/spf13/cobra"
)

// pgbuilderCmd represents the pgbuilder command
var pgbuilderCmd = &cobra.Command{
	Use:   "pgbuilder",
	Short: "Builds the postgres databases.",
	Long: `
        Builds the postgres databases, tables, columns,
        and constraints if they don't exist.
        `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("starting ...")

        qstr := `
        -- groups
        CREATE TABLE IF NOT EXISTS groups (
            id UUID PRIMARY KEY);
        ALTER TABLE groups ADD COLUMN IF NOT EXISTS user_id UUID;
        ALTER TABLE groups ADD COLUMN IF NOT EXISTS parent_id UUID;
        ALTER TABLE groups ADD COLUMN IF NOT EXISTS type VARCHAR (200);
        ALTER TABLE groups ADD COLUMN IF NOT EXISTS name VARCHAR (200);
        ALTER TABLE groups ADD COLUMN IF NOT EXISTS description VARCHAR (200);
        ALTER TABLE groups ADD COLUMN IF NOT EXISTS keywords VARCHAR (200);
        ALTER TABLE groups ADD COLUMN IF NOT EXISTS link_url VARCHAR (200);
        ALTER TABLE groups ADD COLUMN IF NOT EXISTS link_text VARCHAR (200);

        -- images
        CREATE TABLE IF NOT EXISTS images (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid());
        ALTER TABLE images ADD COLUMN IF NOT EXISTS user_id UUID;
        ALTER TABLE images ADD COLUMN IF NOT EXISTS group_id UUID;
        ALTER TABLE images ADD COLUMN IF NOT EXISTS item_id UUID;
        ALTER TABLE images ADD COLUMN IF NOT EXISTS url VARCHAR (200);
        ALTER TABLE images ADD COLUMN IF NOT EXISTS height VARCHAR (200);
        ALTER TABLE images ADD COLUMN IF NOT EXISTS width VARCHAR (200);
        ALTER TABLE images ADD COLUMN IF NOT EXISTS title VARCHAR (200);
        ALTER TABLE images ADD COLUMN IF NOT EXISTS alt VARCHAR (200);
        ALTER TABLE images ADD COLUMN IF NOT EXISTS caption VARCHAR (200);
        ALTER TABLE images ADD COLUMN IF NOT EXISTS position INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE images ADD COLUMN IF NOT EXISTS featured INTEGER NOT NULL DEFAULT 0;

        --items
        CREATE TABLE IF NOT EXISTS items (
            id UUID PRIMARY KEY);
        ALTER TABLE items ADD COLUMN IF NOT EXISTS user_id UUID;
        ALTER TABLE items ADD COLUMN IF NOT EXISTS parent_id UUID;
        ALTER TABLE items ADD COLUMN IF NOT EXISTS unit_id UUID;
        ALTER TABLE items ADD COLUMN IF NOT EXISTS type VARCHAR (50);
        ALTER TABLE items ADD COLUMN IF NOT EXISTS sku VARCHAR (50);
        ALTER TABLE items ADD COLUMN IF NOT EXISTS name VARCHAR (200);
        ALTER TABLE items ADD COLUMN IF NOT EXISTS description VARCHAR (200);
        ALTER TABLE items ADD COLUMN IF NOT EXISTS keywords VARCHAR (200);
        ALTER TABLE items ADD COLUMN IF NOT EXISTS cost BIGINT NOT NULL DEFAULT 0;
        ALTER TABLE items ADD COLUMN IF NOT EXISTS cost_override BIGINT NOT NULL DEFAULT 0;
        -- for parts
        ALTER TABLE items ADD COLUMN IF NOT EXISTS price_class_id UUID;
        ALTER TABLE items ADD COLUMN IF NOT EXISTS price BIGINT NOT NULL DEFAULT 0;
        ALTER TABLE items ADD COLUMN IF NOT EXISTS price_override BIGINT NOT NULL DEFAULT 0;

        -- join_group_item
        CREATE TABLE IF NOT EXISTS join_group_item (
            id UUID PRIMARY KEY);
        ALTER TABLE join_group_item ADD COLUMN IF NOT EXISTS user_id UUID;
        ALTER TABLE join_group_item ADD COLUMN IF NOT EXISTS group_id UUID;
        ALTER TABLE join_group_item ADD COLUMN IF NOT EXISTS item_id UUID;

        -- relationships: images? groups
        -- ALTER TABLE items ADD COLUMN IF NOT EXISTS identifiers_id UUID;
        -- ALTER TABLE items ADD COLUMN IF NOT EXISTS measurements_id UUID;
        -- ALTER TABLE items ADD COLUMN IF NOT EXISTS digital_assetts_id UUID;
        -- ALTER TABLE items ADD COLUMN IF NOT EXISTS images_id UUID;

        -- table constraints
        `

        fmt.Println(qstr)


        // open each db
        devPostgres := tools.PostgresDev{}
        devConn, err := tools.Open(&devPostgres)
        if err != nil {
            log.Print("Connection error", err)
        }
        _, err = devConn.Exec(qstr)
        if err != nil {
            log.Print("Error creating or updating database.", err)
        }


        fmt.Println(devConn)

        devConn.Close()
        
        // run the query on each db

        // close the db


	},
}

func init() {
	rootCmd.AddCommand(pgbuilderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pgbuilderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pgbuilderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
