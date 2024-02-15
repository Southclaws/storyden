package seed

import (
	"context"
	"fmt"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/cluster"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/item"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
)

var Cluster_Ancient = datagraph.Cluster{
	ID:          datagraph.ClusterID(id("00000000000010000000")),
	Owner:       profile.Profile{ID: Account_001_Odin.ID},
	Name:        "Ancient Ships",
	Description: "Ancient ships are ships from the ancient world, which are the ships that were in use from the Bronze Age to the Middle Ages. This includes all vessels those were sail or row boats, powered by oars or sails. The ships of antiquity are mostly classed as ancient warships, trading vessels and sometimes fishing boats.",
	Content:     opt.New("This cluster contains ships from the ancient world, which are the ships that were in use from the Bronze Age to the Middle Ages. This includes all vessels those were sail or row boats, powered by oars or sails. The ships of antiquity are mostly classed as ancient warships, trading vessels and sometimes fishing boats."),
	Assets:      []*asset.Asset{{URL: "https://ih1.redbubble.net/image.499157992.1766/flat,1000x1000,075,f.u1.jpg"}},
	Items: []*datagraph.Item{
		{
			ID:          datagraph.ItemID(id("00000000000010000010")),
			Owner:       profile.Profile{ID: Account_001_Odin.ID},
			Name:        "Hringhorni",
			Assets:      []*asset.Asset{{URL: "https://upload.wikimedia.org/wikipedia/commons/1/18/Thor_kicks_Litr.jpg"}},
			Description: `In Norse mythology, Hringhorni (Old Norse "ship with a circle on the stem") is the name of the ship of the god Baldr, described as the "greatest of all ships".`,
			Content: opt.New(`According to Gylfaginning, following the murder of Baldr by Loki, the other gods brought his body down to the sea and laid him to rest on the ship. They would have launched it out into the water and kindled a funeral pyre for Baldr but were unable to move the great vessel without the help of the giantess Hyrrokkin, who was sent for out of Jötunheim. She then flung the ship so violently down the rollers at the first push that flames appeared and the earth trembled, much to the annoyance of Thor.
Along with Baldr, his wife Nanna was also borne to the funeral pyre after she had died of grief. As Thor was consecrating the fire with his hammer Mjolnir, a dwarf named Litr began cavorting at his feet. Thor then kicked him into the flames and the dwarf was burned up as well. The significance of this seemingly incidental event is speculative but may perhaps find a parallel in religious ritual. Among other artifacts and creatures sacrificed on the pyre of Hringhorni were Odin's gold ring Draupnir and the horse of Baldr with all its trappings.
`),
		},
		{
			ID:          datagraph.ItemID(id("00000000000010000020")),
			Owner:       profile.Profile{ID: Account_001_Odin.ID},
			Name:        "Naglfar",
			Assets:      []*asset.Asset{{URL: "https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/f/a842aadd-58e9-47a3-90f5-5f5b2f1a2722/djnyin-f640fb25-7e1d-4e6b-9abe-5a8a8138fba9.jpg/v1/fit/w_800,h_1232,q_70,strp/naglfar_full_by_faile35_djnyin-414w-2x.jpg?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1cm46YXBwOjdlMGQxODg5ODIyNjQzNzNhNWYwZDQxNWVhMGQyNmUwIiwiaXNzIjoidXJuOmFwcDo3ZTBkMTg4OTgyMjY0MzczYTVmMGQ0MTVlYTBkMjZlMCIsIm9iaiI6W1t7ImhlaWdodCI6Ijw9MTIzMiIsInBhdGgiOiJcL2ZcL2E4NDJhYWRkLTU4ZTktNDdhMy05MGY1LTVmNWIyZjFhMjcyMlwvZGpueWluLWY2NDBmYjI1LTdlMWQtNGU2Yi05YWJlLTVhOGE4MTM4ZmJhOS5qcGciLCJ3aWR0aCI6Ijw9ODAwIn1dXSwiYXVkIjpbInVybjpzZXJ2aWNlOmltYWdlLm9wZXJhdGlvbnMiXX0.Rc_LIxWvVB7HBk0FULN0Ew5p5g4HNhJVKSvMHgtp_L0"}},
			Description: `In Norse mythology, Naglfar or Naglfari (Old Norse "nail farer") is a boat made entirely from the fingernails and toenails of the dead.`,
			Content:     opt.New(`During the events of Ragnarök, Naglfar is foretold to sail to Vígríðr, ferrying hordes of monsters that will do battle with the gods. Naglfar is attested in the Poetic Edda, compiled in the 13th century from earlier traditional sources, and the Prose Edda, also composed in the 13th century. The boat itself has been connected by scholars with a larger pattern of ritual hair and nail disposal among Indo-Europeans, stemming from Proto-Indo-European custom, and it may be depicted on the Tullstorp Runestone in Scania, Sweden.`),
		},
		{
			ID:          datagraph.ItemID(id("00000000000010000030")),
			Owner:       profile.Profile{ID: Account_001_Odin.ID},
			Name:        "Skíðblaðnir",
			Assets:      []*asset.Asset{{URL: "https://upload.wikimedia.org/wikipedia/commons/4/4c/The_third_gift_%E2%80%94_an_enormous_hammer_by_Elmer_Boyd_Smith.jpg"}},
			Description: `Skíðblaðnir (Old Norse: [ˈskiːðˌblɑðnez̠], 'assembled from thin pieces of wood'), sometimes anglicized as Skidbladnir or Skithblathnir, is the best of ships in Norse mythology.`,
			Content:     opt.New(`It is attested in the Poetic Edda, compiled in the 13th century from earlier traditional sources, and in the Prose Edda and Heimskringla, both written in the 13th century by Snorri Sturluson. All sources note that the ship is the finest of ships, and the Poetic Edda and Prose Edda attest that it is owned by the god Freyr, while the euhemerized account in Heimskringla attributes it to the magic of Odin. Both Heimskringla and the Prose Edda attribute to it the ability to be folded up—as cloth may be—into one's pocket when not needed.`),
		},
	},
}

