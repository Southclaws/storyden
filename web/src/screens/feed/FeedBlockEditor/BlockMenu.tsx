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
  index: number;
  onConfigured: () => void;
};

export function BlockMenu({ block, index, onConfigured }: Props) {
  const { feed, removeBlock } = useFeedBlockEditor();
  const canAdd = feed.blocks.length < 8;

  return (
    <Menu.ItemGroup>
      <Menu.ItemGroupLabel>{FeedBlockName[block.type]}</Menu.ItemGroupLabel>
      <Menu.Separator />
      <BlockConfigMenu block={block} onConfigured={onConfigured} />
      <Menu.Item value="delete" onClick={() => void removeBlock(block.type)}>
        <DeleteIcon />
        Delete
      </Menu.Item>
      {canAdd && <CreateBlockMenu index={index} />}
    </Menu.ItemGroup>
  );
}

function BlockConfigMenu({
  block,
  onConfigured,
}: {
  block: FeedBlock;
  onConfigured: () => void;
}) {
  switch (block.type) {
    case "categories":
      return <CategoryLayoutMenu block={block} onConfigured={onConfigured} />;
    case "threads":
      return <ThreadSourceMenu block={block} onConfigured={onConfigured} />;
    case "quick-share":
      return (
        <QuickShareCategoryMenu block={block} onConfigured={onConfigured} />
      );
    case "library":
      return <LibraryConfigMenu block={block} onConfigured={onConfigured} />;
    default:
      return null;
  }
}

function CategoryLayoutMenu({
  block,
  onConfigured,
}: {
  block: Extract<FeedBlock, { type: "categories" }>;
  onConfigured: () => void;
}) {
  const { overwriteBlock } = useFeedBlockEditor();

  function handleSelect({ value }: MenuSelectionDetails) {
    void overwriteBlock({
      ...block,
      layout: value as "grid" | "list",
    });
    onConfigured();
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
  onConfigured,
}: {
  block: Extract<FeedBlock, { type: "threads" }>;
  onConfigured: () => void;
}) {
  const { overwriteBlock } = useFeedBlockEditor();

  function handleSelect({ value }: MenuSelectionDetails) {
    void overwriteBlock({
      ...block,
      source: value as "all" | "uncategorised",
    });
    onConfigured();
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
  onConfigured,
}: {
  block: Extract<FeedBlock, { type: "quick-share" }>;
  onConfigured: () => void;
}) {
  const { overwriteBlock } = useFeedBlockEditor();

  function handleSelect({ value }: MenuSelectionDetails) {
    void overwriteBlock({
      ...block,
      showCategorySelect: value === "show",
    });
    onConfigured();
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
  onConfigured,
}: {
  block: Extract<FeedBlock, { type: "library" }>;
  onConfigured: () => void;
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
                    void overwriteBlock({ ...block, node: node?.id }).then(
                      onConfigured,
                    )
                  }
                />
              </Box>
            </Menu.Content>
          </Menu.Positioner>
        </Portal>
      </Menu.Root>
      {!block.node && (
        <LibraryLayoutMenu block={block} onConfigured={onConfigured} />
      )}
    </>
  );
}

function LibraryLayoutMenu({
  block,
  onConfigured,
}: {
  block: Extract<FeedBlock, { type: "library" }>;
  onConfigured: () => void;
}) {
  const { overwriteBlock } = useFeedBlockEditor();

  function handleSelect({ value }: MenuSelectionDetails) {
    void overwriteBlock({
      ...block,
      layout: value as "grid" | "list",
    });
    onConfigured();
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
