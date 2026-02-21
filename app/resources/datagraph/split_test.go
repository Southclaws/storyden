package datagraph

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/Southclaws/storyden/internal/utils"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func check(t *testing.T, want Content) func(got Content, err error) {
	return func(got Content, err error) {
		require.NoError(t, err)
		assert.Equal(t, want.short, got.short)
		assert.Equal(t, want.links, got.links)
		assert.Equal(t, want.media, got.media)
	}
}

func TestNewRichText(t *testing.T) {
	t.Run("simple_html", func(t *testing.T) {
		check(t, Content{
			short: `Here's a paragraph. It's pretty neat. Here's the rest of the text. neat photo right? This is quite a long post, the summary...`,
			links: []string{},
			media: []string{"http://image.com"},
		})(NewRichText(`<h1>heading</h1>

<p>Here's a paragraph. It's pretty neat.</p>

<p>Here's the rest of the text.</p>

<img src="http://image.com" />

<p>neat photo right?</p>

<p>This is quite a long post, the summary, should just be the first 128 characters rounded down to the nearest space.</p>`))
	})

	t.Run("pull_links", func(t *testing.T) {
		check(t, Content{
			short: `Here's a paragraph. It's pretty neat. here are my favourite ovens here are my favourite trees`,
			links: []string{"https://ao.com/cooking/ovens", "https://tre.ee/trees/favs"},
			media: []string{},
		})(NewRichText(`<h1>heading</h1>

<p>Here's a paragraph. It's pretty neat.</p>

<a href="https://ao.com/cooking/ovens">here are my favourite ovens</a>
<a href="https://tre.ee/trees/favs">here are my favourite trees</a>
`))
	})

	t.Run("pull_images", func(t *testing.T) {
		check(t, Content{
			short: `Here are some cool photos.`,
			links: []string{},
			media: []string{
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fcarters-halt.jpg&w=3840&q=75",
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2F30.jpg&w=3840&q=75",
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fboxes.jpg&w=2048&q=75",
			},
		})(NewRichText(`<h1>heading</h1>

<p>Here are some cool photos.</p>

<img src="https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fcarters-halt.jpg&w=3840&q=75" />
<img src="https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2F30.jpg&w=3840&q=75" />
<img src="https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fboxes.jpg&w=2048&q=75" />
`))
	})

	t.Run("pull_images_relative", func(t *testing.T) {
		check(t, Content{
			short: `Here are some cool photos.`,
			links: []string{},
			media: []string{
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fcarters-halt.jpg&w=3840&q=75",
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2F30.jpg&w=3840&q=75",
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fboxes.jpg&w=2048&q=75",
			},
		})(NewRichTextWithOptions(`<h1>heading</h1>

<p>Here are some cool photos.</p>

<img src="/_next/image?url=%2Fphotography%2Fcity-of-london%2Fcarters-halt.jpg&w=3840&q=75" />
<img src="/_next/image?url=%2Fphotography%2Fcity-of-london%2F30.jpg&w=3840&q=75" />
<img src="/_next/image?url=%2Fphotography%2Fcity-of-london%2Fboxes.jpg&w=2048&q=75" />
`, WithBaseURL("https://barney.is")))
	})

	t.Run("with_uris", func(t *testing.T) {
		mention := utils.Must(xid.FromString("cn2h3gfljatbqvjqctdg"))

		check(t, Content{
			short: `hey @southclaws!`,
			links: []string{},
			media: []string{},
			sdrs: RefList{
				{Kind: KindProfile, ID: mention},
			},
		})(NewRichText(`<h1>heading</h1><p>hey <a href="sdr:profile/cn2h3gfljatbqvjqctdg">@southclaws</a>!</p>`))
	})

	t.Run("json", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		original, err := NewRichText(`<body><p>a</p></body>`)
		r.NoError(err)
		r.NotEmpty(original)

		encoded, err := json.Marshal(original)
		r.NoError(err)
		r.NotEmpty(encoded)

		a.Equal(`"\u003cbody\u003e\u003cp\u003ea\u003c/p\u003e\u003c/body\u003e"`, string(encoded))

		var parsed Content
		err = json.Unmarshal(encoded, &parsed)
		r.NoError(err)
		r.NotEmpty(parsed)

		a.Equal(original, parsed)
	})
}

