package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/jeromewir/flatos/cmd/syncflats/providers"
	"github.com/jeromewir/flatos/internal/config"
	"github.com/jeromewir/flatos/internal/db/queries"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite" // Import the SQLite driver
)

func openSQLConnection() (*sql.DB, error) {
	// Replace with your actual database connection logic
	db, err := sql.Open("sqlite", "./database/flatos.sqlite")

	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	var url string

	config := config.NewConfig()

	if err := config.Load(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	sqlDB, err := openSQLConnection()

	if err != nil {
		fmt.Printf("Error opening SQL connection: %v\n", err)
		return
	}
	defer sqlDB.Close()

	q := queries.New(sqlDB)

	ctx := context.Background()

	var rootCmd = &cobra.Command{
		Use:   "syncflats",
		Short: "Sync flats CLI",
		Run: func(cmd *cobra.Command, args []string) {
			if url == "" {
				fmt.Println("No URL provided. Please provide a URL using the --url flag.")
				return
			}

			provider, err := providers.Get(config, url)

			if err != nil {
				fmt.Printf("Error getting provider: %v\n", err)
				return
			}

			flats, err := provider.GetFlats(url)

			if err != nil {
				fmt.Printf("Error getting flats: %v\n", err)
				return
			}

			fmt.Println("Flats retrieved successfully:", len(flats))

			for _, flat := range flats {
				id, err := q.UpsertFlat(ctx, queries.UpsertFlatParams(flat))

				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					fmt.Printf("Error upserting flat: %v\n", err)
					return
				}

				// Inserted
				if id != "" {
					fmt.Printf("Flat ID: %s, Name: %s, Address: %s, City: %s, Price: %d\n",
						flat.ID, flat.Name, flat.Address, flat.City, flat.Price)
				}

			}
		},
	}

	rootCmd.Flags().StringVarP(&url, "url", "u", "", "Optional URL to use")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
