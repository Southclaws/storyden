import { keyBy } from "lodash";

import { DeleteIcon } from "@/components/ui/icons/Delete";
import * as Menu from "@/components/ui/menu";
import { allBlockTypes } from "@/lib/library/blockTypes";
import { useEmitLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock, LibraryPageBlockName } from "@/lib/library/metadata";

import { useWatch } from "../store";

import { CreateBlockMenu } from "./CreateBlockMenu";
import { LibraryPageAssetsBlockMenuItems } from "./LibraryPageAssetsBlock/LibraryPageAssetsBlockMenuItems";
import { LibraryPageCoverBlockMenuItems } from "./LibraryPageCoverBlock/LibraryPageCoverBlockMenuItems";
import { LibraryPageDirectoryBlockMenuItems } from "./LibraryPageDirectoryBlock/LibraryPageDirectoryBlockMenuItems";
import { LibraryPageTitleBlockMenuItems } from "./LibraryPageTitleBlock/LibraryPageTitleBlockMenuItems";

type Props = {
  block: LibraryPageBlock;
  index: number;
};

export function BlockMenu({ block, index }: Props) {
  const emit = useEmitLibraryBlockEvent();

  const currentMetadata = useWatch((s) => s.draft.meta);

  const existingBlocks = keyBy(currentMetadata.layout?.blocks, (b) => b.type);

  const blockList = allBlockTypes.filter((b) => !existingBlocks[b]);

  const newBlocksAvailable = blockList.length > 0;

  return (
    <Menu.ItemGroup>
      <Menu.ItemGroupLabel>
        <span>{LibraryPageBlockName[block.type]}</span>
      </Menu.ItemGroupLabel>

      <Menu.Separator />

      <Menu.Item
        value="delete"
        onClick={() =>
          emit("library:remove-block", {
            type: block.type,
          })
        }
      >
        <DeleteIcon />
        &nbsp;Delete
      </Menu.Item>
      <BlockConfigMenu index={index} block={block} />
      {newBlocksAvailable && <CreateBlockMenu />}
    </Menu.ItemGroup>
  );
}

function BlockConfigMenu({ block }: Props) {
  switch (block.type) {
    case "assets": {
      return <LibraryPageAssetsBlockMenuItems />;
    }
    case "cover": {
      return <LibraryPageCoverBlockMenuItems />;
    }
    case "directory": {
      return <LibraryPageDirectoryBlockMenuItems />;
    }
    case "title": {
      return <LibraryPageTitleBlockMenuItems />;
    }
    default:
      return null;
  }
}
