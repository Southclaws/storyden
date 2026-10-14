"use client";

import dynamic from "next/dynamic";

import { type Account } from "@/api/openapi-schema";
import { useFeedConfig } from "@/lib/settings/feed-client";
import { type Settings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";
import { useSiteEditorState } from "@/lib/settings/site-editor-client";
import { useWebMCPToolsEnabled } from "@/lib/webmcp/activation";

import { FeedBlockEditorProvider } from "./FeedBlockEditor/Context";
import { FeedPageBlocks } from "./FeedBlockEditor/FeedPageBlocks";
import { InitialData } from "./types";

const IndexPageWebMCPTools = dynamic(
  () =>
    import("./FeedBlockEditor/IndexPageWebMCPTools").then(
      (module) => module.IndexPageWebMCPTools,
    ),
  { ssr: false },
);

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
  const webMCPEnabled = useWebMCPToolsEnabled();
  const { settings } = useSettings(initialSettings, true);
  const currentSettings = settings ?? initialSettings;
  const feed = useFeedConfig(currentSettings, false);
  const { isEditingEnabled, isEditing, handleToggleEditing } =
    useSiteEditorState({
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
      {webMCPEnabled && (
        <IndexPageWebMCPTools
          isEditingEnabled={isEditingEnabled}
          isEditing={isEditing}
          startEditing={handleToggleEditing}
        />
      )}
      <FeedPageBlocks />
    </FeedBlockEditorProvider>
  );
}
