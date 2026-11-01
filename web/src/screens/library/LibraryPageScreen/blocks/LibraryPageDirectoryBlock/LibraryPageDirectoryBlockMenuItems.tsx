import { MenuSelectionDetails, Portal } from "@ark-ui/react";

import { CheckIcon } from "@/components/ui/icons/Check";
import { LayoutIcon } from "@/components/ui/icons/Layout";
import { LayoutGridIcon } from "@/components/ui/icons/LayoutGrid";
import { LayoutListIcon } from "@/components/ui/icons/LayoutList";
import { LayoutTableIcon } from "@/components/ui/icons/LayoutTable";
import { MediaIcon } from "@/components/ui/icons/Media";
import * as Menu from "@/components/ui/menu";
import {
  LibraryPageBlockTypeDirectoryConfig,
  LibraryPageBlockTypeDirectoryLayout,
} from "@/lib/library/metadata";

import { useLibraryPageContext } from "../../Context";

import { useDirectoryBlock } from "./useDirectoryBlock";

export function LibraryPageDirectoryBlockMenuItems() {
  return (
    <>
      <LayoutMenu />
      <CardStyleMenu />
    </>
  );
}

export function LayoutMenu() {
  const { store } = useLibraryPageContext();
  const block = useDirectoryBlock();

  if (block === undefined) {
    throw new Error("LayoutMenu rendered in a page without a Directory block.");
  }

  const { overwriteBlock } = store.getState();

  const handleSelect = ({ value }: MenuSelectionDetails) => {
    const config = block.config ?? {
      columns: [],
    };

    const newBlockConfig = {
      ...config,
      layout: value as LibraryPageBlockTypeDirectoryLayout,
    } satisfies LibraryPageBlockTypeDirectoryConfig;

    overwriteBlock({
      type: "directory",
      config: newBlockConfig,
    });
  };

  return (
    <Menu.Root lazyMount onSelect={handleSelect}>
      <Menu.TriggerItem>
        <LayoutIcon />
        &nbsp;Layout
      </Menu.TriggerItem>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.RadioItemGroup value={block.config?.layout ?? "table"}>
              <Menu.RadioItem value="table">
                <LayoutTableIcon />
                &nbsp;Table
                <Menu.ItemIndicator>
                  <CheckIcon />
                </Menu.ItemIndicator>
              </Menu.RadioItem>
              <Menu.RadioItem value="rows">
                <LayoutListIcon />
                &nbsp;Rows
                <Menu.ItemIndicator>
                  <CheckIcon />
                </Menu.ItemIndicator>
              </Menu.RadioItem>
              <Menu.RadioItem value="grid">
                <LayoutGridIcon />
                &nbsp;Grid
                <Menu.ItemIndicator>
                  <CheckIcon />
                </Menu.ItemIndicator>
              </Menu.RadioItem>
              <Menu.RadioItem value="gallery">
                <MediaIcon />
                &nbsp;Gallery
                <Menu.ItemIndicator>
                  <CheckIcon />
                </Menu.ItemIndicator>
              </Menu.RadioItem>
            </Menu.RadioItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

export function CardStyleMenu() {
  const { store } = useLibraryPageContext();
  const block = useDirectoryBlock();
  const config = block.config;

  if (!config || config.layout === "table" || config.layout === "rows") {
    return null;
  }

  const styles =
    config.layout === "gallery"
      ? (["minimal", "box"] as const)
      : (["minimal", "row", "box", "fill"] as const);
  const labels = {
    minimal: "Minimal",
    row: "Compact",
    box: "Editorial",
    fill: "Image overlay",
  };
  const savedStyle = config.card_style ?? "minimal";
  const selected =
    config.layout === "gallery" && savedStyle !== "minimal"
      ? "box"
      : savedStyle;
  return (
    <Menu.Root
      lazyMount
      onSelect={({ value }) => {
        if (!styles.some((style) => style === value)) {
          return;
        }
        store.getState().overwriteBlock({
          ...block,
          config: { ...config, card_style: value as typeof selected },
        });
      }}
    >
      <Menu.TriggerItem>
        <LayoutGridIcon />
        &nbsp;Card style
      </Menu.TriggerItem>
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="40">
            <Menu.RadioItemGroup value={selected}>
              {styles.map((style) => (
                <Menu.RadioItem key={style} value={style}>
                  {labels[style]}
                  <Menu.ItemIndicator>
                    <CheckIcon />
                  </Menu.ItemIndicator>
                </Menu.RadioItem>
              ))}
            </Menu.RadioItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
