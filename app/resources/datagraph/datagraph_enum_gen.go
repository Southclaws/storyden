// Code generated by enumerator. DO NOT EDIT.

package datagraph

import (
	"database/sql/driver"
	"fmt"
)

type Kind struct {
	v kindEnum
}

var (
	KindThread  = Kind{kindThread}
	KindReply   = Kind{kindReply}
	KindCluster = Kind{kindCluster}
	KindItem    = Kind{kindItem}
	KindLink    = Kind{kindLink}
)

func (r Kind) Format(f fmt.State, verb rune) {
	switch verb {
	case 's':
		fmt.Fprint(f, r.v)
	case 'q':
		fmt.Fprintf(f, "%q", r.String())
	default:
		fmt.Fprint(f, r.v)
	}
}
func (r Kind) String() string {
	return string(r.v)
}
func (r Kind) MarshalText() ([]byte, error) {
	return []byte(r.v), nil
}
func (r *Kind) UnmarshalText(in []byte) error {
	s, err := NewKind(string(in))
	if err != nil {
		return err
	}
	*r = s
	return nil
}
func (r Kind) Value() (driver.Value, error) {
	return r.v, nil
}
func (r *Kind) Scan(in any) error {
	s, err := NewKind(fmt.Sprint(in))
	if err != nil {
		return err
	}
	*r = s
	return nil
}
func NewKind(in string) (Kind, error) {
	switch in {
	case string(kindThread):
		return KindThread, nil
	case string(kindReply):
		return KindReply, nil
	case string(kindCluster):
		return KindCluster, nil
	case string(kindItem):
		return KindItem, nil
	case string(kindLink):
		return KindLink, nil
	default:
		return Kind{}, fmt.Errorf("invalid value for type 'Kind': '%s'", in)
	}
}
