import { MenuSelectionDetails, Portal } from "@ark-ui/react";

import { LayoutIcon } from "@/components/ui/icons/Layout";
import { LayoutGridIcon } from "@/components/ui/icons/LayoutGrid";
import { LayoutTableIcon } from "@/components/ui/icons/LayoutTable";
import * as Menu from "@/components/ui/menu";
import {
  LibraryPageBlockTypeDirectoryConfig,
  LibraryPageBlockTypeDirectoryLayout,
} from "@/lib/library/metadata";

import { useLibraryPageContext } from "../../Context";

import { getDefaultBlockConfig } from "./column";
import { useDirectoryBlock } from "./useDirectoryBlock";

export function LibraryPageDirectoryBlockMenuItems() {
  return (
    <>
      <LayoutMenu />
    </>
  );
}

function LayoutMenu() {
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
      <Menu.Trigger asChild>
        <Menu.Item value="layout">
          <LayoutIcon />
          &nbsp;Layout
        </Menu.Item>
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.Item value="table">
              <LayoutTableIcon />
              &nbsp;Table
            </Menu.Item>
            <Menu.Item value="grid">
              <LayoutGridIcon />
              &nbsp;Grid
            </Menu.Item>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
