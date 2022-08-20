package main

import (
	"context"
	"database/sql"
	"log"
	"text/template"
	"time"

	"ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
)

var dburl = "postgres://default:default@localhost:5432/postgres?sslmode=disable"

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	err := entc.Generate(
		"./api/src/infra/db/schema", &gen.Config{
			Target:  "./api/src/infra/db/model",
			Package: "github.com/Southclaws/storyden/backend/internal/infrastructure/db/model",
		},
		entc.FeatureNames(
			"sql/versioned-migration",
			"sql/upsert",
			"sql/modifier",
			"sql/upsert",
		),
	)
	if err != nil {
		return err
	}

	driver, err := sql.Open("pgx", dburl)
	if err != nil {
		return err
	}

	client := model.NewClient(model.Driver(entsql.OpenDB(dialect.Postgres, driver)))
	defer client.Close()

	// d, err := migrate.NewLocalDir("migrations")
	// if err != nil {
	// 	return err
	// }

	// // Write migration diff.
	// err = client.Schema.Diff(context.Background(), schema.WithDir(d), schema.WithFormatter(GolangMigrateFormatter))
	// if err != nil {
	// 	return err
	// }

	err = client.Schema.Create(context.Background(),
		schema.WithAtlas(true),
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
	)

	return err
}

var (
	// GolangMigrateFormatter is a migrate.Formatter compatible with golang-migrate/migrate.
	GolangMigrateFormatter = templateFormatter(
		"{{ now }}{{ with .Name }}_{{ . }}{{ end }}.up.sql",
		`{{ range .Changes }}{{ with .Comment }}-- {{ println . }}{{ end }}{{ printf "%s;\n" .Cmd }}{{ end }}`,
		"{{ now }}{{ with .Name }}_{{ . }}{{ end }}.down.sql",
		`{{ range rev .Changes }}{{ if .Reverse }}{{ with .Comment }}-- reverse: {{ println . }}{{ end }}{{ printf "%s;\n" .Reverse }}{{ end }}{{ end }}`,
	)
	// funcs contains the template.FuncMap for the different formatters.
	funcs = template.FuncMap{
		"inc": func(x int) int { return x + 1 },
		// now format the current time in a lexicographically ascending order while maintaining human readability.
		"now": func() string { return time.Now().Format("20060102150405") },
		"rev": reverse,
	}
)

// templateFormatter parses the given templates and passes them on to the migrate.NewTemplateFormatter.
func templateFormatter(templates ...string) migrate.Formatter {
	tpls := make([]*template.Template, len(templates))
	for i, t := range templates {
		tpls[i] = template.Must(template.New("").Funcs(funcs).Parse(t))
	}

	fmt, err := migrate.NewTemplateFormatter(tpls...)
	if err != nil {
		panic(err)
	}

	return fmt
}

// reverse changes for the down migration.
func reverse(changes []*migrate.Change) []*migrate.Change {
	n := len(changes)
	rev := make([]*migrate.Change, n)

	if n%2 == 1 {
		rev[n/2] = changes[n/2]
	}

	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = changes[j], changes[i]
	}

	return rev
}
