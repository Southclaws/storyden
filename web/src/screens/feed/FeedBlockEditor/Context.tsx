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
  const mutationQueueRef = useRef<Promise<void> | undefined>(undefined);
  useEffect(() => {
    feedRef.current = feed;
  }, [feed]);

  function getFeed() {
    return feedRef.current;
  }

  function mutateFeed(mutation: (current: FeedConfig) => FeedConfig) {
    const result = (mutationQueueRef.current ?? Promise.resolve()).then(
      async () => {
        const current = feedRef.current;
        const updated = mutation(current);
        if (updated === current) return current;

        feedRef.current = updated;
        try {
          await updateFeed(updated);
          return updated;
        } catch (error) {
          feedRef.current = current;
          throw error;
        }
      },
    );
    mutationQueueRef.current = result.then(
      () => undefined,
      () => undefined,
    );
    return result;
  }

  async function addBlock(type: FeedBlockType, index?: number) {
    return mutateFeed((current) => addFeedBlock(current, type, index));
  }

  async function moveBlock(active: FeedBlockType, over: FeedBlockType) {
    return mutateFeed((current) => reorderFeedBlock(current, active, over));
  }

  async function overwriteBlock(updated: FeedBlock) {
    return mutateFeed((current) => replaceFeedBlock(current, updated));
  }

  async function removeBlock(type: FeedBlockType) {
    return mutateFeed((current) => removeFeedBlock(current, type));
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
