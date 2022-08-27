package seed

import (
	"context"
	"fmt"

	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/pkg/resources/account"
)

var (
	Account_000 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000000")), Name: "Francis J. Underwood", Email: "francis-j-underwood@car.ds", Admin: true}
	Account_001 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000001")), Name: "Claire Hale Underwood", Email: "claire-hale-underwood@car.ds", Admin: true}
	Account_002 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000002")), Name: "Zoe Barnes", Email: "zoe-barnes@car.ds"}
	Account_003 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000003")), Name: "Peter Russo", Email: "peter-russo@car.ds"}
	Account_004 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000004")), Name: "Doug Stamper", Email: "doug-stamper@car.ds"}
	Account_005 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000005")), Name: "Christina Gallagher", Email: "christina-gallagher@car.ds"}
	Account_006 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000006")), Name: "Linda Vasquez", Email: "linda-vasquez@car.ds"}
	Account_007 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000007")), Name: "Gillian Cole", Email: "gillian-cole@car.ds"}
	Account_008 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000008")), Name: "Janine Skorsky", Email: "janine-skorsky@car.ds"}
	Account_009 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000009")), Name: "Garrett Walker", Email: "garrett-walker@car.ds"}
	Account_010 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000010")), Name: "Lucas Goodwin", Email: "lucas-goodwin@car.ds"}
	Account_011 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000011")), Name: "Remy Danton", Email: "remy-danton@car.ds"}
	Account_012 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000012")), Name: "Tom Hammerschmidt", Email: "tom-hammerschmidt@car.ds"}
	Account_013 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000013")), Name: "Edward Meechum", Email: "edward-meechum@car.ds"}
	Account_014 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000014")), Name: "Rachel Posner", Email: "rachel-posner@car.ds"}
	Account_015 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000015")), Name: "Raymond Tusk", Email: "raymond-tusk@car.ds"}
	Account_016 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000016")), Name: "Cathy Durant", Email: "cathy-durant@car.ds"}
	Account_017 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000017")), Name: "Jackie Sharp", Email: "jackie-sharp@car.ds"}
	Account_018 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000018")), Name: "Gavin Orsay", Email: "gavin-orsay@car.ds"}
	Account_019 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000019")), Name: "Ayla Sayyad", Email: "ayla-sayyad@car.ds"}
	Account_020 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000020")), Name: "Seth Grayson", Email: "seth-grayson@car.ds"}
	Account_021 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000021")), Name: "Heather Dunbar", Email: "heather-dunbar@car.ds"}
	Account_022 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000022")), Name: "Thomas Yates", Email: "thomas-yates@car.ds"}
	Account_023 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000023")), Name: "Viktor Petrov", Email: "viktor-petrov@car.ds"}
	Account_024 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000024")), Name: "Kate Baldwin", Email: "kate-baldwin@car.ds"}
	Account_025 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000025")), Name: "LeAnn Harvey", Email: "le-ann-harvey@car.ds"}
	Account_026 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000026")), Name: "Will Conway", Email: "will-conway@car.ds"}
	Account_027 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000027")), Name: "Mark Usher", Email: "mark-usher@car.ds"}
	Account_028 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000028")), Name: "Jane Davis", Email: "jane-davis@car.ds"}
	Account_029 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000029")), Name: "Hannah Conway", Email: "hannah-conway@car.ds"}
	Account_030 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000030")), Name: "Aidan Macallan", Email: "aidan-macallan@car.ds"}
	Account_031 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000031")), Name: "Annette Shepherd", Email: "annette-shepherd@car.ds"}
	Account_032 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000032")), Name: "Bill Shepherd", Email: "bill-shepherd@car.ds"}
	Account_033 = account.Account{ID: account.AccountID(id("00000000-0000-4000-8000-000000000033")), Name: "Duncan Shepherd", Email: "duncan-shepherd@car.ds"}
)

func accounts(r account.Repository) {
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
		utils.Must(r.Create(ctx, v.Email, v.Name, account.WithID(v.ID)))
	}

	fmt.Println("created seed users")
}