var Cluster_Steam = datagraph.Cluster{
	ID:          datagraph.ClusterID(id("00000000000020000000")),
	Owner:       profile.Profile{ID: Account_001_Odin.ID},
	Name:        "Steam ships",
	Assets:      []*asset.Asset{{URL: "https://media.cnn.com/api/v1/images/stellar/prod/140730133850-belle-of-louisville.jpg?q=w_1884,h_1024,x_0,y_0,c_fill"}},
	Description: `Steam ships are ships that are propelled primarily or entirely by steam engines. The steam was produced in a boiler by burning fossil fuels. The steam drives reciprocating pistons which are connected to the paddle wheels or propellers. Steam ships usually use the prefix designations of PS (for paddle steamer) or SS (for screw steamer). As paddle steamers became less common, SS is assumed by many to stand for screw steamer, rather than steam ship.`,
	Content:     opt.New(`Steam ships are ships that are propelled primarily or entirely by steam engines. The steam was produced in a boiler by burning fossil fuels. The steam drives reciprocating pistons which are connected to the paddle wheels or propellers. Steam ships usually use the prefix designations of PS (for paddle steamer) or SS (for screw steamer). As paddle steamers became less common, SS is assumed by many to stand for screw steamer, rather than steam ship. The term steamboat is used mostly in reference to smaller steam-powered boats working on lakes and rivers, particularly riverboats; steamship generally refers to larger steam-powered ships, usually ocean-going, capable of carrying a (ship's) boat. The term steamwheeler is archaic and rarely used. In England, "steam packet", after its sailing predecessor, was the usual term; even "steam barge" could be used for smaller vessels.`),
	Items: []*datagraph.Item{
		{
			ID:          datagraph.ItemID(id("00000000000020000010")),
			Owner:       profile.Profile{ID: Account_001_Odin.ID},
			Name:        "Belle of Louisville",
			Assets:      []*asset.Asset{{URL: "https://lh3.googleusercontent.com/p/AF1QipOTf1f7dzR1NXixsT9caL8gdSrztqMDNM3uCkIB=s1360-w1360-h1020"}},
			Description: `Belle of Louisville is a steamboat owned and operated by the city of Louisville, Kentucky, and moored at its downtown wharf next to the Riverfront Plaza/Belvedere during its annual operational period.`,
			Content:     opt.New(`Originally named Idlewild, the Belle of Louisville was built by James Rees & Sons Company in Pittsburgh, Pennsylvania, for the West Memphis Packet Company in 1914. She initially operated as a passenger ferry between Memphis, Tennessee, and West Memphis, Arkansas. She also hauled cargo such as cotton, lumber, and grain. She then came to Louisville in 1931 and ran trips between the Fontaine Ferry amusement park near downtown Louisville and Rose Island, a resort about 14 miles (23 km) upriver from Louisville. From 1934 through World War II, Idlewild operated a regular excursion schedule. During this time she was outfitted with special equipment to push oil barges along the river. She also served as a floating USO nightclub for troops stationed at military bases along the Mississippi River.`),
		},
		{
			ID:          datagraph.ItemID(id("00000000000020000020")),
			Owner:       profile.Profile{ID: Account_001_Odin.ID},
			Name:        "Katahdin",
			Assets:      []*asset.Asset{{URL: "https://upload.wikimedia.org/wikipedia/commons/3/33/SSKatahdinII.jpg"}},
			Description: `The Katahdin is a historic steamboat berthed on Moosehead Lake in Greenville, Maine. Built in 1914 at the Bath Iron Works, it at first served the tourist trade on the lake before being converted to a towboat hauling lumber.`,
			Content:     opt.New(`It was fully restored in the 1990s by the nonprofit Moosehead Maritime Museum, and is again giving tours on the lake. One of the very few surviving early lake boats in Maine, and the oldest vessel afloat built at Bath, it was listed on the National Register of Historic Places in 1978.`),
		},
		{
			ID:          datagraph.ItemID(id("00000000000020000030")),
			Owner:       profile.Profile{ID: Account_001_Odin.ID},
			Name:        "Minnehaha",
			Assets:      []*asset.Asset{{URL: "https://upload.wikimedia.org/wikipedia/commons/thumb/3/3d/Steamboat_Minnehaha%2C_1906.jpg/440px-Steamboat_Minnehaha%2C_1906.jpg"}},
			Description: `Minnehaha is a steam-powered excursion vessel on Lake Minnetonka in the U.S. state of Minnesota.`,
			Content:     opt.New(`The vessel was originally in service between 1906 and 1926. After being scuttled in 1926, Minnehaha was raised from the bottom of Lake Minnetonka in 1980, restored, and returned to active service in 1996. The vessel operated uninterrupted on Lake Minnetonka until 2019. It is currently stored in a maintenance facility in the town of Excelsior.`),
		},
	},
}

