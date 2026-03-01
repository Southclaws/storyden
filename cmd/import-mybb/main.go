package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/glebarez/go-sqlite"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Southclaws/storyden/cmd/import-mybb/loader"
	"github.com/Southclaws/storyden/cmd/import-mybb/logger"
	"github.com/Southclaws/storyden/cmd/import-mybb/transform"
	"github.com/Southclaws/storyden/cmd/import-mybb/writer"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/db"
)

func main() {
	mybbDSN := flag.String("mybb-dsn", "mybbuser:mybbpassword@tcp(127.0.0.1:3306)/mybb?charset=latin1", "MyBB MySQL connection string")
	storydenDB := flag.String("storyden-db", "sqlite://data/data.db?_pragma=foreign_keys(1)", "Storyden database URL (SQLite or Postgres)")
	dryRun := flag.Bool("dry-run", false, "Validate without writing to database")
	batchSize := flag.Int("batch-size", 1000, "Bulk insert batch size")
	skipAssets := flag.Bool("skip-assets", true, "Skip asset file imports")
	flag.Parse()

	ctx := context.Background()

	if err := run(ctx, *mybbDSN, *storydenDB, *dryRun, *batchSize, *skipAssets); err != nil {
		log.Fatalf("Import failed: %v", err)
	}

	log.Println("Import completed successfully!")
}

func run(ctx context.Context, mybbDSN, storydenDB string, dryRun bool, batchSize int, skipAssets bool) error {
	mybb, err := sql.Open("mysql", mybbDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to MyBB database: %w", err)
	}
	defer mybb.Close()

	if err := mybb.Ping(); err != nil {
		return fmt.Errorf("failed to ping MyBB database: %w", err)
	}

	cfg := config.Config{DatabaseURL: storydenDB}

	sqlDriver, _, err := db.NewSQL(cfg)
	if err != nil {
		return fmt.Errorf("failed to create ent client: %w", err)
	}

	storydenClient, err := db.Connect(ctx, cfg, sqlDriver)
	if err != nil {
		return fmt.Errorf("failed to create ent client: %w", err)
	}
	defer storydenClient.Close()

	log.Println("Running schema migrations...")
	if err := storydenClient.Schema.Create(ctx, schema.WithDropColumn(true), schema.WithDropIndex(true)); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("Loading MyBB data...")
	mybbData, err := loader.LoadAll(ctx, mybb)
	if err != nil {
		return fmt.Errorf("failed to load MyBB data: %w", err)
	}

	log.Printf("Loaded: %d users, %d forums, %d threads, %d posts, %d usergroups",
		len(mybbData.Users), len(mybbData.Forums), len(mybbData.Threads), len(mybbData.Posts), len(mybbData.UserGroups))

	if dryRun {
		log.Println("Dry run mode - skipping database writes")
		return nil
	}

	log.Println("Transforming and importing data...")
	w := writer.New(storydenClient, batchSize)

	logger.Phase(-1, "Clearing existing Storyden data")
	if err = w.DeleteAllData(ctx); err != nil {
		return fmt.Errorf("failed to clear existing data: %w", err)
	}

	logger.Phase(0, "Settings (from mybb_settings)")
	if err := transform.ImportSettings(ctx, w, mybbData); err != nil {
		return fmt.Errorf("failed to import settings: %w", err)
	}

	logger.Phase(1, "Roles (from mybb_usergroups)")
	if err := transform.ImportRoles(ctx, w, mybbData); err != nil {
		return fmt.Errorf("failed to import roles: %w", err)
	}

	logger.Phase(2, "Accounts (from mybb_users)")
	if err := transform.ImportAccounts(ctx, w, mybbData); err != nil {
		return fmt.Errorf("failed to import accounts: %w", err)
	}

	logger.Phase(3, "Categories (from mybb_forums)")
	if err := transform.ImportCategories(ctx, w, mybbData); err != nil {
		return fmt.Errorf("failed to import categories: %w", err)
	}

	logger.Phase(4, "Tags (from mybb_threadprefixes)")
	if err := transform.ImportTags(ctx, w, mybbData); err != nil {
		return fmt.Errorf("failed to import tags: %w", err)
	}

	logger.Phase(5, "Posts - Threads (from mybb_threads)")
	if err := transform.ImportThreads(ctx, w, mybbData); err != nil {
		return fmt.Errorf("failed to import threads: %w", err)
	}

	logger.Phase(6, "Posts - Replies (from mybb_posts)")
	if err := transform.ImportPosts(ctx, w, mybbData); err != nil {
		return fmt.Errorf("failed to import posts: %w", err)
	}

	logger.Phase(7, "Interactions (reacts, likes, reads, reports)")
	if err := transform.ImportInteractions(ctx, w, mybbData); err != nil {
		return fmt.Errorf("failed to import interactions: %w", err)
	}

	if !skipAssets {
		logger.Phase(8, "Assets (from mybb_attachments)")
		if err := transform.ImportAssets(ctx, w, mybbData); err != nil {
			return fmt.Errorf("failed to import assets: %w", err)
		}
	}

	logger.Success("Import completed successfully!")
	return nil
}
