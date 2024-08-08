package seed

import (
	"context"
	"fmt"
	"html"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/minimaxir/big-list-of-naughty-strings/naughtystrings"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/utils"
)

var (
	Post_01_Welcome = thread.Thread{
		Post: post.Post{
			ID:     post.ID(id("00000000000000000010")),
			Author: profile.Public{ID: Account_001_Odin.ID},
			Content: utils.Must(content.NewRichText(`Storyden is a platform for building communities.

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
`)),
		},
		Title:    "Welcome to Storyden!",
		Category: Category_01_General,
		Replies: []*reply.Reply{
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000001010")),
					Author:  profile.Public{ID: Account_004_Loki.ID},
					Content: utils.Must(content.NewRichText("first 😁")),
				},
				RootPostID: post.ID(id("00000000000000000010")),
			},
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000002010")),
					Author:  profile.Public{ID: Account_002_Frigg.ID},
					Content: utils.Must(content.NewRichText("Nice! One question: what kind of formatting can you use in posts? Is it like the old days with [b]tags[/b] and [color=red]cool stuff[/color] like that?")),
				},
				RootPostID: post.ID(id("00000000000000000010")),
			},
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000003010")),
					Author:  profile.Public{ID: Account_001_Odin.ID},
					Content: utils.Must(content.NewRichText("Good question @frigg, we're probably going to use Markdown with some basic extensions but nothing is set in stone yet.")),
				},
				RootPostID: post.ID(id("00000000000000000010")),
			},
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000004010")),
					Author:  profile.Public{ID: Account_008_Heimdallr.ID},
					Content: utils.Must(content.NewRichText("What about images and stuff?")),
				},
				RootPostID: post.ID(id("00000000000000000010")),
			},
			{
				Post: post.Post{
					ID:     post.ID(id("00000000000000005010")),
					Author: profile.Public{ID: Account_004_Loki.ID},
					Content: utils.Must(content.NewRichText(`oh you can do that like this:
					
					![https://i.imgur.com/gl39KB7.png](https://i.imgur.com/gl39KB7.png)
					`)),
				},
				RootPostID: post.ID(id("00000000000000000010")),
			},
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000006010")),
					Author:  profile.Public{ID: Account_005_Þórr.ID},
					Content: utils.Must(content.NewRichText(`how did you do that??`)),
				},
				RootPostID: post.ID(id("00000000000000000010")),
			},
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000007010")),
					Author:  profile.Public{ID: Account_004_Loki.ID},
					Content: utils.Must(content.NewRichText(`haha secret 😉`)),
				},
				RootPostID: post.ID(id("00000000000000000010")),
			},
			{
				Post: post.Post{
					ID:     post.ID(id("00000000000000008010")),
					Author: profile.Public{ID: Account_002_Frigg.ID},
					Content: utils.Must(content.NewRichText(`It was mentioned above, use markdown:
					
					https://daringfireball.net/markdown
					`)),
				},
				RootPostID: post.ID(id("00000000000000000010")),
			},
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000009010")),
					Author:  profile.Public{ID: Account_008_Heimdallr.ID},
					Content: utils.Must(content.NewRichText("Thanks guys!")),
				},
				RootPostID: post.ID(id("00000000000000000010")),
			},
		},
	}
	Post_02_HowToContribute = thread.Thread{
		Post: post.Post{
			ID:     post.ID(id("00000000000000000020")),
			Author: profile.Public{ID: Account_001_Odin.ID},
			Content: utils.Must(content.NewRichText(`This post contains a list of resources for those of you who wish to contribute to Storyden.
	
	What does contribution mean? Anything, large or small! Even if you spot a typo in the home page or in this demo data you can report it or even take a swing at fixing it!
	
	If you're new to open source, don't be shy and ask for guidance on how to solve a problem you or someone else has found.
	
	The main place for reporting issues or making feature requests is here:
	
	https://github.com/Southclaws/storyden/issues
	
	You can also scout out what's in-progress and offer feedback or support here:
	
	https://github.com/Southclaws/storyden/pulls
	
	And there's also a great community that's friends with Storyden called Makeroom, they run a Discord server where you can ask questions and get support for anything related to building digital products:
	
	https://makeroom.club
	
	If I've missed anything, post in this thread and I'll add it here 😃
	`)),
		},
		Title:    "How to contribute",
		Category: Category_01_General,
		Replies: []*reply.Reply{
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000001020")),
					Content: utils.Must(content.NewRichText("Is there a wiki?")),
					Author:  profile.Public{ID: Account_006_Freyja.ID},
				},
				RootPostID: post.ID(id("00000000000000000020")),
			},
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000002020")),
					Content: utils.Must(content.NewRichText("Not yet but they're working on it!")),
					Author:  profile.Public{ID: Account_002_Frigg.ID},
				},
				RootPostID: post.ID(id("00000000000000000020")),
			},
		},
	}

	Post_03_LoremIpsum = thread.Thread{
		Post: post.Post{
			ID:     post.ID(id("00000000000000000030")),
			Author: profile.Public{ID: Account_005_Þórr.ID},
			Content: utils.Must(content.NewRichText(`In this thread:
	
	Try to break storyden with large amounts of text, hacky strings, etc! GO!`)),
		},
		Title:    "The lorem ipsum thread",
		Category: Category_01_General,
		Replies: []*reply.Reply{
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000001030")),
					Content: utils.Must(content.NewRichText("ooh fun! my favourite tool for this is: https://jaspervdj.be/lorem-markdownum/\n\n" + markdownTest01)),
					Author:  profile.Public{ID: Account_006_Freyja.ID},
				},
				RootPostID: post.ID(id("00000000000000000030")),
			},
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000002030")),
					Content: utils.Must(content.NewRichText(markdownTest03)),
					Author:  profile.Public{ID: Account_002_Frigg.ID},
				},
				RootPostID: post.ID(id("00000000000000000030")),
			},
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000003030")),
					Content: utils.Must(content.NewRichText("That's pretty useful, here's mine:\n\n" + markdownTest02)),
					Author:  profile.Public{ID: Account_007_Freyr.ID},
				},
				RootPostID: post.ID(id("00000000000000000030")),
			},
			{
				Post: post.Post{
					ID:      post.ID(id("00000000000000004030")),
					Content: utils.Must(content.NewRichText("nah that's useless, you guys need some real hacky stuff to properly test:\n\n" + strings.Join(naughtystrings.Unencoded(), "\n\n"))),
					Author:  profile.Public{ID: Account_004_Loki.ID},
				},
				RootPostID: post.ID(id("00000000000000000030")),
			},
		},
	}

	Post_04_Photos = thread.Thread{
		Post: post.Post{
			ID:      post.ID(id("00000000000000000040")),
			Author:  profile.Public{ID: Account_005_Þórr.ID},
			Content: utils.Must(content.NewRichText("some pics from my trip!")),
			Assets: []*asset.Asset{
				{
					ID:   utils.Must(xid.FromString("00000000000000000040")),
					Name: asset.NewFilename("00000000000000000040-asset-01"),
					URL:  "https://pub-7b5607a210bc4f0b81cb6ba41e8754f9.r2.dev/test/IMG_1158.jpg",
					Size: 2537802, Metadata: map[string]any{"mime_type": "image/jpeg", "width": 3024, "height": 4032},
				},
				{
					ID:   utils.Must(xid.FromString("00000000000000000040")),
					Name: asset.NewFilename("00000000000000000040-asset-01"),
					URL:  "https://pub-7b5607a210bc4f0b81cb6ba41e8754f9.r2.dev/test/IMG_1174.jpg",
					Size: 1433625, Metadata: map[string]any{"mime_type": "image/jpeg", "width": 3024, "height": 4032},
				},
				{
					ID:   utils.Must(xid.FromString("00000000000000000040")),
					Name: asset.NewFilename("00000000000000000040-asset-01"),
					URL:  "https://pub-7b5607a210bc4f0b81cb6ba41e8754f9.r2.dev/test/IMG_1236.jpg",
					Size: 1828065, Metadata: map[string]any{"mime_type": "image/jpeg", "width": 3024, "height": 4032},
				},
				{
					ID:   utils.Must(xid.FromString("00000000000000000040")),
					Name: asset.NewFilename("00000000000000000040-asset-01"),
					URL:  "https://pub-7b5607a210bc4f0b81cb6ba41e8754f9.r2.dev/test/IMG_1239.jpg",
					Size: 1769497, Metadata: map[string]any{"mime_type": "image/jpeg", "width": 3024, "height": 4032},
				},
				{
					ID:   utils.Must(xid.FromString("00000000000000000040")),
					Name: asset.NewFilename("00000000000000000040-asset-01"),
					URL:  "https://pub-7b5607a210bc4f0b81cb6ba41e8754f9.r2.dev/test/IMG_1243.jpg",
					Size: 1930321, Metadata: map[string]any{"mime_type": "image/jpeg", "width": 4032, "height": 3024},
				},
				{
					ID:   utils.Must(xid.FromString("00000000000000000040")),
					Name: asset.NewFilename("00000000000000000040-asset-01"),
					URL:  "https://pub-7b5607a210bc4f0b81cb6ba41e8754f9.r2.dev/test/IMG_1264.jpg",
					Size: 1724055, Metadata: map[string]any{"mime_type": "image/jpeg", "width": 4032, "height": 3024},
				},
				{
					ID:   utils.Must(xid.FromString("00000000000000000040")),
					Name: asset.NewFilename("00000000000000000040-asset-01"),
					URL:  "https://pub-7b5607a210bc4f0b81cb6ba41e8754f9.r2.dev/test/4b3f6b4eeadb4dcc9358541c1d377588.mov",
					Size: 2552083, Metadata: map[string]any{"mime_type": "video/quicktime", "width": 1080, "height": 1920},
				},
				{
					ID:   utils.Must(xid.FromString("00000000000000000040")),
					Name: asset.NewFilename("00000000000000000040-asset-01"),
					URL:  "https://pub-7b5607a210bc4f0b81cb6ba41e8754f9.r2.dev/test/1631887536125.mov",
					Size: 3155277, Metadata: map[string]any{"mime_type": "video/quicktime", "width": 1080, "height": 1920},
				},
			},
		},
		Title:    "Trip to Iceland",
		Category: Category_02_Photos,

		Replies: []*reply.Reply{
			{
				Post: post.Post{
					Content: utils.Must(content.NewRichText("")),
				},
			},
		},
	}

	Threads = []thread.Thread{
		Post_01_Welcome,
		Post_02_HowToContribute,
		Post_03_LoremIpsum,
		Post_04_Photos,
	}
)

