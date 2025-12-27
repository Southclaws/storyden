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
	Services: opt.New(ServiceSettings{
		Moderation: opt.New(ModerationServiceSettings{
			ThreadBodyLengthMax: opt.New(60000),
			ReplyBodyLengthMax:  opt.New(10000),
		}),
	}),
}
