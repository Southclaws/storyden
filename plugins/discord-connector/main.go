package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	sdk "github.com/Southclaws/storyden/sdk/go/storyden"
)

const (
	latestCommandName        = "latest"
	latestCommandDescription = "Show the latest published Storyden thread."
	saveCommandName          = "save"
	saveCommandDescription   = "Save the latest URL posted in this channel to Storyden Links."
	searchCommandName        = "search"
	searchCommandDescription = "Search Storyden from Discord."
	searchCommandOptionQuery = "query"

	initialConfigRetryInterval = 2 * time.Second
	apiRequestTimeout          = 10 * time.Second
	configureTimeout           = 15 * time.Second
	searchResultLimit          = 5
)

var urlPattern = regexp.MustCompile(`https?://[^\s<>()]+`)

type pluginConfig struct {
	DiscordToken string
	ChannelID    string
}

type discordCommand struct {
	ID      string
	GuildID string
	Name    string
}

type runtimeState struct {
	config pluginConfig

	discord *discordgo.Session
	api     *openapi.ClientWithResponses
	webBase *url.URL

	applicationID   string
	discordCommands []discordCommand
}

type connector struct {
	plugin *sdk.Plugin
	logger *slog.Logger

	mu    sync.RWMutex
	state *runtimeState
}

//go:generate ./package.nu

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("starting discord connector plugin")
	if err := godotenv.Load(); err == nil {
		logger.Info("loaded environment from .env")
	}

	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stopSignal()
	defer logger.Info("discord connector process exiting")

	pl, err := sdk.New(ctx)
	if err != nil {
		logger.Error("failed to create plugin", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("storyden plugin SDK initialised")

	conn := newConnector(pl, logger)
	defer conn.close()

	pl.OnConfigure(conn.handleConfigure)
	pl.OnThreadPublished(conn.handleThreadPublished)

	go conn.syncInitialConfig(ctx)
	go func() {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.Canceled) {
			logger.Info("received shutdown signal")
		}
	}()

	logger.Info("starting plugin runtime loop")
	if err := pl.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Info("plugin runtime loop stopped due to shutdown signal")
			return
		}
		logger.Error("plugin stopped with error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("plugin runtime loop stopped cleanly")
}

func newConnector(plugin *sdk.Plugin, logger *slog.Logger) *connector {
	return &connector{
		plugin: plugin,
		logger: logger,
	}
}

func (c *connector) handleConfigure(ctx context.Context, config map[string]any) error {
	cfg, err := parseConfig(config)
	if err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, configureTimeout)
	defer cancel()

	return c.reconfigure(timeoutCtx, cfg)
}

func (c *connector) syncInitialConfig(ctx context.Context) {
	c.logger.Info("syncing initial plugin configuration")
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("initial configuration sync cancelled")
			return
		default:
		}

		cfgMap, err := c.getConfig(ctx)
		if err != nil {
			c.logger.Debug("initial get_config failed, retrying", slog.String("error", err.Error()))
			time.Sleep(initialConfigRetryInterval)
			continue
		}

		cfg, err := parseConfig(cfgMap)
		if err != nil {
			if err := c.setUnconfigured(); err != nil {
				c.logger.Warn("failed to clear runtime state", slog.String("error", err.Error()))
			}
			c.logger.Info("plugin started without valid configuration", slog.String("error", err.Error()))
			return
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, configureTimeout)
		err = c.reconfigure(timeoutCtx, cfg)
		cancel()
		if err != nil {
			c.logger.Warn("initial configuration rejected", slog.String("error", err.Error()))
			time.Sleep(initialConfigRetryInterval)
			continue
		}

		c.logger.Info("initial configuration loaded", slog.String("channel_id", cfg.ChannelID))
		return
	}
}

func (c *connector) getConfig(ctx context.Context) (map[string]any, error) {
	resp, err := c.plugin.Send(ctx, rpc.RPCRequestGetConfig{
		Jsonrpc: "2.0",
		Method:  "get_config",
	})
	if err != nil {
		return nil, err
	}

	cfgResp, ok := resp.(*rpc.RPCResponseGetConfig)
	if !ok {
		return nil, fmt.Errorf("unexpected get_config response: %T", resp)
	}

	if cfgResp.Config == nil {
		return map[string]any{}, nil
	}

	return cfgResp.Config, nil
}

