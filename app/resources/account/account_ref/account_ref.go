package account_ref

import "github.com/rs/xid"

type ID xid.ID

func (u ID) String() string { return xid.ID(u).String() }

func (u ID) MarshalJSON() ([]byte, error) {
	return xid.ID(u).MarshalJSON()
}

func (u *ID) UnmarshalJSON(data []byte) error {
	var id xid.ID
	if err := id.UnmarshalJSON(data); err != nil {
		return err
	}
	*u = ID(id)
	return nil
}
