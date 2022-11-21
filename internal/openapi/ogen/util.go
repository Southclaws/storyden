package ogen

func NewOptStringPtr(v *string) OptString {
	if v == nil {
		return OptString{}
	}
	return OptString{
		Value: *v,
		Set:   true,
	}
}