func threads(tr thread.Repository, pr reply.Repository, rr react.Repository, ar asset.Repository) {
	ctx := context.Background()

	for _, t := range Threads {
		assetIDs := []asset.AssetID{}
		for i, a := range t.Assets {
			id := fmt.Sprintf("%s-asset-%d", t.ID, i)

			a, err := ar.Add(ctx, t.Author.ID, asset.NewFilename(id), a.URL)
			if err != nil {
				panic(err)
			}
			assetIDs = append(assetIDs, a.ID)
		}

		// TODO: Seed from externally via service or API.
		th, err := tr.Create(ctx,
			t.Title,
			t.Author.ID,
			t.Category.ID,
			t.Tags,
			thread.WithID(t.ID),
			thread.WithContent(t.Content),
			thread.WithVisibility(visibility.VisibilityPublished),
			thread.WithAssets(assetIDs),
		)
		if err != nil {
			if ent.IsConstraintError(err) {
				continue
			}
			panic(err)
		}

		for _, p := range t.Replies {
			p, err = pr.Create(ctx,
				p.Author.ID,
				th.ID,
				reply.WithContent(p.Content),
				reply.WithID(p.ID))
			if err != nil {
				if ent.IsConstraintError(err) {
					continue
				}
				panic(err)
			}

			for i := 0; i < rand.Intn(10); i++ {

				acc := Accounts[rand.Intn(len(Accounts))]

				_, err := rr.Add(ctx,
					acc.ID,
					xid.ID(p.ID),
					randomEmoji())
				if err != nil {
					panic(err)
				}
			}
		}
	}

	fmt.Println("created seed threads and posts")
}

// https://gist.github.com/chiliec/60d815bcbfc56ff62fafe2ff8ce80f6b
func randomEmoji() string {
	rand.Seed(time.Now().UnixNano())
	emoji := [][]int{
		{128513, 128591},
		{128640, 128704},
	}
	r := emoji[rand.Int()%len(emoji)]
	min := r[0]
	max := r[1]
	n := rand.Intn(max-min+1) + min
	return html.UnescapeString("&#" + strconv.Itoa(n) + ";")
}
