"use client";

import { useQueryState } from "nuqs";
import { PropsWithChildren, createContext, useContext } from "react";

import { handle } from "@/api/client";
import { Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import {
  Editing,
  EditingSchema,
} from "@/components/site/SiteContextPane/useSiteContextPane";
import { FeedConfig } from "@/lib/settings/feed";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { DefaultFrontendConfig, Settings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";
import { hasPermission } from "@/utils/permissions";

type SettingsContextProps = {
  feed: FeedConfig;
  isEditingEnabled: boolean;
  isEditing: boolean;
  handleToggleEditing: () => void;
  updateFeed: (c: FeedConfig) => Promise<void>;
};

const context = createContext<SettingsContextProps | null>(null);

export function useSettingsContext(): SettingsContextProps {
  const value = useContext(context);
  if (!value) {
    // throw new Error(
    //   "useSettingsContext must be used within a SettingsContext provider",
    // );
    return {
      feed: DefaultFrontendConfig.feed,
      isEditingEnabled: false,
      isEditing: false,
      handleToggleEditing: () => {},
      updateFeed: async () => {},
    };
  }
  return value;
}

export function SettingsContext({ children }: PropsWithChildren<{}>) {
  const session = useSession();
  // const { settings } = useSettings();
  // const { updateSettings } = useSettingsMutation({});

  // const feed: FeedConfig = settings.metadata.feed;

  const [editing, setEditing] = useQueryState<null | Editing>("editing", {
    defaultValue: null,
    clearOnDefault: true,
    parse: EditingSchema.parse,
  });

  const isEditingEnabled = hasPermission(session, Permission.MANAGE_SETTINGS);

  const isEditing = editing === "feed";

  function handleToggleEditing() {
    if (editing) {
      setEditing(null);
    } else {
      setEditing("feed");
    }
  }

  const updateFeed = async (data: FeedConfig) => {
    await handle(
      async () => {
        // await updateSettings({
        //   metadata: {
        //     feed: data,
        //   },
        // });
      },
      {
        promiseToast: {
          loading: "Updating feed configuration...",
          success: "Updated!",
        },
      },
    );
  };

  return (
    <context.Provider
      value={{
        isEditingEnabled,
        isEditing,
        handleToggleEditing,
        feed: DefaultFrontendConfig.feed,
        updateFeed,
      }}
    >
      {children}
    </context.Provider>
  );
}
