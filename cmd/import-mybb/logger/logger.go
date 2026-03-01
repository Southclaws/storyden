package logger

import (
	"fmt"
	"log"
)

// Phase logs a phase header
func Phase(number int, name string) {
	msg := fmt.Sprintf("Phase %d: %s", number, name)
	log.Println(HeaderStyle.Render(msg))
}

// Success logs a success message
func Success(msg string) {
	log.Println(SuccessStyle.Render("✓ " + msg))
}

// Error logs an error message
func Error(msg string) {
	log.Println(ErrorStyle.Render("✗ " + msg))
}

// Account logs an account import with MyBB username → Storyden handle
func Account(uid int, mybbUsername, storydenHandle string, isAdmin bool) {
	adminBadge := ""
	if isAdmin {
		adminBadge = " " + SuccessStyle.Render("ADMIN")
	}

	msg := fmt.Sprintf("%s %s %s %s %s%s",
		AccountLabel.Render("Account"),
		FieldKey.Render(fmt.Sprintf("uid:%d", uid)),
		FieldValue.Render(mybbUsername),
		Arrow,
		FieldValue.Render("@"+storydenHandle),
		adminBadge,
	)
	log.Println(msg)
}

// Category logs a category import
func Category(fid int, name, slug string) {
	msg := fmt.Sprintf("%s %s %s %s %s",
		CategoryLabel.Render("Category"),
		FieldKey.Render(fmt.Sprintf("fid:%d", fid)),
		FieldValue.Render(name),
		Arrow,
		Dim.Render(slug),
	)
	log.Println(msg)
}

// Thread logs a thread post import
func Thread(tid int, subject, slug string) {
	msg := fmt.Sprintf("%s %s %s %s %s",
		PostLabel.Render("Thread"),
		FieldKey.Render(fmt.Sprintf("tid:%d", tid)),
		FieldValue.Render(truncate(subject, 50)),
		Arrow,
		Dim.Render(truncate(slug, 40)),
	)
	log.Println(msg)
}

// Reply logs a reply post import
func Reply(pid int, subject string) {
	msg := fmt.Sprintf("%s %s %s",
		PostLabel.Render("Reply"),
		FieldKey.Render(fmt.Sprintf("pid:%d", pid)),
		FieldValue.Render(truncate(subject, 60)),
	)
	log.Println(msg)
}

// Role logs a role import
func Role(gid int, title string, permissions []string) {
	permStr := ""
	if len(permissions) > 0 {
		permStr = " " + Dim.Render(fmt.Sprintf("(%d perms)", len(permissions)))
	}

	msg := fmt.Sprintf("%s %s %s%s",
		RoleLabel.Render("Role"),
		FieldKey.Render(fmt.Sprintf("gid:%d", gid)),
		FieldValue.Render(title),
		permStr,
	)
	log.Println(msg)
}

// Tag logs a tag import
func Tag(pid int, prefix string) {
	msg := fmt.Sprintf("%s %s %s",
		TagLabel.Render("Tag"),
		FieldKey.Render(fmt.Sprintf("pid:%d", pid)),
		FieldValue.Render(prefix),
	)
	log.Println(msg)
}

// Info logs a general info message
func Info(msg string) {
	log.Println(Dim.Render(msg))
}

// Skip logs a skipped resource
func Skip(resourceType string, reason string) {
	msg := fmt.Sprintf("Skipping %s: %s", resourceType, reason)
	log.Println(Dim.Render("⊘ " + msg))
}

// truncate truncates a string to maxLen with ellipsis
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
