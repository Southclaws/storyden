// Code generated by enumerator. DO NOT EDIT.

package collection

import (
	"database/sql/driver"
	"fmt"
)

type MembershipType struct {
	v membershipTypeEnum
}

var (
	MembershipTypeNormal             = MembershipType{membershipTypeNormal}
	MembershipTypeSubmissionReview   = MembershipType{membershipTypeSubmissionReview}
	MembershipTypeSubmissionAccepted = MembershipType{membershipTypeSubmissionAccepted}
)

func (r MembershipType) Format(f fmt.State, verb rune) {
	switch verb {
	case 's':
		fmt.Fprint(f, r.v)
	case 'q':
		fmt.Fprintf(f, "%q", r.String())
	default:
		fmt.Fprint(f, r.v)
	}
}
func (r MembershipType) String() string {
	return string(r.v)
}
func (r MembershipType) MarshalText() ([]byte, error) {
	return []byte(r.v), nil
}
func (r *MembershipType) UnmarshalText(__iNpUt__ []byte) error {
	s, err := NewMembershipType(string(__iNpUt__))
	if err != nil {
		return err
	}
	*r = s
	return nil
}
func (r MembershipType) Value() (driver.Value, error) {
	return r.v, nil
}
func (r *MembershipType) Scan(__iNpUt__ any) error {
	s, err := NewMembershipType(fmt.Sprint(__iNpUt__))
	if err != nil {
		return err
	}
	*r = s
	return nil
}
func NewMembershipType(__iNpUt__ string) (MembershipType, error) {
	switch __iNpUt__ {
	case string(membershipTypeNormal):
		return MembershipTypeNormal, nil
	case string(membershipTypeSubmissionReview):
		return MembershipTypeSubmissionReview, nil
	case string(membershipTypeSubmissionAccepted):
		return MembershipTypeSubmissionAccepted, nil
	default:
		return MembershipType{}, fmt.Errorf("invalid value for type 'MembershipType': '%s'", __iNpUt__)
	}
}
