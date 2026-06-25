package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/sdk/go/storyden"
	"github.com/bwmarrin/discordgo"
)

const (
	providerID = "discord"

	configDiscordToken = "discord_token"
	configChannelID    = "channel_id"

	toolListChannels = "list_channels"
	toolSendMessage  = "send_message"

	discordTimeout = 15 * time.Second
	maxMessageLen  = 2000
)

type pluginConfig struct {
	DiscordToken string
	ChannelID    string
}

type discordTools struct {
	plugin *storyden.Plugin
	logger *slog.Logger
}

type channelInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	ParentID string `json:"parent_id,omitempty"`
	Position int    `json:"position"`
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	plugin, err := storyden.New(ctx)
	if err != nil {
		exitError(logger, "create plugin", err)
	}
	defer func() {
		if err := plugin.Shutdown(); err != nil && !errors.Is(err, context.Canceled) {
			logger.Warn("plugin shutdown returned error", slog.String("error", err.Error()))
		}
	}()

	dt := &discordTools{plugin: plugin, logger: logger}
	plugin.OnRobotToolCall(dt.handleRobotToolCall)

	logger.Info("starting Discord Robot tools plugin")
	if err := plugin.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		exitError(logger, "plugin runtime", err)
	}
}

func (dt *discordTools) handleRobotToolCall(ctx context.Context, req rpc.RPCRequestRobotToolCallParams) (rpc.RPCResponseRobotToolCall, error) {
	if req.ProviderID != providerID {
		return toolError(fmt.Sprintf("unsupported provider %q", req.ProviderID)), nil
	}

	callCtx, cancel := context.WithTimeout(ctx, discordTimeout)
	defer cancel()

	switch req.ToolID {
	case toolListChannels:
		return dt.listChannels(callCtx, req.Arguments), nil
	case toolSendMessage:
		return dt.sendMessage(callCtx, req.Arguments), nil
	default:
		return toolError(fmt.Sprintf("unsupported tool %q", req.ToolID)), nil
	}
}

func (dt *discordTools) listChannels(ctx context.Context, args map[string]interface{}) rpc.RPCResponseRobotToolCall {
	cfg, dg, guild, err := dt.discord(ctx)
	if err != nil {
		return toolError(err.Error())
	}

	channels, err := dg.GuildChannels(guild.ID, discordgo.WithContext(ctx))
	if err != nil {
		return toolError(fmt.Sprintf("list Discord channels: %v", err))
	}

	out := make([]channelInfo, 0, len(channels))
	for _, channel := range channels {
		if shouldExposeChannel(channel.Type) {
			out = append(out, channelInfoFromDiscord(channel))
		}
	}

	if includeActiveThreads(args) {
		threads, err := dg.GuildThreadsActive(guild.ID, discordgo.WithContext(ctx))
		if err != nil {
			return toolError(fmt.Sprintf("list Discord active threads: %v", err))
		}
		for _, thread := range threads.Threads {
			out = append(out, channelInfoFromDiscord(thread))
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Position == out[j].Position {
			return out[i].Name < out[j].Name
		}
		return out[i].Position < out[j].Position
	})

	dt.logger.Info("listed Discord channels",
		slog.String("guild_id", guild.ID),
		slog.String("configured_channel_id", cfg.ChannelID),
		slog.Int("channels", len(out)))

	return toolOutput(map[string]interface{}{
		"guild_id":   guild.ID,
		"guild_name": guild.Name,
		"channels":   out,
	})
}

func (dt *discordTools) sendMessage(ctx context.Context, args map[string]interface{}) rpc.RPCResponseRobotToolCall {
	channelID, ok := stringArg(args, "channel_id")
	if !ok {
		return toolError("channel_id is required")
	}
	content, ok := stringArg(args, "message")
	if !ok {
		return toolError("message is required")
	}
	if len(content) > maxMessageLen {
		return toolError(fmt.Sprintf("message is too long: Discord limit is %d characters", maxMessageLen))
	}

	_, dg, guild, err := dt.discord(ctx)
	if err != nil {
		return toolError(err.Error())
	}

	msg, err := dg.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: content,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{},
		},
	}, discordgo.WithContext(ctx))
	if err != nil {
		return toolError(fmt.Sprintf("send Discord message: %v", err))
	}

	dt.logger.Info("sent Discord message",
		slog.String("channel_id", channelID),
		slog.String("message_id", msg.ID))

	return toolOutput(map[string]interface{}{
		"channel_id":  channelID,
		"message_id":  msg.ID,
		"content":     msg.Content,
		"discord_url": discordMessageURL(guild.ID, msg.ChannelID, msg.ID),
	})
}