func (c *connector) reconfigure(ctx context.Context, cfg pluginConfig) error {
	if cfg.DiscordToken == "" {
		return fmt.Errorf("discord_token is required")
	}
	if cfg.ChannelID == "" {
		return fmt.Errorf("channel_id is required")
	}

	current := c.getState()
	if current != nil && current.config == cfg {
		c.logger.Info("configuration unchanged, skipping reconfigure")
		return nil
	}

	c.logger.Info("applying plugin configuration", slog.String("channel_id", cfg.ChannelID))

	apiClient, webBase, err := c.buildAPIClient(ctx)
	if err != nil {
		return err
	}

	discord, appID, commands, err := c.buildDiscordSession(cfg)
	if err != nil {
		return err
	}

	next := &runtimeState{
		config:          cfg,
		discord:         discord,
		api:             apiClient,
		webBase:         webBase,
		applicationID:   appID,
		discordCommands: commands,
	}

	c.mu.Lock()
	old := c.state
	c.state = next
	c.mu.Unlock()

	c.shutdownState(old)

	c.logger.Info("configuration applied",
		slog.String("channel_id", cfg.ChannelID),
		slog.Int("discord_commands", len(commands)),
	)
	return nil
}

func (c *connector) setUnconfigured() error {
	c.logger.Info("clearing plugin runtime state")
	c.mu.Lock()
	old := c.state
	c.state = nil
	c.mu.Unlock()

	c.shutdownState(old)
	return nil
}

func (c *connector) buildAPIClient(ctx context.Context) (*openapi.ClientWithResponses, *url.URL, error) {
	client, err := c.plugin.BuildAPIClient(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build API client: %w", err)
	}

	rawClient, ok := client.ClientInterface.(*openapi.Client)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected API client implementation: %T", client.ClientInterface)
	}

	apiBaseURL, err := url.Parse(rawClient.Server)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse API base URL: %w", err)
	}

	apiBaseURL.Path = strings.TrimSuffix(apiBaseURL.Path, "/")
	apiBaseURL.Path = strings.TrimSuffix(apiBaseURL.Path, "/api")
	if apiBaseURL.Path == "" {
		apiBaseURL.Path = "/"
	}

	c.logger.Info("constructed Storyden API client", slog.String("api_base_url", apiBaseURL.String()))

	return client, apiBaseURL, nil
}

func (c *connector) buildDiscordSession(cfg pluginConfig) (*discordgo.Session, string, []discordCommand, error) {
	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to create discord session: %w", err)
	}

	dg.Client = &http.Client{Timeout: apiRequestTimeout}
	dg.Identify.Intents = discordgo.IntentsGuilds
	dg.AddHandler(c.handleDiscordInteraction)

	channel, err := dg.Channel(cfg.ChannelID)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to access configured Discord channel: %w", err)
	}
	c.logger.Info("connected to configured Discord channel",
		slog.String("channel_id", channel.ID),
		slog.String("guild_id", channel.GuildID),
	)

	if err := dg.Open(); err != nil {
		return nil, "", nil, fmt.Errorf("failed to open discord gateway: %w", err)
	}
	c.logger.Info("opened Discord gateway session")

	self, err := dg.User("@me")
	if err != nil {
		_ = dg.Close()
		return nil, "", nil, fmt.Errorf("failed to lookup bot user: %w", err)
	}

	registeredCommands := make([]discordCommand, 0, 3)
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        latestCommandName,
			Description: latestCommandDescription,
			Type:        discordgo.ChatApplicationCommand,
		},
		{
			Name:        saveCommandName,
			Description: saveCommandDescription,
			Type:        discordgo.ChatApplicationCommand,
		},
		{
			Name:        searchCommandName,
			Description: searchCommandDescription,
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        searchCommandOptionQuery,
					Description: "What to search for.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
	}

	for _, command := range commands {
		created, createErr := dg.ApplicationCommandCreate(self.ID, channel.GuildID, command)
		if createErr != nil {
			for _, registered := range registeredCommands {
				_ = dg.ApplicationCommandDelete(self.ID, registered.GuildID, registered.ID)
			}
			_ = dg.Close()
			return nil, "", nil, fmt.Errorf("failed to register /%s command: %w", command.Name, createErr)
		}

		registeredCommands = append(registeredCommands, discordCommand{
			ID:      created.ID,
			GuildID: channel.GuildID,
			Name:    command.Name,
		})
		c.logger.Info("registered Discord command",
			slog.String("command", command.Name),
			slog.String("command_id", created.ID),
		)
	}

	c.logger.Info("discord session ready", slog.String("application_id", self.ID))

	return dg, self.ID, registeredCommands, nil
}

