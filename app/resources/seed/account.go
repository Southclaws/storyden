package seed

import (
	"context"
	"fmt"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
)

const SeedPassword = `$argon2id$v=19$m=65536,t=1,p=2$MAwllQoeGcxCPOC52OQwZA$jLlzHsmSHmQPbpQ6Y5+877NlacOYeyqEqWoKJJXRcHM`

var (
	Account_001_Odin      = account.Account{ID: account.AccountID(id("00000000000000000010")), Name: "Odin", Handle: "odin", Admin: true}
	Account_002_Frigg     = account.Account{ID: account.AccountID(id("00000000000000000020")), Name: "Frigg", Handle: "frigg", Admin: true}
	Account_003_Baldur    = account.Account{ID: account.AccountID(id("00000000000000000030")), Name: "Baldur", Handle: "baldur"}
	Account_004_Loki      = account.Account{ID: account.AccountID(id("00000000000000000040")), Name: "Loki", Handle: "loki"}
	Account_005_Þórr      = account.Account{ID: account.AccountID(id("00000000000000000050")), Name: "Þórr", Handle: "þórr"}
	Account_006_Freyja    = account.Account{ID: account.AccountID(id("00000000000000000060")), Name: "Freyja", Handle: "freyja"}
	Account_007_Freyr     = account.Account{ID: account.AccountID(id("00000000000000000070")), Name: "Freyr", Handle: "freyr"}
	Account_008_Heimdallr = account.Account{ID: account.AccountID(id("00000000000000000080")), Name: "Heimdallr", Handle: "heimdallr"}
	Account_009_Hel       = account.Account{ID: account.AccountID(id("00000000000000000090")), Name: "Hel", Handle: "hel"}
	Account_010_Víðarr    = account.Account{ID: account.AccountID(id("00000000000000000100")), Name: "Víðarr", Handle: "víðarr"}
	Account_011_Váli      = account.Account{ID: account.AccountID(id("00000000000000000110")), Name: "Váli", Handle: "váli"}
	Account_012_Njörðr    = account.Account{ID: account.AccountID(id("00000000000000000120")), Name: "Njörðr", Handle: "njörðr"}
	Account_013_Annoying  = account.Account{ID: account.AccountID(id("00000000000000000130")), Name: "WWWWWWWWWWWWWWWWWWWWWWWWWWWWWW", Handle: "wwwwwwwwwwwwwwwwwwwwwwwwwwwwww"}

	Accounts = []account.Account{
		Account_001_Odin,
		Account_002_Frigg,
		Account_003_Baldur,
		Account_004_Loki,
		Account_005_Þórr,
		Account_006_Freyja,
		Account_007_Freyr,
		Account_008_Heimdallr,
		Account_009_Hel,
		Account_010_Víðarr,
		Account_011_Váli,
		Account_012_Njörðr,
		Account_013_Annoying,
	}
)

func accounts(r *account_writer.Writer, auth authentication.Repository) {
	ctx := context.Background()

	for _, v := range Accounts {
		SeedAccount(ctx, r, auth, v)
	}

	for i := 0; i < 100; i++ {
		SeedAccount(ctx, r, auth, account.Account{
			ID:     account.AccountID(xid.New()),
			Name:   fmt.Sprintf("Acc%d", i),
			Handle: fmt.Sprintf("acc%d", i),
		})
	}

	fmt.Println("created seed users")
}

func SeedAccount(ctx context.Context, r *account_writer.Writer, auth authentication.Repository, v account.Account) {
	acc, err := r.Create(ctx, v.Handle,
		account_writer.WithID(v.ID),
		account_writer.WithName(v.Name),
		account_writer.WithBio(v.Bio),
		account_writer.WithAdmin(v.Admin),
	)
	if err != nil {
		panic(err)
	}

	// TODO: email+password auth provider.
	// email := acc.Handle + "@storyd.en"

	if _, err = auth.Create(
		ctx,
		acc.ID,
		authentication.ServicePassword,
		authentication.TokenTypePasswordHash,
		acc.ID.String(),
		SeedPassword,
		nil,
	); err != nil {
		panic(err)
	}
}

func SeedAccountUnique(ctx context.Context, r *account_writer.Writer, auth authentication.Repository, opts ...account_writer.Option) {
	id := account.AccountID(xid.New())
	handle := fmt.Sprintf("acc%s", id.String())

	opts = append(opts, account_writer.WithID(id), account_writer.WithName(handle))

	acc, err := r.Create(ctx, handle, opts...)
	if err != nil {
		panic(err)
	}

	if _, err = auth.Create(
		ctx,
		acc.ID,
		authentication.ServicePassword,
		authentication.TokenTypePasswordHash,
		"password",
		SeedPassword,
		nil,
	); err != nil {
		panic(err)
	}
}
