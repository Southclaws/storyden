// Code generated by enumerator. DO NOT EDIT.

package subscription

import (
	"database/sql/driver"
	"fmt"
)

type Channel struct {
	v channelEnum
}

var (
	ChannelThread = Channel{channelThread}
	ChannelNone   = Channel{channelNone}
)

func (r Channel) Format(f fmt.State, verb rune) {
	switch verb {
	case 's':
		fmt.Fprint(f, r.v)
	case 'q':
		fmt.Fprintf(f, "%q", r.String())
	default:
		fmt.Fprint(f, r.v)
	}
}
func (r Channel) String() string {
	return string(r.v)
}
func (r Channel) MarshalText() ([]byte, error) {
	return []byte(r.v), nil
}
func (r *Channel) UnmarshalText(__iNpUt__ []byte) error {
	s, err := NewChannel(string(__iNpUt__))
	if err != nil {
		return err
	}
	*r = s
	return nil
}
func (r Channel) Value() (driver.Value, error) {
	return r.v, nil
}
func (r *Channel) Scan(__iNpUt__ any) error {
	s, err := NewChannel(fmt.Sprint(__iNpUt__))
	if err != nil {
		return err
	}
	*r = s
	return nil
}
func NewChannel(__iNpUt__ string) (Channel, error) {
	switch __iNpUt__ {
	case string(channelThread):
		return ChannelThread, nil
	case string(channelNone):
		return ChannelNone, nil
	default:
		return Channel{}, fmt.Errorf("invalid value for type 'Channel': '%s'", __iNpUt__)
	}
}
