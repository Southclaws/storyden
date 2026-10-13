import {
  AdminSettingsProps,
  AuthMode,
  Info,
  MessageOfTheDay,
  RegistrationMode,
} from "@/api/openapi-schema";
import * as z from "@/lib/validation";
import { FALLBACK_COLOUR } from "@/utils/colour";

import { EditorSettingsSchema } from "./editor";
import { DefaultFeedConfig, FeedConfigSchema } from "./feed";
import { DefaultNavigationConfig, NavigationConfigSchema } from "./navigation";

export const DefaultEditorSettings = {
  mode: "richtext" as const,
} as const;

export const MotdAlertTypeSchema = z.enum([
  "celebration",
  "information",
  "alert",
]);
export type MotdAlertType = z.infer<typeof MotdAlertTypeSchema>;

export const MotdMetadataSchema = z.looseObject({
  type: MotdAlertTypeSchema,
});
export type MotdMetadata = z.infer<typeof MotdMetadataSchema>;

export const DefaultFrontendConfig = {
  feed: DefaultFeedConfig,
  navigation: DefaultNavigationConfig,
  editor: DefaultEditorSettings,
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
  registration_mode: RegistrationMode.public,
  capabilities: [],
  metadata: DefaultFrontendConfig,
  api_address: "http://localhost:8000",
  web_address: "http://localhost:3000",
} satisfies Settings;

// The frontend configuration is stored in Storyden's settings metadata field
// which is untyped on the backend so we provide a schema for validation here.
export const FrontendConfigurationSchema = z._default(
  z.object({
    feed: FeedConfigSchema,
    navigation: NavigationConfigSchema,
    editor: z._default(EditorSettingsSchema, DefaultEditorSettings),
    signatures: z._default(
      z.object({
        enabled: z._default(z.boolean(), true),
        maxHeight: z._default(
          z.number().check(z.int(), z.positive(), z.maximum(2000)),
          160,
        ),
      }),
      { enabled: true, maxHeight: 160 },
    ),
  }),
  DefaultFrontendConfig,
);
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
