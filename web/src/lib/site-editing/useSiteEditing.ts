"use client";

import { useQueryState } from "nuqs";
import z from "zod";

import { Account, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { hasPermission } from "@/utils/permissions";

export const EditingSchema = z.preprocess(
  (v) => {
    if (typeof v === "string" && v === "") {
      return undefined;
    }

    return v;
  },
  z.enum(["settings", "feed"]),
);
export type Editing = z.infer<typeof EditingSchema>;

export function useSiteEditing(session?: Account) {
  const currentSession = useSession(session);
  const [editing, setEditing] = useQueryState<null | Editing>("editing", {
    defaultValue: null,
    clearOnDefault: true,
    parse: EditingSchema.parse,
  });

  const isEditingEnabled = hasPermission(
    currentSession,
    Permission.MANAGE_SETTINGS,
  );

  const isEditingFeed = editing === "feed";
  const isEditingSettings = editing === "settings";

  function toggleFeedEditing() {
    if (editing === "feed") {
      setEditing(null);
    } else {
      setEditing("feed");
    }
  }

  function toggleSettingsEditing() {
    if (editing === "settings") {
      setEditing(null);
    } else {
      setEditing("settings");
    }
  }

  function stopEditing() {
    setEditing(null);
  }

  return {
    editing,
    isEditingEnabled,
    isEditingFeed,
    isEditingSettings,
    toggleFeedEditing,
    toggleSettingsEditing,
    stopEditing,
  };
}
