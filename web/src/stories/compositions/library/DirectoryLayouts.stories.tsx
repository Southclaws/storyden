import {
  DndContext,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  rectSortingStrategy,
  useSortable,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useRouter } from "@storybook/nextjs-vite/navigation.mock";
import { useState } from "react";
import { fn } from "storybook/test";

import { NodeWithChildren } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { IconButton } from "@/components/ui/icon-button";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import { Text } from "@/components/ui/text";
import { LibraryPageBlockTypeDirectoryConfig } from "@/lib/library/metadata";
import { DirectoryRichCard } from "@/screens/library/LibraryPageScreen/blocks/LibraryPageDirectoryBlock/DirectoryRichCard";
import {
  mergeFieldsAndProperties,
  mergeFieldsAndPropertySchema,
} from "@/screens/library/LibraryPageScreen/blocks/LibraryPageDirectoryBlock/column";
import { css } from "@/styled-system/css";
import { Box, HStack, LStack, styled } from "@/styled-system/jsx";

const image = new URL("../../assets/mountains.png", import.meta.url).href;
const schema = [
  { fid: "status", name: "Status", type: "text", sort: "0" },
  { fid: "type", name: "Type", type: "text", sort: "1" },
] as const;
const initialNodes: NodeWithChildren[] = [
  [
    "Colour and contrast",
    "Palettes, colour systems and accessible contrast.",
    "Open for contributions",
    1200,
    700,
  ],
  [
    "Interface patterns",
    "Useful patterns for navigation, forms and everyday interactions.",
    "Open for contributions",
    700,
    1000,
  ],
  [
    "Design reading list",
    "Books and essays recommended by the community.",
    "Open for contributions",
    0,
    0,
  ],
  [
    "Typography",
    "Good typography makes an interface easier to read and gives it a clear voice.",
    "Growing collection",
    1200,
    900,
  ],
].map(([name, description, status, width, height], index) => ({
  id: `design-${index}`,
  name: String(name),
  slug: `design-${index}`,
  description: String(description),
  content: `<p>${description}</p>${"<p>Explore practical examples, compare approaches, and record what works for the community. These field notes include context, decisions, and details worth keeping alongside the reference.</p>".repeat(5)}`,
  createdAt: "2026-11-01T00:00:00Z",
  updatedAt: "2026-11-01T00:00:00Z",
  visibility: "published",
  hide_child_tree: false,
  meta: {},
  assets: [],
  tags: [],
  owner: {
    id: "curator",
    handle: "curator",
    name: "Curator",
    joined: "2026-01-01T00:00:00Z",
    roles: [],
  },
  properties: schema.map((property) => ({
    ...property,
    value: property.fid === "status" ? String(status) : "Reference",
  })),
  child_property_schema: [],
  children: [],
  recomentations: [],
  primary_image: width
    ? {
        id: `image-${index}`,
        path: image,
        filename: "mountains.png",
        mime_type: "image/png",
        width: Number(width),
        height: Number(height),
      }
    : undefined,
}));

const meta = {
  title: "Compositions/Library/Directory Layouts",
  // The installed adapter's App Router boundary predates Next's render-tree
  // context. This isolated composition only needs navigation spies.
  beforeEach: () => {
    useRouter.mockReturnValue({
      back: fn(),
      forward: fn(),
      refresh: fn(),
      push: fn(),
      replace: fn(),
      prefetch: fn(),
      bfcacheId: "directory-story",
    });
  },
  parameters: {
    nextjs: { appDirectory: false },
    layout: "padded",
    docs: {
      description: {
        component:
          "Production Library directory cards across Rows, Grid, and Gallery layouts. Check card styles with mixed image ratios and a text-only entry, and place Status in the header, details, or footer corners. Edit mode supports inline editing and drag reordering; changes stay in this story. Gallery supports Minimal and Editorial cards with natural image ratios; Image overlay falls back to Editorial when a page has no cover.",
      },
    },
  },
} satisfies Meta;
export default meta;
type Story = StoryObj<typeof meta>;

export const Playground: Story = { render: () => <DirectoryPlayground /> };

