package seed

import (
	"context"
	"fmt"
	"strings"

	"github.com/minimaxir/big-list-of-naughty-strings/naughtystrings"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/internal/ent"
)

var (
	Post_01_Welcome = thread.Thread{
		ID:       post.PostID(id("00000000000000000010")),
		Title:    "Welcome to Storyden!",
		Author:   post.Author{ID: Account_001_Odin.ID},
		Category: Category_01_General,
		Posts: []*post.Post{
			{
				Body: `Storyden is a platform for building communities.

But not just another chat app or another forum site. Storyden is a modern take on oldschool bulletin board forums you may remember from the earlier days of the internet.

With a fresh perspective and new objectives, Storyden is designed to be the community platform for the next era of internet culture.

## Why Storyden for people?

There's a huge lack of focus on accessibility with a lot of modern discussion platforms. And this isn't just about alt-tags and screen readers, it's about crafting a standards conformant web application that runs on the bare minimum hardware without grinding to a halt. That means progressively enhanced, server-side-rendered, HTML-first, simple yet extensible and ready to run anywhere.

Privacy is another factor, it's not "becoming more important", it always has been and always will be. Storyden does not use email or phone numbers as the fundamental unit of identification. You can if you want, or you can just go username-only. Or you can go full web3. Or you can enable the new WebAuthn authentication and sign in with your fingerprint. Or sign in with your favourite socials. The options are there but the default is privacy-first.

And finally, we just want to build a quality desktop and web experience that works how you'd expect.

## Why Storyden for system administrators or programmers?

Simple and minimal operational overhead is the primary technical goal of Storyden. No need to compile or build your own Docker image or run various services.

Storyden ships as a single static binary or container image that runs almost anywhere.

And if you don't like the user interface, that's fine too! You or your team/community/organisation can run the API in headless mode and build your own using the fully documented OpenAPI specification.

## Open source

And obviously it's open source with a permissive license. Fork it, find bugs, contribute fixes, spin up a hosting company and run instances for your customers if you want!

The code aims to be simple and accessible for either experienced software engineers to dive in and edit or for newcomers to programming to read, learn from and contribute to.

## Future

Storyden is still in development so please give the repository a watch if you're interested!
`,
			},
			{
				ID:         post.PostID(id("00000000000000001010")),
				Body:       "first 😁",
				RootPostID: post.PostID(id("00000000000000000010")),
				Author:     post.Author{ID: Account_004_Loki.ID},
			},
			{
				ID:         post.PostID(id("00000000000000002010")),
				Body:       "Nice! One question: what kind of formatting can you use in posts? Is it like the old days with [b]tags[/b] and [color=red]cool stuff[/color] like that?",
				RootPostID: post.PostID(id("00000000000000000010")),
				Author:     post.Author{ID: Account_002_Frigg.ID},
			},
			{
				ID:         post.PostID(id("00000000000000003010")),
				Body:       "Good question @frigg, we're probably going to use Markdown with some basic extensions but nothing is set in stone yet.",
				RootPostID: post.PostID(id("00000000000000000010")),
				Author:     post.Author{ID: Account_001_Odin.ID},
			},
			{
				ID:         post.PostID(id("00000000000000004010")),
				Body:       "What about images and stuff?",
				RootPostID: post.PostID(id("00000000000000000010")),
				Author:     post.Author{ID: Account_008_Heimdallr.ID},
			},
			{
				ID: post.PostID(id("00000000000000005010")),
				Body: `oh you can do that like this:

![https://i.imgur.com/gl39KB7.png](https://i.imgur.com/gl39KB7.png)
`,
				RootPostID: post.PostID(id("00000000000000000010")),
				Author:     post.Author{ID: Account_004_Loki.ID},
			},
			{
				ID:         post.PostID(id("00000000000000006010")),
				Body:       `how did you do that??`,
				RootPostID: post.PostID(id("00000000000000000010")),
				Author:     post.Author{ID: Account_005_Þórr.ID},
			},
			{
				ID:         post.PostID(id("00000000000000007010")),
				Body:       `haha secret 😉`,
				RootPostID: post.PostID(id("00000000000000000010")),
				Author:     post.Author{ID: Account_004_Loki.ID},
			},
			{
				ID: post.PostID(id("00000000000000008010")),
				Body: `It was mentioned above, use markdown:

https://daringfireball.net/markdown
`,
				RootPostID: post.PostID(id("00000000000000000010")),
				Author:     post.Author{ID: Account_002_Frigg.ID},
			},
			{
				ID:         post.PostID(id("00000000000000009010")),
				Body:       "Thanks guys!",
				RootPostID: post.PostID(id("00000000000000000010")),
				Author:     post.Author{ID: Account_008_Heimdallr.ID},
			},
		},
	}
	Post_02_HowToContribute = thread.Thread{
		ID:       post.PostID(id("00000000000000000020")),
		Title:    "How to contribute",
		Author:   post.Author{ID: Account_001_Odin.ID},
		Category: Category_01_General,
		Posts: []*post.Post{
			{
				Body: `This post contains a list of resources for those of you who wish to contribute to Storyden.

What does contribution mean? Anything, large or small! Even if you spot a typo in the home page or in this demo data you can report it or even take a swing at fixing it!

If you're new to open source, don't be shy and ask for guidance on how to solve a problem you or someone else has found.

The main place for reporting issues or making feature requests is here:

https://github.com/Southclaws/storyden/issues

You can also scout out what's in-progress and offer feedback or support here:

https://github.com/Southclaws/storyden/pulls

And there's also a great community that's friends with Storyden called Makeroom, they run a Discord server where you can ask questions and get support for anything related to building digital products:

https://makeroom.club

If I've missed anything, post in this thread and I'll add it here 😃
`,
			},
			{
				ID:         post.PostID(id("00000000000000001020")),
				Body:       "Is there a wiki?",
				RootPostID: post.PostID(id("00000000000000000020")),
				Author:     post.Author{ID: Account_006_Freyja.ID},
			},
			{
				ID:         post.PostID(id("00000000000000002020")),
				Body:       "Not yet but they're working on it!",
				RootPostID: post.PostID(id("00000000000000000020")),
				Author:     post.Author{ID: Account_002_Frigg.ID},
			},
		},
	}

	Post_03_LoremIpsum = thread.Thread{
		ID:       post.PostID(id("00000000000000000030")),
		Title:    "The lorem ipsum thread",
		Author:   post.Author{ID: Account_005_Þórr.ID},
		Category: Category_01_General,
		Posts: []*post.Post{
			{
				Body: `In this thread:

Try to break storyden with large amounts of text, hacky strings, etc! GO!`,
			},
			{
				ID:         post.PostID(id("00000000000000001030")),
				Body:       "ooh fun! my favourite tool for this is: https://jaspervdj.be/lorem-markdownum/\n\n" + markdownTest01,
				RootPostID: post.PostID(id("00000000000000000030")),
				Author:     post.Author{ID: Account_006_Freyja.ID},
			},
			{
				ID:         post.PostID(id("00000000000000002030")),
				Body:       "That's pretty useful, here's mine:\n\n" + markdownTest02,
				RootPostID: post.PostID(id("00000000000000000030")),
				Author:     post.Author{ID: Account_007_Freyr.ID},
			},
			{
				ID:         post.PostID(id("00000000000000003030")),
				Body:       "nah that's useless, you guys need some real hacky stuff to properly test:\n\n" + strings.Join(naughtystrings.Unencoded(), "\n\n"),
				RootPostID: post.PostID(id("00000000000000000030")),
				Author:     post.Author{ID: Account_004_Loki.ID},
			},
		},
	}

	Threads = []thread.Thread{
		Post_01_Welcome,
		Post_02_HowToContribute,
		Post_03_LoremIpsum,
	}
)

func threads(tr thread.Repository, pr post.Repository) {
	ctx := context.Background()

	for _, t := range Threads {
		first := t.Posts[0]

		th, err := tr.Create(ctx,
			t.Title,
			first.Body,
			t.Author.ID,
			t.Category.ID,
			t.Tags,
			thread.WithID(t.ID))
		if err != nil {
			if ent.IsConstraintError(err) {
				continue
			}
			panic(err)
		}

		for _, p := range t.Posts[1:] {
			_, err = pr.Create(ctx,
				p.Body,
				p.Author.ID,
				th.ID,
				nil,
				nil,
				post.WithID(p.ID))
			if err != nil {
				if ent.IsConstraintError(err) {
					continue
				}
				panic(err)
			}
		}
	}

	fmt.Println("created seed threads and posts")
}
