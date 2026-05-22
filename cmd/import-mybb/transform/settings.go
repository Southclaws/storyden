package transform

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/cmd/import-mybb/loader"
	"github.com/Southclaws/storyden/cmd/import-mybb/writer"
)

const StorydenPrimarySettingsKey = "storyden_system"

func ImportSettings(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	// Build Storyden settings from MyBB settings and force email auth mode.
	// MyBB password hashes are not migrated as-is, so accounts authenticate via email.
	storydenSettings := buildStorydenSettings(data.Settings)
	if len(data.Settings) == 0 {
		log.Println("No MyBB settings found, using defaults and forcing email authentication mode")
	}

	// Serialize to JSON
	settingsJSON, err := json.Marshal(storydenSettings)
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}

	// Check if settings already exist
	_, err = w.Client().Setting.Get(ctx, StorydenPrimarySettingsKey)
	if err == nil {
		// Update existing settings
		_, err = w.Client().Setting.UpdateOneID(StorydenPrimarySettingsKey).
			SetValue(string(settingsJSON)).
			SetUpdatedAt(time.Now()).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("update settings: %w", err)
		}
		log.Printf("Updated Storyden settings with MyBB forum name: %s", data.Settings["bbname"])
	} else {
		// Create new settings
		_, err = w.Client().Setting.Create().
			SetID(StorydenPrimarySettingsKey).
			SetValue(string(settingsJSON)).
			SetUpdatedAt(time.Now()).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("create settings: %w", err)
		}
		log.Printf("Created Storyden settings with MyBB forum name: %s", data.Settings["bbname"])
	}

	return nil
}

func buildStorydenSettings(mybbSettings map[string]string) settings.Settings {
	storydenSettings := settings.DefaultSettings
	storydenSettings.AuthenticationMode = opt.New(authentication.ModeEmail)

	bbname, hasTitle := mybbSettings["bbname"]
	if hasTitle && strings.TrimSpace(bbname) != "" {
		storydenSettings.Title = opt.New(bbname)
	}

	tagline, hasTagline := mybbSettings["tagline"]
	if hasTagline && strings.TrimSpace(tagline) != "" {
		storydenSettings.Description = opt.New(tagline)
	}

	return storydenSettings
}
