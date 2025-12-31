package action_dispatcher

//go:generate go run github.com/Southclaws/enumerator

type contentTypeEnum string

const (
	contentTypeThreads    contentTypeEnum = "threads"
	contentTypeReplies    contentTypeEnum = "replies"
	contentTypeReacts     contentTypeEnum = "reacts"
	contentTypeLikes      contentTypeEnum = "likes"
	contentTypeNodes      contentTypeEnum = "nodes"
	contentTypeCollections contentTypeEnum = "collections"
	contentTypeProfileBio contentTypeEnum = "profile_bio"
)