var Cluster_Pirate = datagraph.Cluster{
	ID:    datagraph.ClusterID(id("00000000000030000000")),
	Owner: profile.Profile{ID: Account_001_Odin.ID},
	Name:  "Pirate Ships",
	Items: []*datagraph.Item{
		{
			ID:          datagraph.ItemID(id("00000000000030000010")),
			Owner:       profile.Profile{ID: Account_001_Odin.ID},
			Name:        "Golden Hind",
			Assets:      []*asset.Asset{{URL: "https://dynamic-media-cdn.tripadvisor.com/media/photo-o/0b/1a/a1/a4/harbour-photo.jpg?w=1200&h=1200&s=1"}},
			Description: `Golden Hind was a galleon captained by Francis Drake in his circumnavigation of the world between 1577 and 1580.`,
			Content:     opt.New(`She was originally known as Pelican, but Drake renamed her mid-voyage in 1578, in honour of his patron, Sir Christopher Hatton, whose crest was a golden hind (a female red deer). Hatton was one of the principal sponsors of Drake's world voyage. A full-sized, seaworthy reconstruction is in London, on the south bank of the Thames.`),
		},
		{
			ID:          datagraph.ItemID(id("00000000000030000020")),
			Owner:       profile.Profile{ID: Account_001_Odin.ID},
			Name:        "Adventure Galley",
			Assets:      []*asset.Asset{{URL: "https://sillyhistory.files.wordpress.com/2015/05/charles_galley_1688.jpg"}},
			Description: `Adventure Galley, also known as Adventure, was an English sailing ship captained by William Kidd, the notorious privateer.`,
			Content:     opt.New(`She was a type of hybrid ship that combined square rigged sails with oars to give her manoeuvrability in both windy and calm conditions. The vessel was launched at the end of 1695 and was acquired by Kidd the following year to serve in his privateering venture. Between April 1696 and April 1698, she travelled thousands of miles across the Atlantic and Indian Oceans in search of pirates but failed to find any until nearly the end of her travels. Instead, Kidd himself turned pirate in desperation at not having obtained any prizes. Adventure Galley succeeded in capturing two vessels off India and brought them back to Madagascar, but by the spring of 1698 the ship's hull had become so rotten and leaky that she was no longer seaworthy.`),
		},
		{
			ID:          datagraph.ItemID(id("00000000000030000030")),
			Owner:       profile.Profile{ID: Account_001_Odin.ID},
			Name:        "Bachelor's Delight",
			Assets:      []*asset.Asset{{URL: "https://2.bp.blogspot.com/-bbmJLxO90dg/T4ReN9TlzzI/AAAAAAAAFYY/R2uAAi26ko0/s1600/ship_1.jpg"}},
			Description: `The Bachelor's Delight was a 36 gun frigate that became the ship of the pirate William Dampier, a British privateer and explorer who mostly liked to meet other cultures rather than practice piracy.`,
		},
		{
			ID:          datagraph.ItemID(id("00000000000030000040")),
			Owner:       profile.Profile{ID: Account_001_Odin.ID},
			Name:        "Fancy",
			Assets:      []*asset.Asset{{URL: "https://upload.wikimedia.org/wikipedia/commons/thumb/5/5c/Henry_Every.gif/300px-Henry_Every.gif"}},
			Description: `The Fancy was the famous pirate ship of Henry Every.`,
			Content:     opt.New(`The ship was originally named the Charles II and while privateering in May of 1694 off the coast of Spain, Every and some crew mutinied and captured the ship. Following its seizure they renamed the ship the Fancy and set off to commit acts of piracy. The Fancy at this time had about fifty cannons and a crew of 150.`),
		},
		{
			ID:          datagraph.ItemID(id("00000000000030000050")),
			Owner:       profile.Profile{ID: Account_001_Odin.ID},
			Name:        "Royal Fortune",
			Assets:      []*asset.Asset{{URL: "https://64.media.tumblr.com/48f9e8deba3920e5795dba48cca4c839/tumblr_ptc60dYLtJ1tmlms3o1_1280.jpg"}},
			Description: `The Royal Fortune was a famous pirate ship of Bartholomew Roberts during the Post Spanish Succession Period.`,
			Content:     opt.New(`If Bartholomew Roberts fathered any children during his adventures on the high seas, he may or may not have named all of them Royal Fortune. In July 1720, Roberts captured a French brigantine off the coast of Newfoundland. He outfitted the naval frigate with 26 cannons, renamed her the Good Fortune and headed south for the Caribbean, where the ship was repaired and renamed the Royal Fortune. Soon after, Roberts captured a French warship operated by the Governor of Martinique, renamed her the Royal Fortune and made the ship his new flagship. Roberts then set sail for West Africa, where he captured the Onslow, renamed her the Royal Fortune, and, well, you know the rest. Roberts died, and the final Royal Fortune sank, on February 10, 1722, in an attack by the British warship HMS Swallow.`),
		},
	},
}

