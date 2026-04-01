package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"ariga.io/atlas/atlasexec"
)

func RunMigrations(migrationsDir string, dbURL string) error {
	// atlasexec works on a temporary working directory.
	// WithMigrations loads your .sql files from disk into that temp dir.
	// This means the Atlas CLI doesn't need to know where your project root is —
	// it reads the files directly from memory via the fs.FS interface.
	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS(migrationsDir),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to load migrations directory: %w", err)
	}
	// atlasexec copies files to a temp dir internally, so we clean it up when done
	defer workdir.Close()

	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		return fmt.Errorf("failed to initialize atlas client: %w", err)
	}

	res, err := client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL: dbURL,
	})
	if err != nil {
		return fmt.Errorf("atlas migrate apply failed: %w", err)
	}

	log.Printf("Applied %d migration(s)\n", len(res.Applied))
	return nil
}