function DirectoryPlayground() {
  const [layout, setLayout] = useState<"rows" | "grid" | "gallery">("rows");
  const [style, setStyle] = useState<"minimal" | "row" | "box" | "fill">("box");
  const [placement, setPlacement] = useState<
    "header" | "details" | "footer-start" | "footer-end"
  >("footer-start");
  const [editing, setEditing] = useState(false);
  const [showContent, setShowContent] = useState(false);
  const [nodes, setNodes] = useState(initialNodes);
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 4 } }),
  );
  const config: LibraryPageBlockTypeDirectoryConfig = {
    layout,
    card_style: style,
    show_content: showContent,
    columns: [
      { fid: "fixed:name", hidden: false },
      { fid: "fixed:description", hidden: false },
      { fid: "fixed:link", hidden: true },
      { fid: "fixed:primary_image", hidden: false },
      { fid: "status", hidden: false, placement },
      { fid: "type", hidden: false, placement: "header" },
    ],
  };
  function update(nodeID: string, fid: string, value: string) {
    setNodes((nodes) =>
      nodes.map((node) =>
        node.id !== nodeID
          ? node
          : fid.startsWith("fixed:")
            ? { ...node, [fid.slice(6)]: value }
            : {
                ...node,
                properties: node.properties.map((p) =>
                  p.fid === fid ? { ...p, value } : p,
                ),
              },
      ),
    );
  }

  return (
    <LStack w="full" maxW="5xl" gap="4">
      <HStack gap="2" flexWrap="wrap">
        <Checkbox
          checked={showContent}
          onCheckedChange={(details) =>
            setShowContent(details.checked === true)
          }
        >
          Rich content
        </Checkbox>
        {(["rows", "grid", "gallery"] as const).map((value) => (
          <Button
            key={value}
            size="sm"
            variant={layout === value ? "solid" : "outline"}
            onClick={() => setLayout(value)}
          >
            {value}
          </Button>
        ))}
        <Checkbox
          checked={editing}
          onCheckedChange={(details) => setEditing(details.checked === true)}
        >
          Edit and reorder
        </Checkbox>
      </HStack>
      {layout !== "rows" && (
        <HStack flexWrap="wrap" gap="2">
          {(
            [
              ["minimal", "Minimal"],
              ["row", "Compact"],
              ["box", "Editorial"],
              ["fill", "Image overlay"],
            ] as const
          ).map(([value, label]) => (
            <Button
              key={value}
              size="sm"
              variant={style === value ? "subtle" : "ghost"}
              onClick={() => setStyle(value)}
            >
              {label}
            </Button>
          ))}
        </HStack>
      )}
      <HStack gap="2" flexWrap="wrap">
        <Text variant="metadata">Status placement</Text>
        {(
          [
            ["header", "Header"],
            ["details", "Details"],
            ["footer-start", "Bottom-left"],
            ["footer-end", "Bottom-right"],
          ] as const
        ).map(([value, label]) => (
          <Button
            key={value}
            size="sm"
            variant={placement === value ? "subtle" : "ghost"}
            onClick={() => setPlacement(value)}
          >
            {label}
          </Button>
        ))}
      </HStack>
      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragEnd={({ active, over }) => {
          if (!over || active.id === over.id) {
            return;
          }
          setNodes((nodes) =>
            arrayMove(
              nodes,
              nodes.findIndex((n) => n.id === active.id),
              nodes.findIndex((n) => n.id === over.id),
            ),
          );
        }}
      >
        <SortableContext
          items={nodes.map((n) => n.id)}
          strategy={rectSortingStrategy}
        >
          <styled.ul
            aria-label="Directory preview"
            w="full"
            gap="2"
            display={layout === "gallery" ? "block" : "grid"}
            gridTemplateColumns={
              layout === "rows"
                ? "1fr"
                : "repeat(auto-fill, minmax(min(100%, 240px), 1fr))"
            }
            columnWidth={layout === "gallery" ? "[240px]" : undefined}
            columnGap="2"
          >
            {nodes.map((node) => (
              <PreviewItem
                key={node.id}
                node={node}
                config={config}
                editing={editing}
                update={update}
              />
            ))}
          </styled.ul>
        </SortableContext>
      </DndContext>
    </LStack>
  );
}

function PreviewItem({
  node,
  config,
  editing,
  update,
}: {
  node: NodeWithChildren;
  config: LibraryPageBlockTypeDirectoryConfig;
  editing: boolean;
  update: (nodeID: string, fid: string, value: string) => void;
}) {
  const sortable = useSortable({ id: node.id, disabled: !editing });
  const block = { type: "directory", config } as const;
  const values = Object.fromEntries(
    mergeFieldsAndProperties([...schema], node, block).map((c) => [c.fid, c]),
  );
  const overlay =
    config.layout === "grid" &&
    config.card_style === "fill" &&
    !config.columns.find((column) => column.fid === "fixed:primary_image")
      ?.hidden &&
    Boolean(node.primary_image);
  return (
    <styled.li
      ref={sortable.setNodeRef}
      className={css({ "& > article": { height: "full" } })}
      minW="0"
      position="relative"
      breakInside="avoid"
      mb={config.layout === "gallery" ? "2" : "0"}
      style={{
        transform: CSS.Translate.toString(sortable.transform),
        transition: sortable.transition,
        opacity: sortable.isDragging ? 0.5 : 1,
      }}
    >
      {editing && (
        <Box
          position="absolute"
          top={overlay ? "3" : "2"}
          right={overlay ? "3" : "2"}
          display="flex"
          zIndex="[6]"
        >
          <IconButton
            ref={sortable.setActivatorNodeRef}
            {...sortable.attributes}
            {...sortable.listeners}
            size="sm"
            variant="subtle"
            minWidth={overlay ? "6" : "5"}
            width={overlay ? "6" : "5"}
            height={overlay ? "6" : "5"}
            padding="0"
            aria-label={`Move ${node.name}`}
          >
            <DragHandleIcon />
          </IconButton>
        </Box>
      )}
      <DirectoryRichCard
        node={node}
        config={config}
        layout={config.layout as "rows" | "grid" | "gallery"}
        columns={mergeFieldsAndPropertySchema([...schema], block).filter(
          (c) => !c.fixed,
        )}
        values={values}
        editing={editing}
        onFieldValueChange={update}
      />
    </styled.li>
  );
}