func TestNewRichTextFromMarkdown(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		fmd, err := NewRichTextFromMarkdown(`To start with data science, it is essential to begin with a practical and step-by-step approach. First, explore platforms like Kaggle, which provide datasets that are accessible for beginners. While Kaggle may appear daunting at first, you can choose beginner-friendly tutorials and datasets that interest you. Begin by downloading and inspecting these datasets to get familiar with their structure and content. Concurrently, work on crafting questions from the data to guide your exploration—this helps in developing a problem-solving mindset.

Consistency in practicing with data, asking for advice, and seeking support, such as shared links or files, are also key steps. Keep in mind that experience and understanding grow steadily through practice rather than seeking perfection right away.

References:
- sdr:thread/cto7n8ifunp55p1bujv0: Emphasized the importance of staying practical and using beginner tutorials and platforms like Kaggle.
- sdr:thread/cto7nm2funp55p1bujvg: Provided advice on starting with data, forming questions, and the value of consistent practice.
`)

		check(t, Content{
			short: `To start with data science, it is essential to begin with a practical and step-by-step approach. First, explore platforms like...`,
			links: []string{},
			media: []string{},
		})(fmd, err)

		rendered := fmd.HTML()
		assert.Equal(t, `<body><p>To start with data science, it is essential to begin with a practical and step-by-step approach. First, explore platforms like Kaggle, which provide datasets that are accessible for beginners. While Kaggle may appear daunting at first, you can choose beginner-friendly tutorials and datasets that interest you. Begin by downloading and inspecting these datasets to get familiar with their structure and content. Concurrently, work on crafting questions from the data to guide your exploration—this helps in developing a problem-solving mindset.</p>

<p>Consistency in practicing with data, asking for advice, and seeking support, such as shared links or files, are also key steps. Keep in mind that experience and understanding grow steadily through practice rather than seeking perfection right away.</p>

<p>References:</p>

<ul>
<li>sdr:thread/cto7n8ifunp55p1bujv0: Emphasized the importance of staying practical and using beginner tutorials and platforms like Kaggle.</li>
<li>sdr:thread/cto7nm2funp55p1bujvg: Provided advice on starting with data, forming questions, and the value of consistent practice.</li>
</ul>
</body>`, rendered)
	})
}

func TestSplit(t *testing.T) {
	r := require.New(t)

	c, err := NewRichText(`<h1>heading</h1>

<p>Here's a paragraph. It's pretty neat.</p>

<p>Here's the rest of the text.</p>

<img src="http://image.com" />

<p>neat photo right?</p>

<p>This is quite a long post, the summary, should just be the first 128 characters rounded down to the nearest space.</p>`)
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	r.Len(ps, 4)
	r.Contains(ps[0], "heading")
	r.Contains(strings.Join(ps, "\n"), "neat photo right?")
}

func TestSplitLongIsStableWithoutExactBoundaries(t *testing.T) {
	r := require.New(t)

	long := strings.Repeat("Jira tickets are multiplying quickly. The sprint plan needs sharper priorities and cleaner ownership. ", 80)
	c, err := NewRichText("<body><p>Short intro.</p><p>" + long + "</p></body>")
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	r.Greater(len(ps), 2)
	for _, p := range ps {
		r.NotEmpty(strings.TrimSpace(p))
		r.LessOrEqual(len([]rune(p)), roughMaxSentenceSize)
	}
	out := strings.Join(ps, " ")
	r.Contains(out, "Short intro.")
	r.Contains(out, "Jira tickets are multiplying quickly.")
}

