package graphql

import "github.com/Southclaws/storyden/app/services/thread"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	thread_service thread.Service
}
