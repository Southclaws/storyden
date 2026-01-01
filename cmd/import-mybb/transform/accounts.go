package transform

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/cmd/import-mybb/loader"
	"github.com/Southclaws/storyden/cmd/import-mybb/logger"
	"github.com/Southclaws/storyden/cmd/import-mybb/writer"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/otp"
	"github.com/alexedwards/argon2id"
	"github.com/rs/xid"
)

func ImportAccounts(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	if len(data.Users) == 0 {
		log.Println("No users to import")
		return nil
	}

	bannedSet := make(map[int]loader.MyBBBanned)
	for _, b := range data.Banned {
		bannedSet[b.UID] = b
	}

	accountBuilders := make([]*ent.AccountCreate, 0, len(data.Users))
	authBuilders := make([]*ent.AuthenticationCreate, 0, len(data.Users))
	emailBuilders := make([]*ent.EmailCreate, 0, len(data.Users))
	roleBuilders := make([]*ent.AccountRolesCreate, 0)

	handleMap := make(map[string]int)
	emailSet := make(map[string]bool)

	for _, user := range data.Users {
		accountID := xid.New()
		w.AccountIDMap[user.UID] = accountID

		bio := buildBio(user, data)
		createdAt := time.Unix(user.RegDate, 0)

		handle := mark.Slugify(user.Username)

		// If slugify returns empty (username was all special chars/emojis), use a fallback
		if handle == "" {
			handle = fmt.Sprintf("user-%d", user.UID)
			logger.Info(fmt.Sprintf("Generated fallback handle for user '%s' (uid:%d) → @%s", user.Username, user.UID, handle))
		}

		// Handle duplicates
		if count, exists := handleMap[handle]; exists {
			handleMap[handle] = count + 1
			handle = fmt.Sprintf("%s-%d", handle, count+1)
		} else {
			handleMap[handle] = 0
		}

		accountBuilder := w.Client().Account.Create().
			SetID(accountID).
			SetHandle(handle).
			SetName(user.Username).
			SetCreatedAt(createdAt)

		if bio != "" {
			accountBuilder.SetBio(bio)
		}

		// Check if user has admin privileges (cancp permission)
		isAdmin := checkUserIsAdmin(user, data.UserGroups)
		accountBuilder.SetAdmin(isAdmin)

		if banned, isBanned := bannedSet[user.UID]; isBanned {
			deletedAt := time.Unix(banned.DateBan, 0)
			accountBuilder.SetDeletedAt(deletedAt)
		}

		accountBuilders = append(accountBuilders, accountBuilder)

		// Log account import with sexy formatting
		// logger.Account(user.UID, user.Username, handle, isAdmin)

		randomPassword, err := generateRandomPassword()
		if err != nil {
			return fmt.Errorf("generate random password for user %d: %w", user.UID, err)
		}

		passwordHash, err := argon2id.CreateHash(randomPassword, argon2id.DefaultParams)
		if err != nil {
			return fmt.Errorf("hash random password for user %d: %w", user.UID, err)
		}

		authID := xid.New()
		authBuilder := w.Client().Authentication.Create().
			SetID(authID).
			SetService("password").
			SetTokenType("password_hash").
			SetIdentifier(accountID.String()).
			SetToken(passwordHash).
			SetAccountAuthentication(accountID).
			SetCreatedAt(createdAt)

		authBuilders = append(authBuilders, authBuilder)

		if user.Email != "" && !emailSet[user.Email] {
			emailSet[user.Email] = true

			verificationCode, err := otp.Generate()
			if err != nil {
				return fmt.Errorf("generate verification code for user %d: %w", user.UID, err)
			}

			emailBuilder := w.Client().Email.Create().
				SetID(xid.New()).
				SetAccountID(accountID).
				SetEmailAddress(user.Email).
				SetVerificationCode(verificationCode).
				SetVerified(true).
				SetCreatedAt(createdAt)

			emailBuilders = append(emailBuilders, emailBuilder)
		}

		roleIDs := parseUserGroups(user, w)
		displayGroupID := w.RoleIDMap[user.DisplayGroup]
		badgeAssigned := false

		for _, roleID := range roleIDs {
			badge := false
			if !badgeAssigned && roleID == displayGroupID {
				badge = true
				badgeAssigned = true
			}

			roleBuilder := w.Client().AccountRoles.Create().
				SetID(xid.New()).
				SetAccountID(accountID).
				SetRoleID(roleID)

			if badge {
				roleBuilder.SetBadge(true)
			}

			roleBuilders = append(roleBuilders, roleBuilder)
		}
	}

	accounts, err := w.CreateAccounts(ctx, accountBuilders)
	if err != nil {
		return fmt.Errorf("create accounts: %w", err)
	}

	_, err = w.CreateAuthentications(ctx, authBuilders)
	if err != nil {
		return fmt.Errorf("create authentications: %w", err)
	}

	_, err = w.CreateEmails(ctx, emailBuilders)
	if err != nil {
		return fmt.Errorf("create emails: %w", err)
	}

	_, err = w.CreateAccountRoles(ctx, roleBuilders)
	if err != nil {
		return fmt.Errorf("create account roles: %w", err)
	}

	log.Printf("Imported %d accounts with %d authentications, %d emails, and %d role assignments",
		len(accounts), len(authBuilders), len(emailBuilders), len(roleBuilders))
	return nil
}

func buildBio(user loader.MyBBUser, data *loader.MyBBData) string {
	var parts []string

	if user.UserTitle != "" {
		parts = append(parts, user.UserTitle)
	}

	if user.Signature != "" {
		parts = append(parts, user.Signature)
	}

	if userFields, ok := data.UserFields[user.UID]; ok && len(userFields.Fields) > 0 {
		for name, value := range userFields.Fields {
			if value != "" {
				parts = append(parts, fmt.Sprintf("**%s**: %s", name, value))
			}
		}
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, "\n\n")
}

func parseUserGroups(user loader.MyBBUser, w *writer.Writer) []xid.ID {
	roleIDSet := make(map[xid.ID]bool)
	roleIDs := make([]xid.ID, 0)

	if roleID, ok := w.RoleIDMap[user.UserGroup]; ok {
		if !roleIDSet[roleID] {
			roleIDSet[roleID] = true
			roleIDs = append(roleIDs, roleID)
		}
	}

	if user.AdditionalGroups != "" {
		groups := strings.Split(user.AdditionalGroups, ",")
		for _, gidStr := range groups {
			gid, err := strconv.Atoi(strings.TrimSpace(gidStr))
			if err != nil {
				continue
			}
			if roleID, ok := w.RoleIDMap[gid]; ok {
				if !roleIDSet[roleID] {
					roleIDSet[roleID] = true
					roleIDs = append(roleIDs, roleID)
				}
			}
		}
	}

	return roleIDs
}

func generateRandomPassword() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func checkUserIsAdmin(user loader.MyBBUser, userGroups []loader.MyBBUserGroup) bool {
	// Check primary usergroup
	for _, group := range userGroups {
		if group.GID == user.UserGroup && group.CanCP == 1 {
			return true
		}
	}

	// Check additional groups
	if user.AdditionalGroups != "" {
		groups := strings.Split(user.AdditionalGroups, ",")
		for _, gidStr := range groups {
			gid, err := strconv.Atoi(strings.TrimSpace(gidStr))
			if err != nil {
				continue
			}
			for _, group := range userGroups {
				if group.GID == gid && group.CanCP == 1 {
					return true
				}
			}
		}
	}

	return false
}
