package loader

import (
	"context"
	"database/sql"
	"fmt"
)

func LoadAll(ctx context.Context, db *sql.DB) (*MyBBData, error) {
	data := &MyBBData{
		UserFields: make(map[int]MyBBUserField),
		Settings:   make(map[string]string),
	}

	var err error

	data.Settings, err = loadSettings(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load settings: %w", err)
	}

	data.UserGroups, err = loadUserGroups(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load usergroups: %w", err)
	}

	data.ProfileFields, err = loadProfileFields(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load profile fields: %w", err)
	}

	data.UserTitles, err = loadUserTitles(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load user titles: %w", err)
	}

	data.Users, err = loadUsers(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load users: %w", err)
	}

	data.UserFields, err = loadUserFields(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load userfields: %w", err)
	}

	data.Banned, err = loadBanned(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load banned: %w", err)
	}

	data.Forums, err = loadForums(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load forums: %w", err)
	}

	data.ThreadPrefixes, err = loadThreadPrefixes(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load thread prefixes: %w", err)
	}

	data.Threads, err = loadThreads(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load threads: %w", err)
	}

	data.Posts, err = loadPosts(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load posts: %w", err)
	}

	data.Reputation, err = loadReputation(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load reputation: %w", err)
	}

	data.ThreadRatings, err = loadThreadRatings(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load thread ratings: %w", err)
	}

	data.ThreadsRead, err = loadThreadsRead(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load threads read: %w", err)
	}

	data.ReportedContent, err = loadReportedContent(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load reported content: %w", err)
	}

	data.Attachments, err = loadAttachments(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("load attachments: %w", err)
	}

	return data, nil
}

