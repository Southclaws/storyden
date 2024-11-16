import { z } from "zod";

import { AuthMode, Info } from "@/api/openapi-schema";
import { FALLBACK_COLOUR } from "@/utils/colour";

import { DefaultFeedConfig, FeedConfigSchema } from "./feed";

export const DefaultFrontendConfig = {
  feed: DefaultFeedConfig,
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
  })
  .default(DefaultFrontendConfig);
export type FrontendConfiguration = z.infer<typeof FrontendConfigurationSchema>;

// Settings is the union of the backend typed config and the frontend config.
export type Settings = Info & {
  metadata: FrontendConfiguration;
};

export function parseSettings(data: Info): Settings {
  const metadata = FrontendConfigurationSchema.parse(data.metadata);

  const settings = { ...data, metadata } satisfies Settings;

  return settings;
}