func TestSplitPreservesSpacesAcrossInlineNodes(t *testing.T) {
	r := require.New(t)

	c, err := NewRichText(`<body><p>Hello <b>world</b>!</p></body>`)
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	r.Len(ps, 1)
	r.Equal("Hello world!", ps[0])
}

func TestSplitNormalizesWeirdWhitespace(t *testing.T) {
	r := require.New(t)

	c, err := NewRichText(`<body><p>Hello   <i>world</i>   !</p></body>`)
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	r.Len(ps, 1)
	r.Equal("Hello world!", ps[0])
}

func TestSplitRuneSafe(t *testing.T) {
	r := require.New(t)

	long := strings.Repeat("café naïve — 漢字。", 200)
	c, err := NewRichText("<body><p>" + long + "</p></body>")
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	r.Greater(len(ps), 1)
	for i, s := range ps {
		if !utf8.ValidString(s) {
			t.Fatalf("chunk %d is not valid UTF-8", i)
		}
	}
}

func TestSplitHandlesListsAsSteps(t *testing.T) {
	r := require.New(t)

	c, err := NewRichText(`<body>
<h2>Reset password</h2>
<ol>
  <li>Open settings.</li>
  <li>Click "Security".</li>
  <li>Choose <code>Reset</code>.</li>
</ol>
</body>`)
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	out := strings.Join(ps, "\n")
	r.Contains(out, "Reset password")
	r.Contains(out, "Open settings.")
	r.Contains(out, `Click "Security".`)
	r.Contains(out, "Choose Reset.")
}

func TestSplitHeadingCarriesFollowingParagraph(t *testing.T) {
	r := require.New(t)

	c, err := NewRichText(`<body>
<h2>Rate limiting</h2>
<p><code>RATE_LIMIT_PERIOD</code> controls the window size.</p>
<p>Only applies when cache is enabled.</p>
</body>`)
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	joined := strings.Join(ps, "\n---\n")
	r.Contains(joined, "Rate limiting")
	r.Contains(joined, "RATE_LIMIT_PERIOD")
}

func TestSplitPreservesPreNewlines(t *testing.T) {
	r := require.New(t)

	c, err := NewRichText(`<body>
<pre>export RATE_LIMIT=1000
export RATE_LIMIT_PERIOD=60</pre>
</body>`)
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	r.Len(ps, 1)
	r.Equal("export RATE_LIMIT=1000\nexport RATE_LIMIT_PERIOD=60", ps[0])
}

func TestSplitHandlesTablesWithRowBoundaries(t *testing.T) {
	r := require.New(t)

	c, err := NewRichText(`<body>
<table>
  <tr><th>Key</th><th>Meaning</th></tr>
  <tr><td>RATE_LIMIT</td><td>Requests per period</td></tr>
  <tr><td>RATE_LIMIT_PERIOD</td><td>Window length</td></tr>
</table>
</body>`)
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	out := strings.Join(ps, "\n")
	r.Contains(out, "RATE_LIMIT")
	r.Contains(out, "RATE_LIMIT_PERIOD")
	r.Contains(out, "|")
}

func TestSplitIgnoresNavChromeNoise(t *testing.T) {
	r := require.New(t)

	c, err := NewRichText(`<body>
<nav><a href="/">Home</a> <a href="/docs">Docs</a> <a href="/login">Login</a></nav>
<main><h1>Docs</h1><p>How to deploy.</p></main>
<footer>© 2026</footer>
</body>`)
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	out := strings.Join(ps, "\n---\n")
	r.Contains(out, "How to deploy.")
	r.NotContains(out, "Login")
	r.NotContains(out, "©")
}

func TestSplitDoesNotDuplicateText(t *testing.T) {
	r := require.New(t)

	c, err := NewRichText(`<body><p>Hello <b>bold</b> world.</p></body>`)
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	r.Len(ps, 1)
	r.Equal(1, strings.Count(ps[0], "bold"))
}