func (c *connector) handleThreadPublished(ctx context.Context, event *rpc.EventThreadPublished) error {
	state := c.getState()
	if state == nil {
		c.logger.Warn("thread published event ignored because plugin is not configured")
		return nil
	}

	summary, err := c.fetchThreadSummary(ctx, state.api, event.ID.String())
	if err != nil {
		c.logger.Warn("failed to load thread details", slog.String("error", err.Error()), slog.String("thread_id", event.ID.String()))
		summary = threadSummary{ID: event.ID.String()}
	}

	msg := formatThreadPublishedMessage(*state.webBase, summary)
	c.logger.Info("sending thread published notification to Discord",
		slog.String("thread_id", event.ID.String()),
		slog.String("channel_id", state.config.ChannelID),
	)
	if _, err := state.discord.ChannelMessageSendComplex(state.config.ChannelID, &discordgo.MessageSend{
		Content: msg,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{},
		},
	}); err != nil {
		c.logger.Error("failed to send discord notification", slog.String("error", err.Error()), slog.String("channel_id", state.config.ChannelID))
	}

	return nil
}

func (c *connector) handleDiscordInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	state := c.getState()
	if state == nil || i == nil {
		return
	}

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.ApplicationCommandData().Name
	if !isSupportedCommand(commandName) {
		return
	}
	c.logger.Info("received Discord command",
		slog.String("command", commandName),
		slog.String("channel_id", i.ChannelID),
	)

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			AllowedMentions: &discordgo.MessageAllowedMentions{Parse: []discordgo.AllowedMentionType{}},
		},
	}); err != nil {
		c.logger.Error("failed to defer interaction", slog.String("error", err.Error()), slog.String("command", commandName))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), apiRequestTimeout)
	defer cancel()

	content, err := c.handleCommand(ctx, state, i.ChannelID, i.ApplicationCommandData())
	if err != nil {
		switch commandName {
		case latestCommandName:
			content = "Failed to load latest thread: " + sanitizeDiscordMessage(err.Error())
		case saveCommandName:
			content = "Failed to save link: " + sanitizeDiscordMessage(err.Error())
		case searchCommandName:
			content = "Failed to search Storyden: " + sanitizeDiscordMessage(err.Error())
		default:
			content = "Command failed: " + sanitizeDiscordMessage(err.Error())
		}
	}

	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:         &content,
		AllowedMentions: &discordgo.MessageAllowedMentions{Parse: []discordgo.AllowedMentionType{}},
	}); err != nil {
		c.logger.Error("failed to edit interaction response", slog.String("error", err.Error()), slog.String("command", commandName))
		return
	}

	c.logger.Info("handled Discord command", slog.String("command", commandName))
}

func isSupportedCommand(name string) bool {
	return slices.Contains([]string{latestCommandName, saveCommandName, searchCommandName}, name)
}

func (c *connector) handleCommand(ctx context.Context, state *runtimeState, channelID string, data discordgo.ApplicationCommandInteractionData) (string, error) {
	switch data.Name {
	case latestCommandName:
		latest, err := c.fetchLatestThread(ctx, state.api)
		if err != nil {
			return "", err
		}
		return formatLatestThreadMessage(*state.webBase, latest), nil
	case saveCommandName:
		return c.handleSave(ctx, state, channelID)
	case searchCommandName:
		query := strings.TrimSpace(getCommandOptionString(data, searchCommandOptionQuery))
		if query == "" {
			return "Please provide a search query.", nil
		}
		return c.handleSearch(ctx, state, query)
	default:
		return "", fmt.Errorf("unsupported command: %s", data.Name)
	}
}

func getCommandOptionString(data discordgo.ApplicationCommandInteractionData, name string) string {
	for _, option := range data.Options {
		if option.Name == name {
			return option.StringValue()
		}
	}
	return ""
}

func (c *connector) handleSave(ctx context.Context, state *runtimeState, channelID string) (string, error) {
	if strings.TrimSpace(channelID) == "" {
		channelID = state.config.ChannelID
	}

	savedURL, err := c.findLatestChannelURL(state.discord, channelID)
	if err != nil {
		return "", err
	}

	resp, err := state.api.LinkCreateWithResponse(ctx, openapi.LinkCreateJSONRequestBody{
		Url: openapi.URL(savedURL),
	})
	if err != nil {
		return "", err
	}
	if resp.JSON200 == nil {
		return "", fmt.Errorf("unexpected link create response status: %d", resp.StatusCode())
	}

	ref := resp.JSON200
	var b strings.Builder
	b.WriteString("Saved latest link to Storyden.")
	b.WriteString("\n")
	if ref.Title != nil && strings.TrimSpace(string(*ref.Title)) != "" {
		b.WriteString("**")
		b.WriteString(sanitizeDiscordMessage(string(*ref.Title)))
		b.WriteString("**")
		b.WriteString("\n")
	}
	b.WriteString(string(ref.Url))
	b.WriteString("\n")
	b.WriteString("Link ID: `")
	b.WriteString(string(ref.Id))
	b.WriteString("`")

	return truncateDiscordMessage(b.String()), nil
}

