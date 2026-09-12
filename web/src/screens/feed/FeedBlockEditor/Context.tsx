"use client";

import { createContext, useContext, useEffect, useRef } from "react";

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

export type FeedBlockEditorActions = {
  getFeed: () => FeedConfig;
  addBlock: (type: FeedBlockType, index?: number) => Promise<FeedConfig>;
  moveBlock: (
    active: FeedBlockType,
    over: FeedBlockType,
  ) => Promise<FeedConfig>;
  overwriteBlock: (block: FeedBlock) => Promise<FeedConfig>;
  removeBlock: (type: FeedBlockType) => Promise<FeedConfig>;
  updateSite: (patch: AdminSettingsMutableProps) => Promise<void>;
};

type FeedBlockEditorContextValue = FeedBlockEditorActions & {
  feed: FeedConfig;
  initialData: InitialData;
  initialSession?: Account;
  initialSettings?: Settings;
  isEditing: boolean;
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
  const feedRef = useRef(feed);
  useEffect(() => {
    feedRef.current = feed;
  }, [feed]);

  function getFeed() {
    return feedRef.current;
  }

  async function persistFeed(updated: FeedConfig, current: FeedConfig) {
    if (updated === current) return current;

    feedRef.current = updated;
    try {
      await updateFeed(updated);
      return updated;
    } catch (error) {
      feedRef.current = current;
      throw error;
    }
  }

  async function addBlock(type: FeedBlockType, index?: number) {
    const current = getFeed();
    return persistFeed(addFeedBlock(current, type, index), current);
  }

  async function moveBlock(active: FeedBlockType, over: FeedBlockType) {
    const current = getFeed();
    return persistFeed(reorderFeedBlock(current, active, over), current);
  }

  async function overwriteBlock(updated: FeedBlock) {
    const current = getFeed();
    return persistFeed(replaceFeedBlock(current, updated), current);
  }

  async function removeBlock(type: FeedBlockType) {
    const current = getFeed();
    return persistFeed(removeFeedBlock(current, type), current);
  }

  return (
    <Context.Provider
      value={{
        feed,
        initialData,
        initialSession,
        initialSettings,
        isEditing,
        getFeed,
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
