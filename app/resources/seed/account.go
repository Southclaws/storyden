package seed

import (
	"context"
	"fmt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/authentication"
)

const SeedPassword = `$argon2id$v=19$m=65536,t=1,p=2$MAwllQoeGcxCPOC52OQwZA$jLlzHsmSHmQPbpQ6Y5+877NlacOYeyqEqWoKJJXRcHM`

var (
	Account_000 = account.Account{ID: account.AccountID(id("00000000000000000010")), Name: "Francis J. Underwood", Handle: "francis-j-underwood", Admin: true}
	Account_001 = account.Account{ID: account.AccountID(id("00000000000000000020")), Name: "Claire Hale Underwood", Handle: "claire-hale-underwood", Admin: true}
	Account_002 = account.Account{ID: account.AccountID(id("00000000000000000030")), Name: "Zoe Barnes", Handle: "zoe-barnes"}
	Account_003 = account.Account{ID: account.AccountID(id("00000000000000000040")), Name: "Peter Russo", Handle: "peter-russo"}
	Account_004 = account.Account{ID: account.AccountID(id("00000000000000000050")), Name: "Doug Stamper", Handle: "doug-stamper"}
	Account_005 = account.Account{ID: account.AccountID(id("00000000000000000060")), Name: "Christina Gallagher", Handle: "christina-gallagher"}
	Account_006 = account.Account{ID: account.AccountID(id("00000000000000000070")), Name: "Linda Vasquez", Handle: "linda-vasquez"}
	Account_007 = account.Account{ID: account.AccountID(id("00000000000000000080")), Name: "Gillian Cole", Handle: "gillian-cole"}
	Account_008 = account.Account{ID: account.AccountID(id("00000000000000000090")), Name: "Janine Skorsky", Handle: "janine-skorsky"}
	Account_009 = account.Account{ID: account.AccountID(id("00000000000000000100")), Name: "Garrett Walker", Handle: "garrett-walker"}
	Account_010 = account.Account{ID: account.AccountID(id("00000000000000000110")), Name: "Lucas Goodwin", Handle: "lucas-goodwin"}
	Account_011 = account.Account{ID: account.AccountID(id("00000000000000000120")), Name: "Remy Danton", Handle: "remy-danton"}
	Account_012 = account.Account{ID: account.AccountID(id("00000000000000000130")), Name: "Tom Hammerschmidt", Handle: "tom-hammerschmidt"}
	Account_013 = account.Account{ID: account.AccountID(id("00000000000000000140")), Name: "Edward Meechum", Handle: "edward-meechum"}
	Account_014 = account.Account{ID: account.AccountID(id("00000000000000000150")), Name: "Rachel Posner", Handle: "rachel-posner"}
	Account_015 = account.Account{ID: account.AccountID(id("00000000000000000160")), Name: "Raymond Tusk", Handle: "raymond-tusk"}
	Account_016 = account.Account{ID: account.AccountID(id("00000000000000000170")), Name: "Cathy Durant", Handle: "cathy-durant"}
	Account_017 = account.Account{ID: account.AccountID(id("00000000000000000180")), Name: "Jackie Sharp", Handle: "jackie-sharp"}
	Account_018 = account.Account{ID: account.AccountID(id("00000000000000000190")), Name: "Gavin Orsay", Handle: "gavin-orsay"}
	Account_019 = account.Account{ID: account.AccountID(id("00000000000000000200")), Name: "Ayla Sayyad", Handle: "ayla-sayyad"}
	Account_020 = account.Account{ID: account.AccountID(id("00000000000000000210")), Name: "Seth Grayson", Handle: "seth-grayson"}
	Account_021 = account.Account{ID: account.AccountID(id("00000000000000000220")), Name: "Heather Dunbar", Handle: "heather-dunbar"}
	Account_022 = account.Account{ID: account.AccountID(id("00000000000000000230")), Name: "Thomas Yates", Handle: "thomas-yates"}
	Account_023 = account.Account{ID: account.AccountID(id("00000000000000000240")), Name: "Viktor Petrov", Handle: "viktor-petrov"}
	Account_024 = account.Account{ID: account.AccountID(id("00000000000000000250")), Name: "Kate Baldwin", Handle: "kate-baldwin"}
	Account_025 = account.Account{ID: account.AccountID(id("00000000000000000260")), Name: "LeAnn Harvey", Handle: "le-ann-harvey"}
	Account_026 = account.Account{ID: account.AccountID(id("00000000000000000270")), Name: "Will Conway", Handle: "will-conway"}
	Account_027 = account.Account{ID: account.AccountID(id("00000000000000000280")), Name: "Mark Usher", Handle: "mark-usher"}
	Account_028 = account.Account{ID: account.AccountID(id("00000000000000000290")), Name: "Jane Davis", Handle: "jane-davis"}
	Account_029 = account.Account{ID: account.AccountID(id("00000000000000000300")), Name: "Hannah Conway", Handle: "hannah-conway"}
	Account_030 = account.Account{ID: account.AccountID(id("00000000000000000310")), Name: "Aidan Macallan", Handle: "aidan-macallan"}
	Account_031 = account.Account{ID: account.AccountID(id("00000000000000000320")), Name: "Annette Shepherd", Handle: "annette-shepherd"}
	Account_032 = account.Account{ID: account.AccountID(id("00000000000000000330")), Name: "Bill Shepherd", Handle: "bill-shepherd"}
	Account_033 = account.Account{ID: account.AccountID(id("00000000000000000340")), Name: "Duncan Shepherd", Handle: "duncan-shepherd"}
)

func accounts(r account.Repository, auth authentication.Repository) {
	ctx := context.Background()

	for _, v := range []account.Account{
		Account_000,
		Account_001,
		Account_002,
		Account_003,
		Account_004,
		Account_005,
		Account_006,
		Account_007,
		Account_008,
		Account_009,
		Account_010,
		Account_011,
		Account_012,
		Account_013,
		Account_014,
		Account_015,
		Account_016,
		Account_017,
		Account_018,
		Account_019,
		Account_020,
		Account_021,
		Account_022,
		Account_023,
		Account_024,
		Account_025,
		Account_026,
		Account_027,
		Account_028,
		Account_029,
		Account_030,
		Account_031,
		Account_032,
		Account_033,
	} {
		acc, err := r.Create(ctx, v.Handle,
			account.WithID(v.ID),
			account.WithName(v.Name),
			account.WithBio(v.Bio.ElseZero()),
		)
		if err != nil {
			panic(err)
		}

		if _, err = auth.Create(ctx, acc.ID, authentication.Service("password"), acc.Handle+"@storyd.en", SeedPassword, nil); err != nil {
			panic(err)
		}

	}

	fmt.Println("created seed users")
}
