"use client";

import { type Account } from "@/api/openapi-schema";
import { useFeedConfig } from "@/lib/settings/feed-client";
import { type Settings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";
import { useSiteEditorState } from "@/lib/settings/site-editor-client";

import { FeedBlockEditorProvider } from "./FeedBlockEditor/Context";
import { FeedPageBlocks } from "./FeedBlockEditor/FeedPageBlocks";
import { InitialData } from "./types";

type Props = {
  initialData: InitialData;
  initialSettings?: Settings;
  initialSession?: Account;
};

export function FeedScreenContent({
  initialData,
  initialSettings,
  initialSession,
}: Props) {
  const { settings } = useSettings(initialSettings, true);
  const currentSettings = settings ?? initialSettings;
  const feed = useFeedConfig(currentSettings, false);
  const { isEditing } = useSiteEditorState({
    initialSession,
    initialSettings: currentSettings,
  });

  return (
    <FeedBlockEditorProvider
      feed={feed}
      initialData={initialData}
      initialSession={initialSession}
      initialSettings={currentSettings}
      isEditing={isEditing}
    >
      <FeedPageBlocks />
    </FeedBlockEditorProvider>
  );
}