var Cluster_Cruise = datagraph.Cluster{
	ID:    datagraph.ClusterID(id("00000000000040000000")),
	Owner: profile.Profile{ID: Account_001_Odin.ID},
	Name:  "Cruise Ships",
}

var Cluster_Container = datagraph.Cluster{
	ID:    datagraph.ClusterID(id("00000000000050000000")),
	Owner: profile.Profile{ID: Account_001_Odin.ID},
	Name:  "Container Ships",
}

var Cluster_Battle = datagraph.Cluster{
	ID:    datagraph.ClusterID(id("00000000000060000000")),
	Owner: profile.Profile{ID: Account_001_Odin.ID},
	Name:  "Battle Ships",
}

var Clusters = []datagraph.Cluster{
	Cluster_Ancient,
	Cluster_Steam,
	Cluster_Pirate,
	Cluster_Cruise,
	Cluster_Container,
	Cluster_Battle,
}

func clusters_items(cluster_repo cluster.Repository, item_repo item.Repository, ar asset.Repository) {
	ctx := context.Background()

	for _, c := range Clusters {
		s := slug.Make(c.Name)
		_, err := cluster_repo.Create(ctx,
			c.Owner.ID,
			c.Name,
			s,
			c.Description,
			cluster.WithID(c.ID),
			cluster.WithAssets(assets(ar, c.Owner.ID, c.ID.String(), c.Assets)),
			cluster.WithContent(c.Content.String()),
			cluster.WithVisibility(post.VisibilityPublished),
		)
		if err != nil {
			panic(err)
		}

		for _, i := range c.Items {
			s := slug.Make(i.Name)
			_, err := item_repo.Create(ctx,
				i.Owner.ID,
				i.Name,
				s,
				i.Description,
				item.WithID(i.ID),
				item.WithContent(i.Content.String()),
				item.WithAssets(assets(ar, i.Owner.ID, i.ID.String(), i.Assets)),
				item.WithParentClusterAdd(xid.ID(c.ID)),
				item.WithVisibility(post.VisibilityPublished),
			)
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("created seed datagraph")
}

func assets(ar asset.Repository, owner account.AccountID, id string, assets []*asset.Asset) (ids []asset.AssetID) {
	ids, err := dt.MapErr(assets, func(a *asset.Asset) (asset.AssetID, error) {
		a, err := ar.Add(context.Background(), owner, asset.NewExistingFilename(xid.New(), uuid.NewString()), a.URL)
		if err != nil {
			return xid.NilID(), err
		}

		return a.ID, nil
	})
	if err != nil {
		panic(err)
	}

	return ids
}
