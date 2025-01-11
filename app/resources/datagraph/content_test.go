package datagraph

import (
	"encoding/json"
	"testing"

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
	a := assert.New(t)

	c, err := NewRichText(`<h1>heading</h1>

<p>Here's a paragraph. It's pretty neat.</p>

<p>Here's the rest of the text.</p>

<img src="http://image.com" />

<p>neat photo right?</p>

<p>This is quite a long post, the summary, should just be the first 128 characters rounded down to the nearest space.</p>`)
	r.NoError(err)
	r.NotNil(r)

	ps := c.Split()
	a.Len(ps, 5)
}

func TestSplitMinimal(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	c, err := NewRichText(`I&#39;ve tried everything, but for some reason it seems impossible to find datasets that includes a simple list of the councils in England sorted by their group, and a list of covid cases also sorted by councils.  I&#39;m not British so it may be a lack of knowledge of how their government sites work. 

Anyone know a place to find these?`)
	r.NoError(err)
	r.NotNil(r)

	ps := c.Split()
	a.Len(ps, 1)
}

func TestSplitLong(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	c, err := NewRichText(`<body>
<p>A very short paragraph.</p>

<p>A very long single paragraph parato Solis gemitus nefandam munus cupidisque luminis Fuge fuit vestra undis laudando in aristas Lorem markdownum trepidantum genetricis [late](http://hic.com/nunc), steterat delendaque summique non? Domum si cursus supremaque aeraque [manibus](http://www.hoc-quod.io/); Iovis illi fert Bacchum, vulneraque __Oleniae suis__ increpat. Sanguine raucis albet Martius infantemque est _parili multam_ auditaque Caencu inferior augent, vix dote telae volat nec horto sceleratus. Mirator ambage! _Tale_ hic Diomedis, arva sonum factas maxima et relicto Longa incustoditam fixus. Praestantissima operata ardere semine per formae quod: in accipe quamvis amoribus, aquis medio confido puer stridore et clamavit? Video invidiosa Glauci flaventibus funeribus _tenet viscera boumque_ tacuit mearum unda interea calorem poscunt primum. Materia et Ethemon attonitae pactaeque in auctor, furiosior miscuerat? Anguem Arctonque receptus ait vacuus vestigia vapore praetulit, moves? var compression = array + pretestVrmlLion.mediaNetbios(file(topologyFlash, 4), 625909); fat_encoding(4, 1); honeypot += cisc; tigerOutboxFile -= repository_memory_symbolic(skin_dlc(5, sampling_excel)); var floppySmtp = core(5, fiber_opengl_touchscreen, gigabitProxy + rdf_cycle_scrolling); In orare ut finitimi Cum essemus spisso inferias: quam post: __unus cum tendebat__, ad, ecquem. Putaret vulnera noxque de vadorum materiam nomen, nives inroravere equidem Hippason hic convivia: per. var reimage_pup_case = saas + 2; cache = delExecutableMotion.subdirectory_service_status(-3, nosqlBin, analog_thick_rte) - adf; var ddlKerning = sector / meme + 3; Inania classes te lavere; _decimo iter prominet_ Scirone, contra harena! Neu multaque vocatur raucum dux pia, fruges illam Cupidinis huius, corruit aurea crudelius structis. Velamina manebit, unda desint saetas, deteriora domumque: e haec Ceres. Sic duro qua quidque victi Ino non, valentes fuit. Ventisque habeto [aliter feto vinci](http://partibus.org/)! Tonitrumque fidas quaerensque proles temptatum citharae iuguli duae patrio coniugis est mea genus dominae, ut nisi [ignoro](http://letiferos.net/auratis-humili.php)? Aetatis lactantia; ad et per clivo cognovit pretium. Adorat Solem illa flumine nobis patriam auxiliare illa Theseus dubioque lunae nactus discedere obiecit, e. Where do I begin? Easy. There are lots of problems with current forms of government. Clearly Democracy is not working (see [https:&#x2F;&#x2F;www.youtube.com&#x2F;watch?v=QFgcqB8-AxE] for a pretty short summary of why). Based on the IQ distribution, we have that 50% of people are below 100 IQ. If you go to a university or other decently high skill job, think about the dumbest person at that job. They probably have a 110+ IQ. Now consider that over 50% of people are dumber than that idiot who you hate with all your being. He is probably your boss, that project manager who wants you to play hopscotch with strangers as a team building exercise. Imagine if people just like him and even more stupid could singlehandedly decide on who the next ruler of your country will be. Terrible.So clearly the next best option is to have the decision be made by a more intelligent group of people. Who better than Rabbis? They have studied all their lives, are genetically more likely to have a very high IQ, have shown immense dedication, work ethic, and pure intentions (aside from pricking the penis of male converts, not sure why they do that). It&#x27;s common for them to engage in debates and intellectual discussions with each other, and they are chosen by G-d as His favored people to lead the way forward for humanity.Imagine a society where they are able to choose amongst themselves. Personally I think it would be amazing. The person they choose doesn&#x27;t even have to be a Rabbi or Jewish at all, it could be some random kid. But we need to all trust in their judgment because it is the best one available to us. To keep things fresh it&#x27;s probably best to rotate different Rabbis every year, maybe have one year be Conservative, the next one be Reform, etc. just for the variety and to give them a break. Many of them are senior citizens, we don&#x27;t want them getting exhausted or accelerating neurological issues they might have.</p>
</body>
`)
	r.NoError(err)
	r.NotNil(r)

	ps := c.Split()
	r.Len(ps, 15)
	a.Equal("A very short paragraph.", ps[0])
	a.Equal("A very long single paragraph parato Solis gemitus nefandam munus cupidisque luminis Fuge fuit vestra undis laudando in aristas Lorem markdownum trepidantum genetricis [late](http://hic.com/nunc), steterat delendaque summique non? Domum si cursus supremaque aeraque [manibus](http://www.hoc-quod.io/)", ps[1])
	a.Equal("Iovis illi fert Bacchum, vulneraque __Oleniae suis__ increpat. Sanguine raucis albet Martius infantemque est _parili multam_ auditaque Caencu inferior augent, vix dote telae volat nec horto sceleratus. Mirator ambage! _Tale_ hic Diomedis, arva sonum factas maxima et relicto Longa incustoditam fixus", ps[2])
	a.Equal("Praestantissima operata ardere semine per formae quod: in accipe quamvis amoribus, aquis medio confido puer stridore et clamavit? Video invidiosa Glauci flaventibus funeribus _tenet viscera boumque_ tacuit mearum unda interea calorem poscunt primum. Materia et Ethemon attonitae pactaeque in auctor, furiosior miscuerat", ps[3])
	a.Equal("Anguem Arctonque receptus ait vacuus vestigia vapore praetulit, moves? var compression = array + pretestVrmlLion.mediaNetbios(file(topologyFlash, 4), 625909); fat_encoding(4, 1); honeypot += cisc; tigerOutboxFile -= repository_memory_symbolic(skin_dlc(5, sampling_excel))", ps[4])
	a.Equal("var floppySmtp = core(5, fiber_opengl_touchscreen, gigabitProxy + rdf_cycle_scrolling); In orare ut finitimi Cum essemus spisso inferias: quam post: __unus cum tendebat__, ad, ecquem. Putaret vulnera noxque de vadorum materiam nomen, nives inroravere equidem Hippason hic convivia: per. var reimage_pup_case = saas + 2; cache = delExecutableMotion", ps[5])
	a.Equal("subdirectory_service_status(-3, nosqlBin, analog_thick_rte) - adf; var ddlKerning = sector / meme + 3; Inania classes te lavere; _decimo iter prominet_ Scirone, contra harena! Neu multaque vocatur raucum dux pia, fruges illam Cupidinis huius, corruit aurea crudelius structis. Velamina manebit, unda desint saetas, deteriora domumque: e haec Ceres", ps[6])
	a.Equal("Sic duro qua quidque victi Ino non, valentes fuit. Ventisque habeto [aliter feto vinci](http://partibus.org/)! Tonitrumque fidas quaerensque proles temptatum citharae iuguli duae patrio coniugis est mea genus dominae, ut nisi [ignoro](http://letiferos.net/auratis-humili.php)? Aetatis lactantia; ad et per clivo cognovit pretium", ps[7])
	a.Equal("Adorat Solem illa flumine nobis patriam auxiliare illa Theseus dubioque lunae nactus discedere obiecit, e. Where do I begin? Easy. There are lots of problems with current forms of government. Clearly Democracy is not working (see [https://www.youtube.com/watch?v=QFgcqB8-AxE] for a pretty short summary of why)", ps[8])
	a.Equal("Based on the IQ distribution, we have that 50% of people are below 100 IQ. If you go to a university or other decently high skill job, think about the dumbest person at that job. They probably have a 110+ IQ. Now consider that over 50% of people are dumber than that idiot who you hate with all your being", ps[9])
	a.Equal("He is probably your boss, that project manager who wants you to play hopscotch with strangers as a team building exercise. Imagine if people just like him and even more stupid could singlehandedly decide on who the next ruler of your country will be. Terrible", ps[10])
	a.Equal("So clearly the next best option is to have the decision be made by a more intelligent group of people. Who better than Rabbis? They have studied all their lives, are genetically more likely to have a very high IQ, have shown immense dedication, work ethic, and pure intentions (aside from pricking the penis of male converts, not sure why they do", ps[11])
	a.Equal("that). It's common for them to engage in debates and intellectual discussions with each other, and they are chosen by G-d as His favored people to lead the way forward for humanity.Imagine a society where they are able to choose amongst themselves. Personally I think it would be amazing", ps[12])
	a.Equal("The person they choose doesn't even have to be a Rabbi or Jewish at all, it could be some random kid. But we need to all trust in their judgment because it is the best one available to us. To keep things fresh it's probably best to rotate different Rabbis every year, maybe have one year be Conservative, the next one be Reform, etc", ps[13])
	a.Equal("just for the variety and to give them a break. Many of them are senior citizens, we don't want them getting exhausted or accelerating neurological issues they might have.", ps[14])
}

func TestSplitLong2(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	c, err := NewRichText(`<body>Hello friends,

I need help. I am currently using Europeam Values Study data set (2017) and i did crosstab for two variables - country code and political party support.

The problem is that i have been given all the countries and all the political parties.

I would like to sort varibles in a way that i see only a specific country and the support for the political parties only in that country.

Thank you im advance</body>
`)
	r.NoError(err)
	r.NotNil(r)

	ps := c.Split()
	r.Len(ps, 2)
	a.Equal("Hello friends,\n\nI need help. I am currently using Europeam Values Study data set (2017) and i did crosstab for two variables - country code and political party support.\n\nThe problem is that i have been given all the countries and all the political parties", ps[0])
	a.Equal("I would like to sort varibles in a way that i see only a specific country and the support for the political parties only in that country.\n\nThank you im advance", ps[1])
}
