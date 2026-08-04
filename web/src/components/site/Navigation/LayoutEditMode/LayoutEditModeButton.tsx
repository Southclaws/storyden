"use client";

import { type Account } from "@/api/openapi-schema";
import { IconButton } from "@/components/ui/icon-button";
import { EditIcon } from "@/components/ui/icons/Edit";
import * as Tooltip from "@/components/ui/tooltip";
import { type Settings } from "@/lib/settings/settings";
import { useSiteEditorState } from "@/lib/settings/site-editor-client";

type Props = {
  initialSession?: Account;
  initialSettings?: Settings;
};

export function LayoutEditModeButton({
  initialSession,
  initialSettings,
}: Props) {
  const { isEditingEnabled, isEditing, handleToggleEditing } =
    useSiteEditorState({ initialSession, initialSettings });

  if (!isEditingEnabled) {
    return null;
  }

  const label = isEditing ? "Exit edit mode" : "Edit site";

  return (
    <div className="navigation-layout-edit-mode__trigger">
      <Tooltip.Root openDelay={300}>
        <Tooltip.Trigger asChild>
          <span className="navigation-pane__tooltip-anchor">
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
          </span>
        </Tooltip.Trigger>
        <Tooltip.Positioner>
          <Tooltip.Content>{label}</Tooltip.Content>
        </Tooltip.Positioner>
      </Tooltip.Root>
    </div>
  );
}
