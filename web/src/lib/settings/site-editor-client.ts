"use client";

import { useQueryState } from "nuqs";

import { type Account, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { type Editing, EditingSchema } from "@/components/site/editing";
import { hasPermission } from "@/utils/permissions";

import { type Settings } from "./settings";

type SiteEditorStateOptions = {
  initialSession?: Account;
  initialSettings?: Settings;
};

export function useSiteEditorState({
  initialSession,
  initialSettings,
}: SiteEditorStateOptions = {}) {
  const session = useSession(initialSession, initialSettings);
  const [editing, setEditing] = useQueryState<null | Editing>("editing", {
    defaultValue: null,
    clearOnDefault: true,
    parse: EditingSchema.parse,
  });

  const isEditingEnabled = hasPermission(session, Permission.MANAGE_SETTINGS);
  const isEditingRequested = editing === "site";
  const isEditing = isEditingEnabled && isEditingRequested;

  function handleToggleEditing() {
    void setEditing(isEditingRequested ? null : "site");
  }

  return {
    session,
    isEditingEnabled,
    isEditing,
    handleToggleEditing,
  };
}
