package moderation

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/moderation/checker"
	"github.com/Southclaws/storyden/app/services/moderation/length_checker"
	"github.com/Southclaws/storyden/app/services/moderation/spam_checker"
	"github.com/Southclaws/storyden/app/services/moderation/wordblock_checker"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			// Spam detector
			spam_checker.New,

			// Register individual checkers
			length_checker.NewLengthChecker,
			spam_checker.NewSpamChecker,
			wordblock_checker.NewWordBlockChecker,

			// Build the registry with all checkers
			func(
				lengthChecker *length_checker.LengthChecker,
				spamChecker *spam_checker.SpamChecker,
				wordBlockChecker *wordblock_checker.WordBlockChecker,
			) *checker.Registry {
				return checker.NewRegistry(
					lengthChecker,
					spamChecker,
					wordBlockChecker,
				)
			},

			// Provide the manager
			New,
		),
	)
}
