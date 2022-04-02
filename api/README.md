# Server

This directory contains the API codebase. It's written in Go and is fully cross-platform (does not depend on any platform-specific libraries or C code).

For a list of required environment variables, see the `Config` struct in `src/app/app.go`. You can also set `LOG_LEVEL` to `debug` while developing for verbose log output.

## Architecture

### Layout

The codebase follows idiomatic Go guidelines. Interfaces are preferred as a means of sharing functionality in order to facilitate easier mocking and testing. `src/app/app.go` contains the main setup code in a function named `Initialise`. The app itself starts up and blocks when `App.Start` is called. This fires up each of the components in the app in a goroutine then blocks until the root context is cancelled.

### HTTP API

The library of choice for HTTP routing etc is go-chi. In `src/api/` there is a directory for each API "Module" and the names of these must correspond to its respective API route.

Each handler goes in its own file, prefixed with `h_` followed by a simple description or HTTP method of what it does. The setup code for each API module lives in a file named `api.go`.

For example:

```
api/
├module1/
│ ├api.go
│ ├h_get.go
│ └h_post.go
└module2/
  ├api.go
  ├h_get.go
  ├h_post.go
  └h_search.go
```

The naming for these files is not super strict, most of the time a HTTP method such as "get" or "post" should be fine, but other simple verbs are fine too. What matters is consistency and ease of navigation.

### Database

The database code is managed by Prisma. The schema is declared in the `prisma` directory in the root.

The workflow to make changes to the data model is:

1. Make the change
2. Generate new Go code
3. Generate migrations
4. Apply migrations to dev instance

Once a pull request is approved, if it has data model changes, an admin will apply those changes to the live DB.