func loadUsers(ctx context.Context, db *sql.DB) ([]MyBBUser, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT uid, username, email, usergroup, additionalgroups, displaygroup, usertitle,
		       regdate, lastactive, signature, avatar, website, birthday, reputation, regip, lastip, timezone
		FROM mybb_users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []MyBBUser
	for rows.Next() {
		var u MyBBUser
		if err := rows.Scan(&u.UID, &u.Username, &u.Email, &u.UserGroup, &u.AdditionalGroups,
			&u.DisplayGroup, &u.UserTitle, &u.RegDate, &u.LastActive, &u.Signature,
			&u.Avatar, &u.Website, &u.Birthday, &u.Reputation, &u.RegIP, &u.LastIP, &u.Timezone); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func loadUserFields(ctx context.Context, db *sql.DB) (map[int]MyBBUserField, error) {
	rows, err := db.QueryContext(ctx, `SELECT ufid FROM mybb_userfields`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fields := make(map[int]MyBBUserField)
	for rows.Next() {
		var ufid int
		if err := rows.Scan(&ufid); err != nil {
			return nil, err
		}
		fields[ufid] = MyBBUserField{UFID: ufid, Fields: make(map[string]string)}
	}
	return fields, rows.Err()
}

func loadUserGroups(ctx context.Context, db *sql.DB) ([]MyBBUserGroup, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT gid, title, description, type, canviewthreads, canviewprofiles,
		       canpostthreads, canpostreplys, canratethreads, caneditposts,
		       candeleteposts, candeletethreads, cancp, issupermod,
		       canuploadavatars, canmanageannounce, canmanagemodqueue, canbanusers
		FROM mybb_usergroups
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []MyBBUserGroup
	for rows.Next() {
		var g MyBBUserGroup
		if err := rows.Scan(&g.GID, &g.Title, &g.Description, &g.Type,
			&g.CanViewThreads, &g.CanViewProfiles, &g.CanPostThreads, &g.CanPostReplys,
			&g.CanRateThreads, &g.CanEditPosts, &g.CanDeletePosts, &g.CanDeleteThreads,
			&g.CanCP, &g.IsSuperMod, &g.CanUploadAvatars, &g.CanManageAnnounce,
			&g.CanManageModQueue, &g.CanBanUsers); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func loadForums(ctx context.Context, db *sql.DB) ([]MyBBForum, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT fid, name, description, pid, parentlist, disporder, active, type
		FROM mybb_forums
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forums []MyBBForum
	for rows.Next() {
		var f MyBBForum
		if err := rows.Scan(&f.FID, &f.Name, &f.Description, &f.PID, &f.ParentList, &f.DispOrder, &f.Active, &f.Type); err != nil {
			return nil, err
		}
		forums = append(forums, f)
	}
	return forums, rows.Err()
}

func loadThreads(ctx context.Context, db *sql.DB) ([]MyBBThread, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT tid, fid, subject, prefix, uid, username, dateline, firstpost, lastpost,
		       views, replies, sticky, visible, deletetime
		FROM mybb_threads
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []MyBBThread
	for rows.Next() {
		var t MyBBThread
		if err := rows.Scan(&t.TID, &t.FID, &t.Subject, &t.Prefix, &t.UID, &t.Username,
			&t.DateLine, &t.FirstPost, &t.LastPost, &t.Views, &t.Replies, &t.Sticky, &t.Visible, &t.DeleteTime); err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}
	return threads, rows.Err()
}

func loadPosts(ctx context.Context, db *sql.DB) ([]MyBBPost, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT pid, tid, replyto, fid, subject, uid, username, dateline, message, visible
		FROM mybb_posts
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []MyBBPost
	for rows.Next() {
		var p MyBBPost
		if err := rows.Scan(&p.PID, &p.TID, &p.ReplyTo, &p.FID, &p.Subject, &p.UID,
			&p.Username, &p.DateLine, &p.Message, &p.Visible); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func loadThreadPrefixes(ctx context.Context, db *sql.DB) ([]MyBBThreadPrefix, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT pid, prefix
		FROM mybb_threadprefixes
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prefixes []MyBBThreadPrefix
	for rows.Next() {
		var p MyBBThreadPrefix
		if err := rows.Scan(&p.PID, &p.Prefix); err != nil {
			return nil, err
		}
		prefixes = append(prefixes, p)
	}
	return prefixes, rows.Err()
}

func loadReputation(ctx context.Context, db *sql.DB) ([]MyBBReputation, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT rid, uid, adduid, pid, dateline, comments
		FROM mybb_reputation
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reps []MyBBReputation
	for rows.Next() {
		var r MyBBReputation
		if err := rows.Scan(&r.RID, &r.UID, &r.AddUID, &r.PID, &r.DateAdded, &r.Comments); err != nil {
			return nil, err
		}
		reps = append(reps, r)
	}
	return reps, rows.Err()
}

func loadThreadRatings(ctx context.Context, db *sql.DB) ([]MyBBThreadRating, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT rid, tid, uid, rating
		FROM mybb_threadratings
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []MyBBThreadRating
	for rows.Next() {
		var r MyBBThreadRating
		if err := rows.Scan(&r.RID, &r.TID, &r.UID, &r.Rating); err != nil {
			return nil, err
		}
		ratings = append(ratings, r)
	}
	return ratings, rows.Err()
}

func loadThreadsRead(ctx context.Context, db *sql.DB) ([]MyBBThreadRead, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT tid, uid, dateline
		FROM mybb_threadsread
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reads []MyBBThreadRead
	for rows.Next() {
		var r MyBBThreadRead
		if err := rows.Scan(&r.TID, &r.UID, &r.DateLine); err != nil {
			return nil, err
		}
		reads = append(reads, r)
	}
	return reads, rows.Err()
}

func loadReportedContent(ctx context.Context, db *sql.DB) ([]MyBBReportedContent, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT rid, id, type, uid, dateline, reason
		FROM mybb_reportedcontent
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []MyBBReportedContent
	for rows.Next() {
		var r MyBBReportedContent
		if err := rows.Scan(&r.RID, &r.ID, &r.Type, &r.UID, &r.DateLine, &r.Reason); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, rows.Err()
}

func loadBanned(ctx context.Context, db *sql.DB) ([]MyBBBanned, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT uid, gid, dateline, reason, lifted
		FROM mybb_banned
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banned []MyBBBanned
	for rows.Next() {
		var b MyBBBanned
		var lifted int
		if err := rows.Scan(&b.UID, &b.GID, &b.DateBan, &b.Reason, &lifted); err != nil {
			return nil, err
		}
		banned = append(banned, b)
	}
	return banned, rows.Err()
}

func loadAttachments(ctx context.Context, db *sql.DB) ([]MyBBAttachment, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT aid, pid, filename, filetype, filesize
		FROM mybb_attachments
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []MyBBAttachment
	for rows.Next() {
		var a MyBBAttachment
		if err := rows.Scan(&a.AID, &a.PID, &a.FileName, &a.FileType, &a.FileSize); err != nil {
			return nil, err
		}
		attachments = append(attachments, a)
	}
	return attachments, rows.Err()
}

func loadProfileFields(ctx context.Context, db *sql.DB) ([]MyBBProfileField, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT fid, name, description, type
		FROM mybb_profilefields
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []MyBBProfileField
	for rows.Next() {
		var f MyBBProfileField
		if err := rows.Scan(&f.FID, &f.Name, &f.Description, &f.Type); err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, rows.Err()
}

func loadUserTitles(ctx context.Context, db *sql.DB) ([]MyBBUserTitle, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT utid, posts, title
		FROM mybb_usertitles
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var titles []MyBBUserTitle
	for rows.Next() {
		var t MyBBUserTitle
		if err := rows.Scan(&t.UTID, &t.Posts, &t.Title); err != nil {
			return nil, err
		}
		titles = append(titles, t)
	}
	return titles, rows.Err()
}

func loadSettings(ctx context.Context, db *sql.DB) (map[string]string, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT name, value
		FROM mybb_settings
		WHERE name IN ('bbname', 'bburl', 'tagline')
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		settings[name] = value
	}
	return settings, rows.Err()
}
