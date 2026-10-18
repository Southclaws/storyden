package settings

import (
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/datagraph"
)

const (
	DefaultTitle       = "Storyden"
	DefaultDescription = "A forum for the modern age"
	DefaultContent     = `<body>
<p>Welcome to your new community!</p>
<p>You can edit this content by clicking Edit below.</p>
<p>This is a <em>rich text section</em> for telling visitors what your community is about.</p>
<p>Add a link to your <a href="https://discord.gg/XF6ZBGF9XF">Discord</a> or other sites.</p>
<p>Enjoy!</p>
</body>`
)
const DefaultColour = "hsl(157, 65%, 44%)"

// skip error check, we know it's correct, it's literally above ^^
var defaultContent, _ = datagraph.NewRichText(DefaultContent)

var DefaultSettings = Settings{
	Title:              opt.New(DefaultTitle),
	Description:        opt.New(DefaultDescription),
	Content:            opt.New(defaultContent),
	AccentColour:       opt.New(DefaultColour),
	AuthenticationMode: opt.New(authentication.ModeHandle),
	RegistrationMode:   opt.New(RegistrationModePublic),
	Services: opt.New(ServiceSettings{
		Media: opt.New(MediaServiceSettings{
			ImageMaxSize:   opt.New[int64](8 << 20),
			ImageMaxPixels: opt.New[int64](1500000),
		}),
		Discovery: opt.New(DiscoveryServiceSettings{
			SearchIndexing:   opt.New(true),
			TrainingCrawling: opt.New(false),
			LLMSTxt:          opt.New(false),
			AuthMD:           opt.New(false),
			AICatalog:        opt.New(false),
			APICatalog:       opt.New(false),
			MCPServerCard:    opt.New(false),
			RobotsTxt:        opt.New(""),
		}),
		Moderation: opt.New(ModerationServiceSettings{
			ThreadBodyLengthMax: opt.New(60000),
			ReplyBodyLengthMax:  opt.New(10000),
			SignatureLengthMax:  opt.New(500),
		}),
	}),
}
