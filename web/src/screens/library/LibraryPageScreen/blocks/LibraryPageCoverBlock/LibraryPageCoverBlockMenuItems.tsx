import { Portal } from "@ark-ui/react";

import { LayoutIcon } from "@/components/ui/icons/Layout";
import { MediaRemoveIcon } from "@/components/ui/icons/Media";
import * as Menu from "@/components/ui/menu";
import { LibraryPageBlockTypeCover } from "@/lib/library/metadata";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useBlock } from "../useBlock";

export function LibraryPageCoverBlockMenuItems() {
  const { store } = useLibraryPageContext();
  const primary = useWatch((s) => s.draft.primary_image);
  const block = useBlock("cover");
  const { overwriteBlock, removePrimaryImage } = store.getState();

  return (
    <>
      <Menu.Root
        onSelect={({ value }) => {
          if (block)
            overwriteBlock({
              ...block,
              config: {
                ...block.config,
                layout: value as NonNullable<
                  LibraryPageBlockTypeCover["config"]
                >["layout"],
              },
            });
        }}
      >
        <Menu.TriggerItem>
          <LayoutIcon />
          Page layout
        </Menu.TriggerItem>
        <Portal>
          <Menu.Positioner>
            <Menu.Content minW="64">
              <Menu.Item value="banner">Banner · wide page header</Menu.Item>
              <Menu.Item value="image">Image · natural proportions</Menu.Item>
              <Menu.Item value="artwork">
                Artwork · centered logo or album cover
              </Menu.Item>
            </Menu.Content>
          </Menu.Positioner>
        </Portal>
      </Menu.Root>
      {primary?.mime_type === "image/gif" && (
        <Menu.ItemGroupLabel>
          Cropping saves a still image; the original animation is retained.
        </Menu.ItemGroupLabel>
      )}
      {primary && (
        <Menu.Item value="remove-cover" onSelect={() => removePrimaryImage()}>
          <MediaRemoveIcon />
          Remove main image
        </Menu.Item>
      )}
    </>
  );
}
