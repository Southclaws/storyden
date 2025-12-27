package moderation

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/moderation/checker"
	"github.com/Southclaws/storyden/app/services/moderation/length_checker"
	"github.com/Southclaws/storyden/app/services/moderation/spam_checker"
	"github.com/Southclaws/storyden/app/services/moderation/word_checker"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			// Spam detector
			spam_checker.New,

			// Register individual checkers
			length_checker.NewLengthChecker,
			spam_checker.NewSpamChecker,
			word_checker.NewWordChecker,

			// Build the registry with all checkers
			func(
				lengthChecker *length_checker.LengthChecker,
				spamChecker *spam_checker.SpamChecker,
				wordChecker *word_checker.WordChecker,
			) *checker.Registry {
				return checker.NewRegistry(
					lengthChecker,
					spamChecker,
					wordChecker,
				)
			},

			// Provide the manager
			New,
		),
	)
}
