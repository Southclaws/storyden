package library

//go:generate go run github.com/Southclaws/enumerator

type propertyTypeEnum string

const (
	propertyTypeEnumText      propertyTypeEnum = "text"
	propertyTypeEnumNumber    propertyTypeEnum = "number"
	propertyTypeEnumTimestamp propertyTypeEnum = "timestamp"
	propertyTypeEnumBoolean   propertyTypeEnum = "boolean"
)
