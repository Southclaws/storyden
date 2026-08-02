import { MenuSelectionDetails, Portal } from "@ark-ui/react";

import { LibraryPageSelect } from "@/components/library/LibraryPageSelect";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { LayoutIcon } from "@/components/ui/icons/Layout";
import { LayoutGridIcon } from "@/components/ui/icons/LayoutGrid";
import { LayoutListIcon } from "@/components/ui/icons/LayoutList";
import { SelectIcon } from "@/components/ui/icons/Select";
import * as Menu from "@/components/ui/menu";
import { FeedBlock, FeedBlockName } from "@/lib/settings/feed";
import { Box } from "@/styled-system/jsx";

import { useFeedBlockEditor } from "./Context";
import { CreateBlockMenu } from "./CreateBlockMenu";

type Props = {
  block: FeedBlock;
  children: React.ReactNode;
  index: number;
  open?: boolean;
};

export function BlockMenu({ block, children, index, open }: Props) {
  const { feed, removeBlock } = useFeedBlockEditor();
  const canAdd = feed.blocks.length < 8;

  function handleSelect({ value }: MenuSelectionDetails) {
    if (value === "delete") {
      void removeBlock(block.type);
    }
  }

  return (
    <Menu.Root
      lazyMount
      open={open}
      onSelect={handleSelect}
      positioning={{ placement: "right-start", gutter: 0 }}
    >
      <Menu.Trigger asChild>{children}</Menu.Trigger>
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="40">
            <Menu.ItemGroup>
              <Menu.ItemGroupLabel>
                {FeedBlockName[block.type]}
              </Menu.ItemGroupLabel>
              <Menu.Separator />
              <BlockConfigMenu block={block} />
              <Menu.Item value="delete">
                <DeleteIcon />
                Delete
              </Menu.Item>
              {canAdd && <CreateBlockMenu index={index} />}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

function BlockConfigMenu({ block }: { block: FeedBlock }) {
  switch (block.type) {
    case "categories":
      return <CategoryLayoutMenu block={block} />;
    case "threads":
      return <ThreadSourceMenu block={block} />;
    case "quick-share":
      return <QuickShareCategoryMenu block={block} />;
    case "library":
      return <LibraryConfigMenu block={block} />;
    default:
      return null;
  }
}

function CategoryLayoutMenu({
  block,
}: {
  block: Extract<FeedBlock, { type: "categories" }>;
}) {
  const { overwriteBlock } = useFeedBlockEditor();

  function handleSelect({ value }: MenuSelectionDetails) {
    void overwriteBlock({
      ...block,
      layout: value as "grid" | "list",
    });
  }

  return (
    <Menu.Root lazyMount onSelect={handleSelect}>
      <Menu.Trigger asChild>
        <Menu.Item value="layout">
          <LayoutIcon />
          Layout
        </Menu.Item>
      </Menu.Trigger>
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.Item value="list">
              <LayoutListIcon />
              List
            </Menu.Item>
            <Menu.Item value="grid">
              <LayoutGridIcon />
              Grid
            </Menu.Item>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

function ThreadSourceMenu({
  block,
}: {
  block: Extract<FeedBlock, { type: "threads" }>;
}) {
  const { overwriteBlock } = useFeedBlockEditor();

  function handleSelect({ value }: MenuSelectionDetails) {
    void overwriteBlock({
      ...block,
      source: value as "all" | "uncategorised",
    });
  }

  return (
    <Menu.Root lazyMount onSelect={handleSelect}>
      <Menu.Trigger asChild>
        <Menu.Item value="source">
          <SelectIcon />
          Source
        </Menu.Item>
      </Menu.Trigger>
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="40">
            <Menu.Item value="all">All threads</Menu.Item>
            <Menu.Item value="uncategorised">Uncategorised threads</Menu.Item>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

function QuickShareCategoryMenu({
  block,
}: {
  block: Extract<FeedBlock, { type: "quick-share" }>;
}) {
  const { overwriteBlock } = useFeedBlockEditor();

  function handleSelect({ value }: MenuSelectionDetails) {
    void overwriteBlock({
      ...block,
      showCategorySelect: value === "show",
    });
  }

  return (
    <Menu.Root lazyMount onSelect={handleSelect}>
      <Menu.Trigger asChild>
        <Menu.Item value="category-select">
          <SelectIcon />
          Category picker
        </Menu.Item>
      </Menu.Trigger>
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="40">
            <Menu.Item value="show">Show category picker</Menu.Item>
            <Menu.Item value="hide">Hide category picker</Menu.Item>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

function LibraryConfigMenu({
  block,
}: {
  block: Extract<FeedBlock, { type: "library" }>;
}) {
  const { overwriteBlock } = useFeedBlockEditor();

  return (
    <>
      <Menu.Root lazyMount>
        <Menu.Trigger asChild>
          <Menu.Item value="page">
            <SelectIcon />
            Page
          </Menu.Item>
        </Menu.Trigger>
        <Portal>
          <Menu.Positioner>
            <Menu.Content minW="64">
              <Box p="2">
                <LibraryPageSelect
                  value={block.node}
                  defaultValue={block.node}
                  onChange={(node) =>
                    void overwriteBlock({ ...block, node: node?.id })
                  }
                />
              </Box>
            </Menu.Content>
          </Menu.Positioner>
        </Portal>
      </Menu.Root>
      {!block.node && <LibraryLayoutMenu block={block} />}
    </>
  );
}

function LibraryLayoutMenu({
  block,
}: {
  block: Extract<FeedBlock, { type: "library" }>;
}) {
  const { overwriteBlock } = useFeedBlockEditor();

  function handleSelect({ value }: MenuSelectionDetails) {
    void overwriteBlock({
      ...block,
      layout: value as "grid" | "list",
    });
  }

  return (
    <Menu.Root lazyMount onSelect={handleSelect}>
      <Menu.Trigger asChild>
        <Menu.Item value="layout">
          <LayoutIcon />
          Layout
        </Menu.Item>
      </Menu.Trigger>
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.Item value="list">
              <LayoutListIcon />
              List
            </Menu.Item>
            <Menu.Item value="grid">
              <LayoutGridIcon />
              Grid
            </Menu.Item>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
