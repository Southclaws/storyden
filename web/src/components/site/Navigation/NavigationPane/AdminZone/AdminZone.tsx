"use client";

import { type Account } from "@/api/openapi-schema";
import { IconButton } from "@/components/ui/icon-button";
import { EditIcon } from "@/components/ui/icons/Edit";
import * as Tooltip from "@/components/ui/tooltip";
import { type Settings } from "@/lib/settings/settings";
import { useSiteEditorState } from "@/lib/settings/site-editor-client";

import { AdminAnchor } from "../../Anchors/Admin";
import { Route, useRoute } from "../../useRoute";

const editableRoute: Record<Route["name"], boolean> = {
  index: true,
  library: false,
  admin: false,
  settings: false,
};

type Props = {
  initialSession?: Account;
  initialSettings?: Settings;
};

export function AdminZone({ initialSession, initialSettings }: Props) {
  const { isEditingEnabled, isEditing, handleToggleEditing } =
    useSiteEditorState({ initialSession, initialSettings });
  const route = useRoute();
  const isRouteEditable = route && editableRoute[route.name];

  if (!isEditingEnabled) {
    return null;
  }

  const label = isRouteEditable
    ? isEditing
      ? "Exit edit mode"
      : "Edit site"
    : "Admin";

  return (
    <div className="navigation-admin-zone__trigger">
      <Tooltip.Root openDelay={300}>
        <Tooltip.Trigger asChild>
          <span className="navigation-pane__tooltip-anchor">
            {isRouteEditable ? (
              <IconButton
                aria-label={label}
                aria-pressed={isEditing}
                size="sm"
                type="button"
                variant={isEditing ? "subtle" : "ghost"}
                onClick={handleToggleEditing}
              >
                <EditIcon />
              </IconButton>
            ) : (
              <AdminAnchor hideLabel size="sm" />
            )}
          </span>
        </Tooltip.Trigger>
        <Tooltip.Positioner>
          <Tooltip.Content>{label}</Tooltip.Content>
        </Tooltip.Positioner>
      </Tooltip.Root>
    </div>
  );
}