func TestSplitMultilingualInlineSpacing(t *testing.T) {
	r := require.New(t)

	cases := []struct {
		name string
		html string
		want string
	}{
		{name: "chinese", html: `<body><p>机器<b>学习</b>入门。</p></body>`, want: "机器学习入门。"},
		{name: "russian", html: `<body><p>Изучение <b>Python</b> сегодня.</p></body>`, want: "Изучение Python сегодня."},
		{name: "arabic", html: `<body><p>تطوير <b>تطبيقات</b> الويب الحديثة.</p></body>`, want: "تطوير تطبيقات الويب الحديثة."},
		{name: "spanish", html: `<body><p>Recetas <b>mediterráneas</b> tradicionales.</p></body>`, want: "Recetas mediterráneas tradicionales."},
		{name: "french", html: `<body><p>Histoire de l'<b>architecture</b> gothique.</p></body>`, want: "Histoire de l'architecture gothique."},
		{name: "german", html: `<body><p>Wandern in den <b>Alpen</b> heute.</p></body>`, want: "Wandern in den Alpen heute."},
		{name: "greek", html: `<body><p>Αρχαία <b>Ελληνική</b> φιλοσοφία.</p></body>`, want: "Αρχαία Ελληνική φιλοσοφία."},
		{name: "turkish", html: `<body><p>Geleneksel <b>Türk</b> mutfağı.</p></body>`, want: "Geleneksel Türk mutfağı."},
		{name: "georgian", html: `<body><p>ქართული <b>ხალხური</b> სიმღერები.</p></body>`, want: "ქართული ხალხური სიმღერები."},
		{name: "hindi", html: `<body><p>भारतीय <b>शास्त्रीय</b> संगीत.</p></body>`, want: "भारतीय शास्त्रीय संगीत."},
		{name: "hebrew", html: `<body><p>ספרות עברית <b>מודרנית</b>.</p></body>`, want: "ספרות עברית מודרנית."},
		{name: "persian", html: `<body><p>شعر <b>کلاسیک</b> فارسی.</p></body>`, want: "شعر کلاسیک فارسی."},
		{name: "urdu", html: `<body><p>اردو <b>شاعری</b> کی روایت.</p></body>`, want: "اردو شاعری کی روایت."},
		{name: "punjabi", html: `<body><p>ਪੰਜਾਬੀ <b>ਲੋਕ</b> ਗੀਤ.</p></body>`, want: "ਪੰਜਾਬੀ ਲੋਕ ਗੀਤ."},
		{name: "nepali", html: `<body><p>नेपाली <b>पर्वतारोहण</b> इतिहास.</p></body>`, want: "नेपाली पर्वतारोहण इतिहास."},
		{name: "yoruba", html: `<body><p>Àṣà <b>Yorùbá</b> lónìí.</p></body>`, want: "Àṣà Yorùbá lónìí."},
		{name: "igbo", html: `<body><p>Omenala <b>Igbo</b> taa.</p></body>`, want: "Omenala Igbo taa."},
		{name: "hausa", html: `<body><p>Tarihin <b>Hausawa</b> a Arewa.</p></body>`, want: "Tarihin Hausawa a Arewa."},
		{name: "akan", html: `<body><p>Akanfo <b>Atetesɛm</b> ne amammerɛ.</p></body>`, want: "Akanfo Atetesɛm ne amammerɛ."},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewRichText(tc.html)
			r.NoError(err)

			ps := c.Split()
			fmt.Println(ps)
			r.Len(ps, 1)
			r.Equal(tc.want, ps[0])
		})
	}
}

func TestNeedsSpaceCJKRegression(t *testing.T) {
	r := require.New(t)

	c, err := NewRichText(`<body><p>深度<b>学习</b>与神经网络。</p></body>`)
	r.NoError(err)

	ps := c.Split()
	fmt.Println(ps)
	r.Len(ps, 1)
	r.Equal("深度学习与神经网络。", ps[0])
}
