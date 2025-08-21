import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { LayoutIcon, Table2Icon } from "lucide-react";

import { LayoutGridIcon } from "@/components/ui/icons/LayoutGrid";
import * as Menu from "@/components/ui/menu";
import {
  LibraryPageBlockTypeDirectoryConfig,
  LibraryPageBlockTypeDirectoryLayout,
} from "@/lib/library/metadata";

import { useLibraryPageContext } from "../../Context";
import { useBlock } from "../useBlock";

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
    const newBlockConfig = {
      ...block.config!, // TODO: Do something more clever?
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
              <Table2Icon />
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
