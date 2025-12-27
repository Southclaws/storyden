"use client";

import { usePathname } from "next/navigation";
import { useQueryState } from "nuqs";
import { PropsWithChildren, createContext, useContext } from "react";

import { handle } from "@/api/client";
import { Account, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import {
  Editing,
  EditingSchema,
} from "@/components/site/SiteContextPane/useSiteContextPane";
import { FeedConfig } from "@/lib/settings/feed";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { Settings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";
import { hasPermission } from "@/utils/permissions";

type SettingsContextProps = {
  session: Account | undefined;
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
    throw new Error(
      "useSettingsContext must be used within a SettingsContext provider",
    );
  }
  return value;
}

export function SettingsContext({
  initialSession,
  initialSettings,
  children,
}: PropsWithChildren<{
  initialSession: Account | undefined;
  initialSettings: Settings;
}>) {
  const session = useSession(initialSession);
  const { settings } = useSettings(initialSettings);
  const { updateSettings } = useSettingsMutation();

  const feed: FeedConfig = (settings ?? initialSettings).metadata.feed;

  const [editing, setEditing] = useQueryState<null | Editing>("editing", {
    defaultValue: null,
    clearOnDefault: true,
    parse: EditingSchema.parse,
  });

  const isEditingEnabled = hasPermission(
    initialSession,
    Permission.MANAGE_SETTINGS,
  );

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
        await updateSettings({
          metadata: {
            feed: data,
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

  return (
    <context.Provider
      value={{
        session,
        isEditingEnabled,
        isEditing,
        handleToggleEditing,
        feed,
        updateFeed,
      }}
    >
      {children}
    </context.Provider>
  );
}
