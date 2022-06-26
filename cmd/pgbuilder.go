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
    "github.com/robertsmoto/skustor/configs"
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
        -- ##################################
        --  table: sv_user
        -- ##################################
        CREATE TABLE IF NOT EXISTS sv_user (
            id UUID PRIMARY KEY);
        ALTER TABLE sv_user ADD COLUMN IF NOT EXISTS auth UUID;
        ALTER TABLE sv_user ADD COLUMN IF NOT EXISTS key UUID;
        ALTER TABLE sv_user ADD COLUMN IF NOT EXISTS username VARCHAR (100);
        ALTER TABLE sv_user ADD COLUMN IF NOT EXISTS firstname VARCHAR (100);
        ALTER TABLE sv_user ADD COLUMN IF NOT EXISTS lastname VARCHAR (100);
        ALTER TABLE sv_user ADD COLUMN IF NOT EXISTS nickname VARCHAR (100);

        -- ##################################
        --  table: unit
        -- ##################################
        CREATE TABLE IF NOT EXISTS unit (
            id UUID PRIMARY KEY);
        ALTER TABLE unit
            ADD COLUMN IF NOT EXISTS sv_user_id UUID
            REFERENCES sv_user(id);
        ALTER TABLE unit ADD COLUMN IF NOT EXISTS singular VARCHAR (100);
        ALTER TABLE unit ADD COLUMN IF NOT EXISTS singular_display VARCHAR (100);
        ALTER TABLE unit ADD COLUMN IF NOT EXISTS plural VARCHAR (100);
        ALTER TABLE unit ADD COLUMN IF NOT EXISTS plural_display VARCHAR (100);

        -- ##################################
        --  table: price_class
        -- ##################################
        CREATE TABLE IF NOT EXISTS price_class (
            id UUID PRIMARY KEY);
        ALTER TABLE price_class
            ADD COLUMN IF NOT EXISTS sv_user_id UUID
            REFERENCES sv_user(id);
        ALTER TABLE price_class ADD COLUMN IF NOT EXISTS type VARCHAR (100);
        ALTER TABLE price_class ADD COLUMN IF NOT EXISTS name VARCHAR (100);
        ALTER TABLE price_class ADD COLUMN IF NOT EXISTS amount
            BIGINT NOT NULL DEFAULT 0;
        ALTER TABLE price_class ADD COLUMN IF NOT EXISTS note VARCHAR (100);

        -- ##################################
        --  table: aaa_collection
        -- ##################################
        CREATE TABLE IF NOT EXISTS aaa_collection (
            id UUID PRIMARY KEY);
        ALTER TABLE aaa_collection
            ADD COLUMN IF NOT EXISTS sv_user_id UUID
            REFERENCES sv_user(id);
        ALTER TABLE aaa_collection
            ADD COLUMN IF NOT EXISTS parent_id UUID
            REFERENCES aaa_collection(id);
        ALTER TABLE aaa_collection ADD COLUMN IF NOT EXISTS document JSONB;
        -- indexes

        CREATE INDEX IF NOT EXISTS aaa_collection_id_idx ON aaa_collection (id);
        CREATE INDEX IF NOT EXISTS aaa_collection_sv_user_id_idx ON aaa_collection (sv_user_id);


        -- ##################################
        --  table: collection
        -- ##################################
        CREATE TABLE IF NOT EXISTS collection (
            id UUID PRIMARY KEY);
        ALTER TABLE collection
            ADD COLUMN IF NOT EXISTS sv_user_id UUID
            REFERENCES sv_user(id);
        ALTER TABLE collection
            ADD COLUMN IF NOT EXISTS parent_id UUID
            REFERENCES collection(id);
        ALTER TABLE collection ADD COLUMN IF NOT EXISTS type VARCHAR (200);
        ALTER TABLE collection ADD COLUMN IF NOT EXISTS position
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE collection ADD COLUMN IF NOT EXISTS name VARCHAR (200);
        ALTER TABLE collection ADD COLUMN IF NOT EXISTS description VARCHAR (200);
        ALTER TABLE collection ADD COLUMN IF NOT EXISTS keywords VARCHAR (200);
        ALTER TABLE collection ADD COLUMN IF NOT EXISTS link_url VARCHAR (200);
        ALTER TABLE collection ADD COLUMN IF NOT EXISTS link_text VARCHAR (200);
        -- indexes
        CREATE INDEX IF NOT EXISTS collection_sv_user_id_idx ON collection (sv_user_id);
        CREATE INDEX IF NOT EXISTS collection_type_idx ON collection (type);

        -- ##################################
        --  table: place
        -- ##################################
        CREATE TABLE IF NOT EXISTS place (
            id UUID PRIMARY KEY);
        ALTER TABLE place
            ADD COLUMN IF NOT EXISTS sv_user_id UUID
            REFERENCES sv_user(id);
        ALTER TABLE place ADD COLUMN IF NOT EXISTS type VARCHAR (100);
        ALTER TABLE place ADD COLUMN IF NOT EXISTS name VARCHAR (100);
        ALTER TABLE place ADD COLUMN IF NOT EXISTS description VARCHAR (100);
        ALTER TABLE place ADD COLUMN IF NOT EXISTS phone VARCHAR (20);
        ALTER TABLE place ADD COLUMN IF NOT EXISTS email VARCHAR (100);
        ALTER TABLE place ADD COLUMN IF NOT EXISTS website VARCHAR (100);
        ALTER TABLE place ADD COLUMN IF NOT EXISTS domain VARCHAR (100);
        -- indexes
        CREATE INDEX IF NOT EXISTS place_sv_user_id_idx ON place (sv_user_id);
        CREATE INDEX IF NOT EXISTS place_type_idx ON place (type);

        -- ##################################
        --  table: address
        -- ##################################
        CREATE TABLE IF NOT EXISTS address (
            id UUID PRIMARY KEY);
        ALTER TABLE address
            ADD COLUMN IF NOT EXISTS sv_user_id UUID
            REFERENCES sv_user(id);
        ALTER TABLE address
            ADD COLUMN IF NOT EXISTS place_id UUID
            REFERENCES place(id);
        ALTER TABLE address ADD COLUMN IF NOT EXISTS type VARCHAR (100);
        ALTER TABLE address ADD COLUMN IF NOT EXISTS street1 VARCHAR (100);
        ALTER TABLE address ADD COLUMN IF NOT EXISTS street2 VARCHAR (100);
        ALTER TABLE address ADD COLUMN IF NOT EXISTS city VARCHAR (100);
        ALTER TABLE address ADD COLUMN IF NOT EXISTS state VARCHAR (50);
        ALTER TABLE address ADD COLUMN IF NOT EXISTS zipcode VARCHAR (20);
        ALTER TABLE address ADD COLUMN IF NOT EXISTS country VARCHAR (50);
        -- indexes
        CREATE INDEX IF NOT EXISTS address_sv_user_id_idx ON address (sv_user_id);
        CREATE INDEX IF NOT EXISTS address_type_idx ON address (type);

        -- ##################################
        -- table: item
        -- ##################################
        CREATE TABLE IF NOT EXISTS item (
            id UUID PRIMARY KEY);
        ALTER TABLE item
            ADD COLUMN IF NOT EXISTS sv_user_id UUID
            REFERENCES sv_user(id);
        ALTER TABLE item
            ADD COLUMN IF NOT EXISTS parent_id UUID
            REFERENCES item(id);
        ALTER TABLE item
            ADD COLUMN IF NOT EXISTS place_id UUID
            REFERENCES place(id);
        ALTER TABLE item
            ADD COLUMN IF NOT EXISTS price_class_id UUID;
        ALTER TABLE item
            ADD COLUMN IF NOT EXISTS unit_id UUID;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS type VARCHAR (50);
        ALTER TABLE item ADD COLUMN IF NOT EXISTS is_variable
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS is_bundle
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS position
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS sku VARCHAR (50);
        ALTER TABLE item ADD COLUMN IF NOT EXISTS name VARCHAR (200);
        ALTER TABLE item ADD COLUMN IF NOT EXISTS description VARCHAR (200);
        ALTER TABLE item ADD COLUMN IF NOT EXISTS keywords VARCHAR (200);
        ALTER TABLE item ADD COLUMN IF NOT EXISTS cost
            BIGINT NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS cost_override
            BIGINT NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS price BIGINT NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS price_override
            BIGINT NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS price_is_fixed
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS price_discount
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS quantity_available
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS quantity_min
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS quantity_max
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS length
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS width
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS height
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS weight
            INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS file_name VARCHAR (100);
        ALTER TABLE item ADD COLUMN IF NOT EXISTS file_path VARCHAR (100);
        ALTER TABLE item ADD COLUMN IF NOT EXISTS download_code VARCHAR (100);
        ALTER TABLE item ADD COLUMN IF NOT EXISTS download_expiration VARCHAR (100);
        -- indexes
        CREATE INDEX IF NOT EXISTS item_sv_user_id_idx ON item (sv_user_id);
        CREATE INDEX IF NOT EXISTS item_priceclassid_idx ON item (price_class_id);
        CREATE INDEX IF NOT EXISTS item_unitid_idx ON item (unit_id);
        CREATE INDEX IF NOT EXISTS item_type_idx ON item (type);

        -- ##################################
        -- table: join_collection_item
        -- ##################################
        CREATE TABLE IF NOT EXISTS join_collection_item (id UUID PRIMARY KEY);
        ALTER TABLE join_collection_item
            ADD COLUMN IF NOT EXISTS sv_user_id UUID
            REFERENCES sv_user(id);
        ALTER TABLE join_collection_item
            ADD COLUMN IF NOT EXISTS collection_id UUID
            REFERENCES collection(id);
        ALTER TABLE join_collection_item
            ADD COLUMN IF NOT EXISTS item_id UUID
            REFERENCES item(id);
        ALTER TABLE join_collection_item ADD COLUMN IF NOT EXISTS position
            INTEGER NOT NULL DEFAULT 0;
        -- indexes
        CREATE INDEX IF NOT EXISTS jci_sv_user_id_idx ON join_collection_item (sv_user_id);
        CREATE INDEX IF NOT EXISTS jci_itemid_idx ON join_collection_item (item_id);
        CREATE INDEX IF NOT EXISTS jci_collectionid_idx ON join_collection_item (collection_id);

        -- ##################################
        -- table: image
        -- ##################################
        CREATE TABLE IF NOT EXISTS image (id UUID PRIMARY KEY);
        ALTER TABLE image
            ADD COLUMN IF NOT EXISTS sv_user_id UUID NOT NULL
            REFERENCES sv_user(id);
        ALTER TABLE image
            ADD COLUMN IF NOT EXISTS collection_id UUID
            REFERENCES collection(id);
        ALTER TABLE image
            ADD COLUMN IF NOT EXISTS item_id UUID
            REFERENCES item(id);
        ALTER TABLE image ADD COLUMN IF NOT EXISTS url VARCHAR (200);
        ALTER TABLE image ADD COLUMN IF NOT EXISTS size VARCHAR (20);
        ALTER TABLE image ADD COLUMN IF NOT EXISTS position INTEGER NOT NULL DEFAULT 0;
        ALTER TABLE image ADD COLUMN IF NOT EXISTS height VARCHAR (20);
        ALTER TABLE image ADD COLUMN IF NOT EXISTS width VARCHAR (20);
        ALTER TABLE image ADD COLUMN IF NOT EXISTS title VARCHAR (200);
        ALTER TABLE image ADD COLUMN IF NOT EXISTS alt VARCHAR (200);
        ALTER TABLE image ADD COLUMN IF NOT EXISTS caption VARCHAR (200);
        -- indexes
        CREATE INDEX IF NOT EXISTS image_sv_user_id_idx ON image (sv_user_id);
        CREATE INDEX IF NOT EXISTS image_collectionid_idx ON image (collection_id);
        CREATE INDEX IF NOT EXISTS image_itemid_idx ON image (item_id);
        CREATE UNIQUE INDEX
            IF NOT EXISTS image_collection_item_position_size_idx
            ON image (collection_id, item_id, position, size);

        `

        conf := configs.Config{}
        err := configs.Load(&conf)
        if err != nil {
            log.Print("Configs error", err)
        }



        // open each db
        devPostgres := tools.PostgresDb{}
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