func (dt *discordTools) discord(ctx context.Context) (pluginConfig, *discordgo.Session, *discordgo.Guild, error) {
	cfg, err := dt.config(ctx)
	if err != nil {
		return pluginConfig{}, nil, nil, err
	}

	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return pluginConfig{}, nil, nil, fmt.Errorf("create Discord session: %w", err)
	}
	dg.Client = &http.Client{Timeout: discordTimeout}
	dg.Identify.Intents = discordgo.IntentsGuilds

	channel, err := dg.Channel(cfg.ChannelID, discordgo.WithContext(ctx))
	if err != nil {
		return pluginConfig{}, nil, nil, fmt.Errorf("load configured Discord channel: %w", err)
	}
	if strings.TrimSpace(channel.GuildID) == "" {
		return pluginConfig{}, nil, nil, fmt.Errorf("configured Discord channel %q is not a guild channel", cfg.ChannelID)
	}

	guild, err := dg.Guild(channel.GuildID, discordgo.WithContext(ctx))
	if err != nil {
		return pluginConfig{}, nil, nil, fmt.Errorf("load Discord guild: %w", err)
	}

	return cfg, dg, guild, nil
}

func (dt *discordTools) config(ctx context.Context) (pluginConfig, error) {
	config, err := dt.plugin.GetConfig(ctx, configDiscordToken, configChannelID)
	if err != nil {
		return pluginConfig{}, fmt.Errorf("get plugin config: %w", err)
	}
	return parseConfig(config)
}

func parseConfig(in map[string]any) (pluginConfig, error) {
	token, ok := stringArg(in, configDiscordToken)
	if !ok {
		return pluginConfig{}, fmt.Errorf("%s is required", configDiscordToken)
	}
	channelID, ok := stringArg(in, configChannelID)
	if !ok {
		return pluginConfig{}, fmt.Errorf("%s is required", configChannelID)
	}
	return pluginConfig{DiscordToken: token, ChannelID: channelID}, nil
}

func shouldExposeChannel(kind discordgo.ChannelType) bool {
	switch kind {
	case discordgo.ChannelTypeGuildText,
		discordgo.ChannelTypeGuildNews,
		discordgo.ChannelTypeGuildForum,
		discordgo.ChannelTypeGuildMedia,
		discordgo.ChannelTypeGuildPublicThread,
		discordgo.ChannelTypeGuildPrivateThread,
		discordgo.ChannelTypeGuildNewsThread:
		return true
	default:
		return false
	}
}

func channelInfoFromDiscord(channel *discordgo.Channel) channelInfo {
	return channelInfo{
		ID:       channel.ID,
		Name:     channel.Name,
		Type:     channelTypeName(channel.Type),
		ParentID: channel.ParentID,
		Position: channel.Position,
	}
}

func channelTypeName(kind discordgo.ChannelType) string {
	switch kind {
	case discordgo.ChannelTypeGuildText:
		return "text"
	case discordgo.ChannelTypeGuildNews:
		return "news"
	case discordgo.ChannelTypeGuildForum:
		return "forum"
	case discordgo.ChannelTypeGuildMedia:
		return "media"
	case discordgo.ChannelTypeGuildPublicThread:
		return "public_thread"
	case discordgo.ChannelTypeGuildPrivateThread:
		return "private_thread"
	case discordgo.ChannelTypeGuildNewsThread:
		return "news_thread"
	default:
		return fmt.Sprintf("discord_channel_type_%d", kind)
	}
}

func includeActiveThreads(args map[string]interface{}) bool {
	raw, ok := args["include_active_threads"]
	if !ok {
		return false
	}
	value, ok := raw.(bool)
	return ok && value
}

func stringArg(in map[string]interface{}, key string) (string, bool) {
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
	return value, value != ""
}

func discordMessageURL(guildID, channelID, messageID string) string {
	if guildID == "" {
		guildID = "@me"
	}
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, channelID, messageID)
}

func toolOutput(output map[string]interface{}) rpc.RPCResponseRobotToolCall {
	return rpc.RPCResponseRobotToolCall{
		Method: "robot_tool_call",
		Output: output,
	}
}

func toolError(message string) rpc.RPCResponseRobotToolCall {
	return rpc.RPCResponseRobotToolCall{
		Method: "robot_tool_call",
		Error:  opt.New(message),
	}
}

func exitError(logger *slog.Logger, action string, err error) {
	logger.Error(action+" failed", slog.String("error", err.Error()))
	os.Exit(1)
}
