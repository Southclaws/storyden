"use client";

import { handle } from "@/api/client";

import { DefaultFeedConfig, type FeedConfig } from "./feed";
import { useSettingsMutation } from "./mutation";
import { type Settings } from "./settings";
import { useSettings } from "./settings-client";

export function useFeedConfig(
  initialSettings?: Settings,
  revalidateOnMount = false,
): FeedConfig {
  const { settings } = useSettings(initialSettings, revalidateOnMount);

  return (
    settings?.metadata.feed ??
    initialSettings?.metadata.feed ??
    DefaultFeedConfig
  );
}

export function useFeedMutation() {
  const { updateSettings } = useSettingsMutation();

  const updateFeed = async (feed: FeedConfig) => {
    await handle(
      async () => {
        await updateSettings({
          metadata: {
            feed,
          },
        });
      },
      {
        promiseToast: {
          loading: "Updating feed configuration...",
          success: "Updated!",
        },
      },
    );
  };

  return {
    updateFeed,
  };
}