func (c *connector) findLatestChannelURL(dg *discordgo.Session, channelID string) (string, error) {
	messages, err := dg.ChannelMessages(channelID, 50, "", "", "")
	if err != nil {
		return "", fmt.Errorf("failed to read channel messages: %w", err)
	}

	for _, message := range messages {
		if message == nil || message.Author == nil || message.Author.Bot {
			continue
		}

		for _, candidate := range urlPattern.FindAllString(message.Content, -1) {
			sanitized := sanitizeURLCandidate(candidate)
			if sanitized == "" {
				continue
			}
			return sanitized, nil
		}
	}

	return "", fmt.Errorf("no URL found in recent channel messages")
}

func sanitizeURLCandidate(candidate string) string {
	trimmed := strings.TrimSpace(candidate)
	trimmed = strings.Trim(trimmed, "<>[]{}()\"'`")
	trimmed = strings.TrimRight(trimmed, ".,;:!?")
	if trimmed == "" {
		return ""
	}

	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return ""
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ""
	}

	return parsed.String()
}

func (c *connector) handleSearch(ctx context.Context, state *runtimeState, query string) (string, error) {
	page := openapi.PaginationQuery("1")
	searchQuery := openapi.SearchQuery(query)

	resp, err := state.api.LinkListWithResponse(ctx, &openapi.LinkListParams{
		Q:    &searchQuery,
		Page: &page,
	})
	if err != nil {
		return "", err
	}
	if resp.JSON200 == nil {
		return "", fmt.Errorf("unexpected link search response status: %d", resp.StatusCode())
	}

	links := resp.JSON200.Links
	if len(links) == 0 {
		return fmt.Sprintf("No results for `%s`.", sanitizeDiscordMessage(query)), nil
	}

	var b strings.Builder
	b.WriteString("Link results for `")
	b.WriteString(sanitizeDiscordMessage(query))
	b.WriteString("`")

	count := min(searchResultLimit, len(links))
	for i := range count {
		line := formatLinkSearchResultLine(links[i])
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("%d. %s", i+1, line))
	}

	return truncateDiscordMessage(b.String()), nil
}

func formatLinkSearchResultLine(link openapi.LinkReference) string {
	title := ""
	if link.Title != nil {
		title = strings.TrimSpace(string(*link.Title))
	}

	label := strings.TrimSpace(title)
	if label == "" {
		label = string(link.Domain)
	}

	var b strings.Builder
	b.WriteString(sanitizeDiscordMessage(label))
	b.WriteString(" - ")
	b.WriteString(string(link.Url))

	return b.String()
}

func respondImmediate(s *discordgo.Session, i *discordgo.Interaction, message string, ephemeral bool) error {
	flags := discordgo.MessageFlags(0)
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}

	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sanitizeDiscordMessage(message),
			Flags:   flags,
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		},
	})
}

func (c *connector) close() {
	c.logger.Info("shutting down connector")
	c.mu.Lock()
	state := c.state
	c.state = nil
	c.mu.Unlock()

	c.shutdownState(state)
}

func (c *connector) shutdownState(state *runtimeState) {
	if state == nil || state.discord == nil {
		return
	}

	c.logger.Info("shutting down runtime state",
		slog.String("application_id", state.applicationID),
		slog.Int("discord_commands", len(state.discordCommands)),
	)

	for _, command := range state.discordCommands {
		if command.ID == "" || state.applicationID == "" {
			continue
		}

		if err := state.discord.ApplicationCommandDelete(state.applicationID, command.GuildID, command.ID); err != nil {
			c.logger.Warn("failed to delete Discord command",
				slog.String("error", err.Error()),
				slog.String("command_id", command.ID),
				slog.String("command_name", command.Name),
			)
		}
	}

	if err := state.discord.Close(); err != nil {
		c.logger.Warn("failed to close Discord session", slog.String("error", err.Error()))
		return
	}

	c.logger.Info("closed Discord session")
}

func (c *connector) getState() *runtimeState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

type threadSummary struct {
	ID     string
	Title  string
	Slug   string
	Author string
}

