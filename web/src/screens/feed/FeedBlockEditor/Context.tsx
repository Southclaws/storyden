"use client";

import { createContext, useContext } from "react";

import { type Account, AdminSettingsMutableProps } from "@/api/openapi-schema";
import {
  FeedBlock,
  FeedBlockType,
  FeedConfig,
  addFeedBlock,
  removeFeedBlock,
  reorderFeedBlock,
  replaceFeedBlock,
} from "@/lib/settings/feed";
import { useFeedMutation } from "@/lib/settings/feed-client";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { Settings } from "@/lib/settings/settings";

import { InitialData } from "../types";

type FeedBlockEditorContextValue = {
  feed: FeedConfig;
  initialData: InitialData;
  initialSession?: Account;
  initialSettings?: Settings;
  isEditing: boolean;
  addBlock: (type: FeedBlockType, index?: number) => Promise<void>;
  moveBlock: (active: FeedBlockType, over: FeedBlockType) => Promise<void>;
  overwriteBlock: (block: FeedBlock) => Promise<void>;
  removeBlock: (type: FeedBlockType) => Promise<void>;
  updateSite: (patch: AdminSettingsMutableProps) => Promise<void>;
};

const Context = createContext<FeedBlockEditorContextValue | null>(null);

export function useFeedBlockEditor() {
  const context = useContext(Context);
  if (!context) {
    throw new Error(
      "useFeedBlockEditor must be used within FeedBlockEditorProvider",
    );
  }
  return context;
}

type Props = {
  children: React.ReactNode;
  feed: FeedConfig;
  initialData: InitialData;
  initialSession?: Account;
  initialSettings?: Settings;
  isEditing: boolean;
};

export function FeedBlockEditorProvider({
  children,
  feed,
  initialData,
  initialSession,
  initialSettings,
  isEditing,
}: Props) {
  const { updateFeed } = useFeedMutation();
  const { updateSettings } = useSettingsMutation();

  async function addBlock(type: FeedBlockType, index?: number) {
    const updated = addFeedBlock(feed, type, index);
    if (updated === feed) return;
    await updateFeed(updated);
  }

  async function moveBlock(active: FeedBlockType, over: FeedBlockType) {
    const updated = reorderFeedBlock(feed, active, over);
    if (updated === feed) return;
    await updateFeed(updated);
  }

  async function overwriteBlock(updated: FeedBlock) {
    const next = replaceFeedBlock(feed, updated);
    if (next === feed) return;
    await updateFeed(next);
  }

  async function removeBlock(type: FeedBlockType) {
    const updated = removeFeedBlock(feed, type);
    if (updated === feed) return;
    await updateFeed(updated);
  }

  return (
    <Context.Provider
      value={{
        feed,
        initialData,
        initialSession,
        initialSettings,
        isEditing,
        addBlock,
        moveBlock,
        overwriteBlock,
        removeBlock,
        updateSite: updateSettings,
      }}
    >
      {children}
    </Context.Provider>
  );
}
