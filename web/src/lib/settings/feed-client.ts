"use client";

import { useQueryState } from "nuqs";

import { handle } from "@/api/client";
import { type Account, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { type Editing, EditingSchema } from "@/components/site/editing";
import { hasPermission } from "@/utils/permissions";

import { DefaultFeedConfig, type FeedConfig } from "./feed";
import { useSettingsMutation } from "./mutation";
import { type Settings } from "./settings";
import { useSettings } from "./settings-client";

type FeedEditorStateOptions = {
  initialSession?: Account;
  initialSettings?: Settings;
};

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

export function useFeedEditorState({
  initialSession,
  initialSettings,
}: FeedEditorStateOptions = {}) {
  const session = useSession(initialSession, initialSettings);
  const [editing, setEditing] = useQueryState<null | Editing>("editing", {
    defaultValue: null,
    clearOnDefault: true,
    parse: EditingSchema.parse,
  });

  const isEditingEnabled = hasPermission(session, Permission.MANAGE_SETTINGS);
  const isEditing = editing === "feed";

  function handleToggleEditing() {
    if (editing === "feed") {
      setEditing(null);
    } else {
      setEditing("feed");
    }
  }

  return {
    session,
    isEditingEnabled,
    isEditing,
    handleToggleEditing,
  };
}