func (c *connector) fetchThreadSummary(ctx context.Context, api *openapi.ClientWithResponses, threadID string) (threadSummary, error) {
	resp, err := api.ThreadGetWithResponse(ctx, openapi.ThreadMarkParam(threadID), nil)
	if err != nil {
		return threadSummary{}, err
	}
	if resp.JSON200 == nil {
		return threadSummary{}, fmt.Errorf("unexpected thread response status: %d", resp.StatusCode())
	}

	thread := resp.JSON200
	return threadSummary{
		ID:     string(thread.Id),
		Title:  string(thread.Title),
		Slug:   string(thread.Slug),
		Author: string(thread.Author.Handle),
	}, nil
}

func (c *connector) fetchLatestThread(ctx context.Context, api *openapi.ClientWithResponses) (threadSummary, error) {
	page := openapi.PaginationQuery("1")
	ignorePinned := openapi.ThreadsIgnorePinnedQuery(true)

	resp, err := api.ThreadListWithResponse(ctx, &openapi.ThreadListParams{
		Page:         &page,
		IgnorePinned: &ignorePinned,
	})
	if err != nil {
		return threadSummary{}, err
	}
	if resp.JSON200 == nil {
		return threadSummary{}, fmt.Errorf("unexpected thread list response status: %d", resp.StatusCode())
	}
	if len(resp.JSON200.Threads) == 0 {
		return threadSummary{}, fmt.Errorf("no threads found")
	}

	thread := resp.JSON200.Threads[0]
	return threadSummary{
		ID:     string(thread.Id),
		Title:  string(thread.Title),
		Slug:   string(thread.Slug),
		Author: string(thread.Author.Handle),
	}, nil
}

func parseConfig(in map[string]any) (pluginConfig, error) {
	token, ok := extractString(in, "discord_token")
	if !ok {
		return pluginConfig{}, fmt.Errorf("discord_token must be a non-empty string")
	}

	channelID, ok := extractString(in, "channel_id")
	if !ok {
		return pluginConfig{}, fmt.Errorf("channel_id must be a non-empty string")
	}

	return pluginConfig{
		DiscordToken: token,
		ChannelID:    channelID,
	}, nil
}

func extractString(in map[string]any, key string) (string, bool) {
	if in == nil {
		return "", false
	}

	raw, ok := in[key]
	if !ok {
		return "", false
	}

	value, ok := raw.(string)
	if !ok {
		return "", false
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}

	return value, true
}

func formatThreadPublishedMessage(baseURL url.URL, summary threadSummary) string {
	if summary.Title == "" {
		summary.Title = "Untitled thread"
	}

	var b strings.Builder
	b.WriteString("New thread published")
	b.WriteString("\n")
	b.WriteString("**")
	b.WriteString(sanitizeDiscordMessage(summary.Title))
	b.WriteString("**")

	if summary.Author != "" {
		b.WriteString(" by @")
		b.WriteString(sanitizeDiscordMessage(summary.Author))
	}

	if summary.Slug != "" {
		b.WriteString("\n")
		b.WriteString(strings.TrimRight(baseURL.String(), "/"))
		b.WriteString("/t/")
		b.WriteString(summary.Slug)
	} else if summary.ID != "" {
		b.WriteString("\n")
		b.WriteString("Thread ID: `")
		b.WriteString(summary.ID)
		b.WriteString("`")
	}

	message := b.String()
	if len(message) > 2000 {
		message = message[:1997] + "..."
	}

	return message
}

func formatLatestThreadMessage(baseURL url.URL, summary threadSummary) string {
	if summary.Title == "" {
		summary.Title = "Untitled thread"
	}

	var b strings.Builder
	b.WriteString("Latest thread:")
	b.WriteString("\n")
	b.WriteString("**")
	b.WriteString(sanitizeDiscordMessage(summary.Title))
	b.WriteString("**")

	if summary.Author != "" {
		b.WriteString(" by @")
		b.WriteString(sanitizeDiscordMessage(summary.Author))
	}

	if summary.Slug != "" {
		b.WriteString("\n")
		b.WriteString(strings.TrimRight(baseURL.String(), "/"))
		b.WriteString("/t/")
		b.WriteString(summary.Slug)
	} else if summary.ID != "" {
		b.WriteString("\n")
		b.WriteString("Thread ID: `")
		b.WriteString(summary.ID)
		b.WriteString("`")
	}

	return truncateDiscordMessage(b.String())
}

func truncateDiscordMessage(message string) string {
	if len(message) > 2000 {
		return message[:1997] + "..."
	}
	return message
}

func sanitizeDiscordMessage(input string) string {
	replacer := strings.NewReplacer(
		"@everyone", "@\u200beveryone",
		"@here", "@\u200bhere",
		"<@", "<\u200b@",
	)
	return replacer.Replace(input)
}
