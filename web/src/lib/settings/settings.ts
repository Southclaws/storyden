import { z } from "zod";

import {
  AdminSettingsProps,
  AuthMode,
  Info,
  MessageOfTheDay,
} from "@/api/openapi-schema";
import { FALLBACK_COLOUR } from "@/utils/colour";

import { EditorSettingsSchema } from "./editor";
import { DefaultFeedConfig, FeedConfigSchema } from "./feed";
import { DefaultSidebarSettings, SidebarSettingsSchema } from "./sidebar";

export const DefaultEditorSettings = {
  mode: "richtext" as const,
} as const;

export const MotdAlertTypeSchema = z.enum([
  "celebration",
  "information",
  "alert",
]);
export type MotdAlertType = z.infer<typeof MotdAlertTypeSchema>;

export const MotdMetadataSchema = z
  .object({
    type: MotdAlertTypeSchema,
  })
  .passthrough();
export type MotdMetadata = z.infer<typeof MotdMetadataSchema>;

export const DefaultFrontendConfig = {
  feed: DefaultFeedConfig,
  editor: DefaultEditorSettings,
  sidebar: DefaultSidebarSettings,
  signatures: {
    enabled: true,
    maxHeight: 160,
  },
} as const;

export const DefaultSettings = {
  title: "Storyden",
  description: "A forum for the modern age.",
  content: "",
  accent_colour: FALLBACK_COLOUR,
  onboarding_status: "complete",
  authentication_mode: AuthMode.handle,
  capabilities: [],
  metadata: DefaultFrontendConfig,
} satisfies Settings;

// The frontend configuration is stored in Storyden's settings metadata field
// which is untyped on the backend so we provide a schema for validation here.
export const FrontendConfigurationSchema = z
  .object({
    feed: FeedConfigSchema,
    editor: EditorSettingsSchema.default(DefaultEditorSettings),
    sidebar: SidebarSettingsSchema.default(DefaultSidebarSettings),
    signatures: z
      .object({
        enabled: z.boolean().default(true),
        maxHeight: z.number().int().positive().max(2000).default(160),
      })
      .default({ enabled: true, maxHeight: 160 }),
  })
  .default(DefaultFrontendConfig);
export type FrontendConfiguration = z.infer<typeof FrontendConfigurationSchema>;
export type SignatureConfig = FrontendConfiguration["signatures"];

// Settings is the union of the backend typed config and the frontend config.
export type ParsedMotd = Omit<MessageOfTheDay, "metadata"> & {
  metadata?: MotdMetadata;
};
export type Settings = Omit<Info, "motd"> & {
  motd?: ParsedMotd;
  metadata: FrontendConfiguration;
};

// AdminSettings is non-public administration settings + parsed metadata config.
export type AdminSettings = Omit<AdminSettingsProps, "motd"> & {
  motd?: ParsedMotd;
  metadata: FrontendConfiguration;
};

export function parseSettings(data: Info): Settings {
  const parsed = FrontendConfigurationSchema.safeParse(data.metadata);
  const metadata = parsed.success ? parsed.data : DefaultFrontendConfig;
  if (parsed.error) {
    console.warn("Failed to parse frontend configuration:", parsed.error);
  }

  const settings = {
    ...data,
    metadata,
    motd: parseMotd(data.motd),
  } satisfies Settings;

  return settings;
}

export function parseAdminSettings(data: AdminSettingsProps): AdminSettings {
  const parsed = FrontendConfigurationSchema.safeParse(data.metadata);
  const metadata = parsed.success ? parsed.data : DefaultFrontendConfig;
  if (parsed.error) {
    console.warn("Failed to parse frontend configuration:", parsed.error);
  }

  const settings = {
    ...data,
    metadata,
    motd: parseMotd(data.motd),
  } satisfies AdminSettings;

  return settings;
}

function parseMotd(motd: MessageOfTheDay | undefined): ParsedMotd | undefined {
  if (!motd) return undefined;

  const parsed = MotdMetadataSchema.safeParse(motd.metadata);

  return {
    ...motd,
    metadata: parsed.success ? parsed.data : undefined,
  };
}
