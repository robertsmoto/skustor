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

    "github.com/robertsmoto/skustor/internal/postgres"
    "github.com/robertsmoto/skustor/internal/configs"
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
        --  table: account
        -- ##################################
        CREATE TABLE IF NOT EXISTS account (
            id UUID PRIMARY KEY);
        ALTER TABLE account
            ADD COLUMN IF NOT EXISTS parent_id UUID
            REFERENCES account(id)
            ON DELETE CASCADE;
        ALTER TABLE account ADD COLUMN IF NOT EXISTS auth UUID;
        ALTER TABLE account ADD COLUMN IF NOT EXISTS key UUID;
        ALTER TABLE account ADD COLUMN IF NOT EXISTS document JSONB;
        -- indexes
        CREATE INDEX IF NOT EXISTS account_parent_id_idx ON account (parent_id);
        CREATE INDEX IF NOT EXISTS account_auth_idx ON account (auth);
        CREATE INDEX IF NOT EXISTS account_key_idx ON account (key);

        -- ##################################
        --  table: collection
        -- ##################################
        CREATE TABLE IF NOT EXISTS collection (
            id UUID PRIMARY KEY);
        ALTER TABLE collection
            ADD COLUMN IF NOT EXISTS account_id UUID
            REFERENCES account(id)
            ON DELETE CASCADE;
        ALTER TABLE collection
            ADD COLUMN IF NOT EXISTS parent_id UUID
            REFERENCES collection(id)
            ON DELETE CASCADE;
        ALTER TABLE collection ADD COLUMN IF NOT EXISTS type VARCHAR (20);
        ALTER TABLE collection ADD COLUMN IF NOT EXISTS document JSONB;
        -- indexes

        CREATE INDEX IF NOT EXISTS collection_id_idx ON collection (id);
        CREATE INDEX IF NOT EXISTS collection_account_id_idx ON collection (account_id);
        CREATE INDEX IF NOT EXISTS collection_parent_id_idx ON collection (parent_id);
        CREATE INDEX IF NOT EXISTS collection_type_idx ON collection (type);

        -- ##################################
        --  table: image
        -- ##################################
        CREATE TABLE IF NOT EXISTS image (id UUID PRIMARY KEY);
        ALTER TABLE image
            ADD COLUMN IF NOT EXISTS account_id UUID
            REFERENCES account(id)
            ON DELETE CASCADE;
        ALTER TABLE image
            ADD COLUMN IF NOT EXISTS parent_id UUID
            REFERENCES image(id)
            ON DELETE CASCADE;
        ALTER TABLE image ADD COLUMN IF NOT EXISTS type VARCHAR (20);
        ALTER TABLE image ADD COLUMN IF NOT EXISTS document JSONB;
        -- indexes

        CREATE INDEX IF NOT EXISTS image_id_idx ON image (id);
        CREATE INDEX IF NOT EXISTS image_parent_id_idx ON image (parent_id);
        CREATE INDEX IF NOT EXISTS image_account_id_idx ON image (account_id);
        CREATE INDEX IF NOT EXISTS image_type_idx ON image (type);

        -- ##################################
        --  table: content
        -- ##################################
        CREATE TABLE IF NOT EXISTS content (
            id UUID PRIMARY KEY);
        ALTER TABLE content
            ADD COLUMN IF NOT EXISTS account_id UUID
            REFERENCES account(id)
            ON DELETE CASCADE;
        ALTER TABLE content
            ADD COLUMN IF NOT EXISTS parent_id UUID
            REFERENCES content(id)
            ON DELETE CASCADE;
        ALTER TABLE content ADD COLUMN IF NOT EXISTS type VARCHAR (20);
        ALTER TABLE content ADD COLUMN IF NOT EXISTS document JSONB;
        -- indexes

        CREATE INDEX IF NOT EXISTS content_id_idx ON content (id);
        CREATE INDEX IF NOT EXISTS content_parent_id_idx ON content (parent_id);
        CREATE INDEX IF NOT EXISTS content_account_id_idx ON content (account_id);
        CREATE INDEX IF NOT EXISTS content_type_idx ON content (type);
        
        -- ##################################
        --  table: place
        -- ##################################
        CREATE TABLE IF NOT EXISTS place (
            id UUID PRIMARY KEY);
        ALTER TABLE place
            ADD COLUMN IF NOT EXISTS account_id UUID
            REFERENCES account(id)
            ON DELETE CASCADE;
        ALTER TABLE place
            ADD COLUMN IF NOT EXISTS parent_id UUID
            REFERENCES content(id)
            ON DELETE CASCADE;
        ALTER TABLE place ADD COLUMN IF NOT EXISTS type VARCHAR (20);
        ALTER TABLE place ADD COLUMN IF NOT EXISTS document JSONB;
        -- indexes

        CREATE INDEX IF NOT EXISTS place_id_idx ON place (id);
        CREATE INDEX IF NOT EXISTS place_parent_id_idx ON place (parent_id);
        CREATE INDEX IF NOT EXISTS place_account_id_idx ON place (account_id);
        CREATE INDEX IF NOT EXISTS place_type_idx ON place (type);

        -- ##################################
        --  table: person
        -- ##################################
        CREATE TABLE IF NOT EXISTS person (
            id UUID PRIMARY KEY);
        ALTER TABLE person
            ADD COLUMN IF NOT EXISTS account_id UUID
            REFERENCES account(id)
            ON DELETE CASCADE;
        ALTER TABLE person
            ADD COLUMN IF NOT EXISTS parent_id UUID
            REFERENCES content(id)
            ON DELETE CASCADE;
        ALTER TABLE person
            ADD COLUMN IF NOT EXISTS place_id UUID
            REFERENCES place(id)
            ON DELETE CASCADE;
        ALTER TABLE person ADD COLUMN IF NOT EXISTS type VARCHAR (20);
        ALTER TABLE person ADD COLUMN IF NOT EXISTS document JSONB;
        -- indexes

        CREATE INDEX IF NOT EXISTS person_id_idx ON person (id);
        CREATE INDEX IF NOT EXISTS person_parent_id_idx ON person (parent_id);
        CREATE INDEX IF NOT EXISTS person_account_id_idx ON person (account_id);
        CREATE INDEX IF NOT EXISTS person_type_idx ON person (type);

        -- ##################################
        --  table: item
        -- ##################################
        CREATE TABLE IF NOT EXISTS item (
            id UUID PRIMARY KEY);
        ALTER TABLE item
            ADD COLUMN IF NOT EXISTS account_id UUID
            REFERENCES account(id)
            ON DELETE CASCADE;
        ALTER TABLE item
            ADD COLUMN IF NOT EXISTS parent_id UUID
            REFERENCES item(id)
            ON DELETE CASCADE;
        ALTER TABLE item ADD COLUMN IF NOT EXISTS type VARCHAR (20);
        ALTER TABLE item ADD COLUMN IF NOT EXISTS document JSONB;
        -- indexes
        CREATE INDEX IF NOT EXISTS item_id_idx ON item (id);
        CREATE INDEX IF NOT EXISTS item_account_id_idx ON item (account_id);
        CREATE INDEX IF NOT EXISTS item_parent_id_idx ON item (parent_id);
        CREATE INDEX IF NOT EXISTS item_type_idx ON item (type);

        -- ##################################
        -- table: joins
        -- ##################################
        CREATE TABLE IF NOT EXISTS joins (id UUID PRIMARY KEY);
        ALTER TABLE joins
            ADD COLUMN IF NOT EXISTS account_id UUID
            REFERENCES account(id)
            ON DELETE CASCADE;
        ALTER TABLE joins
            ADD COLUMN IF NOT EXISTS collection_id UUID
            REFERENCES collection(id);
        ALTER TABLE joins
            ADD COLUMN IF NOT EXISTS content_id UUID
            REFERENCES content(id);
        ALTER TABLE joins
            ADD COLUMN IF NOT EXISTS image_id UUID
            REFERENCES image(id);
        ALTER TABLE joins
            ADD COLUMN IF NOT EXISTS item_id UUID
            REFERENCES item(id);
        ALTER TABLE joins
            ADD COLUMN IF NOT EXISTS person_id UUID
            REFERENCES person(id);
        ALTER TABLE joins
            ADD COLUMN IF NOT EXISTS place_id UUID
            REFERENCES place(id);
        ALTER TABLE joins ADD COLUMN IF NOT EXISTS attributes JSONB;
        -- indexes
        CREATE INDEX IF NOT EXISTS joins_account_id_idx
            ON joins (account_id);
        CREATE INDEX IF NOT EXISTS joins_collection_id_idx
            ON joins (collection_id);
        CREATE INDEX IF NOT EXISTS joins_content_id_idx
            ON joins (item_id);
        CREATE INDEX IF NOT EXISTS joins_image_id_idx
            ON joins (item_id);
        CREATE INDEX IF NOT EXISTS joins_item_id_idx
            ON joins (item_id);
        CREATE INDEX IF NOT EXISTS joins_person_id_idx
            ON joins (person_id);
        CREATE INDEX IF NOT EXISTS joins_place_id_idx
            ON joins (place_id);
        `

        // load config env variables
        configs.Load(&configs.Config{})

        // open db
        devPostgres := postgres.PostgresDb{}
        pgDb, err := postgres.Open(&devPostgres)

        fmt.Println("## dbconn ", pgDb)
        if err != nil {
            log.Print("Connection error", err)
        }
        _, err = pgDb.Exec(qstr)
        if err != nil {
            log.Print("Error creating or updating database.", err)
        } else {
            log.Print("Successfully updated postgresDb.")
        }

        pgDb.Close()
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
