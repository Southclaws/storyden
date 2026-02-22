package loader

import "time"

type MyBBData struct {
	Users           []MyBBUser
	UserFields      map[int]MyBBUserField
	UserGroups      []MyBBUserGroup
	Forums          []MyBBForum
	Threads         []MyBBThread
	Posts           []MyBBPost
	ThreadPrefixes  []MyBBThreadPrefix
	Reputation      []MyBBReputation
	ThreadRatings   []MyBBThreadRating
	ThreadsRead     []MyBBThreadRead
	ReportedContent []MyBBReportedContent
	Banned          []MyBBBanned
	Attachments     []MyBBAttachment
	ProfileFields   []MyBBProfileField
	UserTitles      []MyBBUserTitle
	Settings        map[string]string
}

type MyBBUser struct {
	UID              int
	Username         string
	Email            string
	UserGroup        int
	AdditionalGroups string
	DisplayGroup     int
	UserTitle        string
	RegDate          int64
	LastActive       int64
	Signature        string
	Avatar           string
	Website          string
	Birthday         string
	Reputation       int
	RegIP            []byte
	LastIP           []byte
	Timezone         string
}

type MyBBUserField struct {
	UFID   int
	Fields map[string]string
}

type MyBBUserGroup struct {
	GID                int
	Title              string
	Description        string
	Type               int
	CanViewThreads     int
	CanViewProfiles    int
	CanPostThreads     int
	CanPostReplys      int
	CanRateThreads     int
	CanEditPosts       int
	CanDeletePosts     int
	CanDeleteThreads   int
	CanCP              int
	IsSuperMod         int
	CanUploadAvatars   int
	CanManageAnnounce  int
	CanManageModQueue  int
	CanBanUsers        int
}

type MyBBForum struct {
	FID         int
	Name        string
	Description string
	PID         int
	ParentList  string
	DispOrder   int
	Active      int
	Type        string
}

type MyBBThread struct {
	TID         int
	FID         int
	Subject     string
	Prefix      int
	UID         int
	Username    string
	DateLine    int64
	FirstPost   int
	LastPost    int64
	Views       int
	Replies     int
	Sticky      int
	Visible     int
	DeleteTime  int64
}

type MyBBPost struct {
	PID      int
	TID      int
	ReplyTo  int
	FID      int
	Subject  string
	UID      int
	Username string
	DateLine int64
	Message  string
	Visible  int
}

type MyBBThreadPrefix struct {
	PID    int
	Prefix string
}

type MyBBReputation struct {
	RID       int
	UID       int
	AddUID    int
	PID       int
	DateAdded int64
	Comments  string
}

type MyBBThreadRating struct {
	RID    int
	TID    int
	UID    int
	Rating int
}

type MyBBThreadRead struct {
	TID      int
	UID      int
	DateLine int64
}

type MyBBReportedContent struct {
	RID      int
	ID       int
	Type     string
	UID      int
	DateLine int64
	Reason   string
}

type MyBBBanned struct {
	UID     int
	GID     int
	DateBan int64
	Reason  string
}

type MyBBAttachment struct {
	AID      int
	PID      int
	FileName string
	FileType string
	FileSize int
}

type MyBBProfileField struct {
	FID         int
	Name        string
	Description string
	Type        string
}

type MyBBUserTitle struct {
	UTID  int
	Posts int
	Title string
}

func unixToTime(unix int64) time.Time {
	if unix == 0 {
		return time.Time{}
	}
	return time.Unix(unix, 0)
}
