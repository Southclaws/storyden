import { useDndContext } from "@dnd-kit/core";
import {
  SortableContext,
  rectSortingStrategy,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { keyBy } from "lodash";
import { CSSProperties } from "react";

import {
  Identifier,
  NodeWithChildren,
  PropertySchemaList,
  Visibility,
} from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";
import { DragItemNode } from "@/lib/dragdrop/provider";
import { visibilityColour } from "@/lib/library/visibilityColours";
import { css, cx, sva } from "@/styled-system/css";
import { Center, GridItem, styled } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";

import { DirectoryCardControls } from "./DirectoryCardControls";
import { DirectoryPropertyLabel } from "./DirectoryCardFields";
import { DirectoryRichCard } from "./DirectoryRichCard";
import {
  ColumnDefinition,
  mergeFieldsAndProperties,
  mergeFieldsAndPropertySchema,
} from "./column";
import {
  DirectoryCardProps,
  applyDirectoryDraft,
  getRichCardShape,
} from "./directoryCard";
import { useDirectoryBlock } from "./useDirectoryBlock";

type Props = {
  nodes: NodeWithChildren[];
  block: ReturnType<typeof useDirectoryBlock>;
  currentChildPropertySchema: PropertySchemaList;
  layout: "rows" | "grid" | "gallery";
};

export function LibraryPageDirectoryBlockGrid({
  nodes,
  block,
  currentChildPropertySchema,
  layout,
}: Props) {
  const { nodeID, store } = useLibraryPageContext();
  const { isDirectEditing } = useEditState();
  const draftChildren = useWatch((state) => state.draft.children);
  const draftChildrenByID = keyBy(draftChildren, "id");

  const { setChildPropertyValue } = store.getState();

  const fullSchema = mergeFieldsAndPropertySchema(
    currentChildPropertySchema,
    block,
  );

  const columns = fullSchema.filter(
    // We don't show fixed fields in the grid view. They are laid
    // out by the component itself and can be toggled on or off.
    (c) => !c._fixedFieldName,
  );
  const config = block.config ?? { layout, columns: fullSchema };

  function handleChildFieldValueChange(
    nodeID: Identifier,
    fid: string,
    value: string,
  ) {
    setChildPropertyValue(nodeID, fid, value);
  }

  if (nodes.length === 0) {
    return (
      <Center w="full">
        <EmptyState hideContributionLabel>There are no pages here.</EmptyState>
      </Center>
    );
  }

  let sortingStrategy = rectSortingStrategy;

  if (layout === "rows") {
    sortingStrategy = verticalListSortingStrategy;
  }

  const items = nodes.map((child) => `${nodeID}.${child.id}`);

  return (
    <SortableContext items={items} strategy={sortingStrategy}>
      <styled.ul
        aria-label="Library pages"
        data-layout={layout}
        className={directoryGrid({ layout }).list}
      >
        {nodes.map((sourceNode) => {
          let draftNode: NodeWithChildren | undefined;

          if (isDirectEditing) {
            draftNode = draftChildrenByID[sourceNode.id];
          }

          const node = applyDirectoryDraft(sourceNode, draftNode);
          const columnValues = mergeFieldsAndProperties(
            currentChildPropertySchema,
            node,
            block,
          );

          const valueTable = keyBy(columnValues, "fid");

          return (
            <GridCard
              key={node.id}
              layout={layout}
              config={config}
              nodeID={nodeID}
              node={node}
              editing={isDirectEditing}
              values={valueTable}
              columns={columns}
              onFieldValueChange={handleChildFieldValueChange}
            />
          );
        })}
      </styled.ul>
    </SortableContext>
  );
}

type GridCardProps = DirectoryCardProps & {
  config: NonNullable<ReturnType<typeof useDirectoryBlock>["config"]>;
  nodeID: Identifier;
};

function GridCard(props: GridCardProps) {
  const {
    config,
    layout,
    nodeID,
    node,
    editing,
    values,
    columns,
    onFieldValueChange,
  } = props;
  const {
    attributes,
    listeners,
    setNodeRef,
    setActivatorNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: `${nodeID}.${node.id}`,
    data: {
      type: "node",
      node: node,
      parentID: nodeID,
      context: "node-children",
    } satisfies DragItemNode,
  });

  const { active } = useDndContext();
  function renderPropertyLabel(column: ColumnDefinition) {
    return <DirectoryPropertyLabel column={column} editing={editing} />;
  }

  const cardProps = {
    node,
    layout,
    columns,
    values,
    editing,
    onFieldValueChange,
    renderPropertyLabel,
  };
  const showCover = !config.columns.find(
    (column) => column.fid === "fixed:primary_image",
  )?.hidden;
  const hasImage = showCover && Boolean(node.primary_image);
  const shape = getRichCardShape(layout, config.card_style, hasImage);
  const overlay = shape === "fill";

  const dragStyle: CSSProperties = {
    transform: CSS.Translate.toString(transform),
    transition,
    opacity: 1,
  };

  if (isDragging) {
    dragStyle.opacity = 0.5;
    dragStyle.zIndex = 10;
  }

  const visibilityStyles = css({
    colorPalette: visibilityColour(node.visibility),
    "& > article": { height: "full" },
  });
  let borderStyles = css({
    borderWidth: "none",
    borderColor: "transparent",
    borderStyle: "none",
  });

  if (node.visibility !== Visibility.published) {
    borderStyles = css({
      borderWidth: "thin",
      borderColor: "colorPalette.border",
      borderStyle: "dashed",
    });
  }

  return (
    <GridItem
      as="li"
      ref={setNodeRef}
      style={dragStyle}
      className={cx(visibilityStyles, borderStyles)}
      position="relative"
      borderRadius="lg"
      data-shape={shape}
      css={{ "&[data-shape=minimal]": { borderRadius: "md" } }}
      display="flex"
      width="full"
      flexDirection="column"
      justifyContent="space-between"
      overflow="hidden"
      gap="2"
    >
      {editing && (
        <DirectoryCardControls
          node={node}
          parentID={nodeID}
          overlay={overlay}
          isDragging={isDragging}
          isDraggingAnything={active !== null}
          attributes={attributes}
          listeners={listeners}
          setActivatorNodeRef={setActivatorNodeRef}
        />
      )}
      <DirectoryRichCard {...cardProps} config={config} />
    </GridItem>
  );
}

const directoryGrid = sva({
  slots: ["list"],
  base: {
    list: {
      width: "full",
      gap: "2",
      "& > li": { minWidth: "0" },
    },
  },
  variants: {
    layout: {
      rows: {
        list: { display: "flex", flexDirection: "column" },
      },
      grid: {
        list: {
          display: "grid",
          gridTemplateColumns:
            "repeat(auto-fill, minmax(min(100%, 200px), 1fr))",
        },
      },
      gallery: {
        list: {
          display: "block",
          columnWidth: "[200px]",
          columnGap: "2",
          columnFill: "balance",
          "& > li": { breakInside: "avoid", marginBottom: "2" },
        },
      },
    },
  },
});
