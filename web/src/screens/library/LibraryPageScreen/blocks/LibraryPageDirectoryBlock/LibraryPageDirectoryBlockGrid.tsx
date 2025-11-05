import { Portal } from "@ark-ui/react";
import { useDndContext } from "@dnd-kit/core";
import {
  SortableContext,
  rectSortingStrategy,
  useSortable,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { keyBy } from "lodash";
import Link from "next/link";
import { ChangeEvent, useState } from "react";

import {
  Identifier,
  NodeWithChildren,
  PropertySchemaList,
  Visibility,
} from "@/api/openapi-schema";
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
import { EmptyState } from "@/components/site/EmptyState";
import { useSortIndicator } from "@/components/site/SortIndicator";
import { IconButton } from "@/components/ui/icon-button";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import { EmptyIcon } from "@/components/ui/icons/Empty";
import * as Tooltip from "@/components/ui/tooltip";
import { DragItemNode } from "@/lib/dragdrop/provider";
import { visibilityColour } from "@/lib/library/visibilityColours";
import { css, cx } from "@/styled-system/css";
import {
  Box,
  Center,
  Grid,
  GridItem,
  HStack,
  LStack,
  styled,
} from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { linkDisabledProps } from "@/utils/anchor";
import { getAssetURL } from "@/utils/asset";

import { useLibraryPageContext } from "../../Context";
import { useEditState } from "../../useEditState";

import { ColumnMenu } from "./ColumnMenu";
import {
  ColumnValue,
  MappableNodeField,
  mergeFieldsAndProperties,
  mergeFieldsAndPropertySchema,
} from "./column";
import { useDirectoryBlock } from "./useDirectoryBlock";

type Props = {
  nodes: NodeWithChildren[];
  block: ReturnType<typeof useDirectoryBlock>;
  currentChildPropertySchema: PropertySchemaList;
};

export function LibraryPageDirectoryBlockGrid({
  nodes,
  block,
  currentChildPropertySchema,
}: Props) {
  const { nodeID, store } = useLibraryPageContext();
  const { sort, handleSort } = useSortIndicator();
  const { editing } = useEditState();

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

  const nameColumn = fullSchema.find((c) => c.fid === "fixed:name");
  const nameColumnHiddenState = nameColumn?.hidden ?? true;
  const linkColumn = fullSchema.find((c) => c.fid === "fixed:link");
  const linkColumnHiddenState = linkColumn?.hidden ?? true;
  const descColumn = fullSchema.find((c) => c.fid === "fixed:description");
  const descColumnHiddenState = descColumn?.hidden ?? true;

  const coverImageHiddenState =
    block.config?.columns.find((c) => c.fid === "fixed:primary_image")
      ?.hidden ?? false;

  const fullBleedCover =
    coverImageHiddenState === false &&
    nameColumnHiddenState === true &&
    linkColumnHiddenState === true &&
    descColumnHiddenState === true;

  function handleChildFieldValueChange(
    nodeID: Identifier,
    fid: MappableNodeField,
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

  const items = nodes.map((child) => `${nodeID}.${child.id}`);

  return (
    <SortableContext items={items} strategy={rectSortingStrategy}>
      <Grid
        w="full"
        gap="2"
        gridTemplateColumns="repeat(auto-fill, minmax(200px, 1fr))"
      >
        {nodes.map((node) => {
          const columnValues = mergeFieldsAndProperties(
            currentChildPropertySchema,
            node,
            block,
          );

          const valueTable = keyBy(columnValues, "fid");

          // When the cover image is displayed and the node has no primary image
          // and the name column is hidden, show the node name as a placeholder.
          const coverImagePlaceholder =
            node.primary_image === undefined &&
            !coverImageHiddenState &&
            nameColumnHiddenState;

          return (
            <GridCard
              key={node.id}
              nodeID={nodeID}
              node={node}
              editing={editing}
              columnValues={columnValues}
              valueTable={valueTable}
              coverImagePlaceholder={coverImagePlaceholder}
              fullBleedCover={fullBleedCover}
              coverImageHiddenState={coverImageHiddenState}
              nameColumnHiddenState={nameColumnHiddenState}
              descColumnHiddenState={descColumnHiddenState}
              linkColumnHiddenState={linkColumnHiddenState}
              columns={columns}
              sort={sort}
              handleSort={handleSort}
              onFieldValueChange={handleChildFieldValueChange}
            />
          );
        })}
      </Grid>
    </SortableContext>
  );
}

type GridCardProps = {
  nodeID: Identifier;
  node: NodeWithChildren;
  editing: boolean;
  columnValues: ColumnValue[];
  valueTable: Record<string, ColumnValue>;
  coverImagePlaceholder: boolean;
  fullBleedCover: boolean;
  coverImageHiddenState: boolean;
  nameColumnHiddenState: boolean;
  descColumnHiddenState: boolean;
  linkColumnHiddenState: boolean;
  columns: ReturnType<typeof mergeFieldsAndPropertySchema>;
  sort: ReturnType<typeof useSortIndicator>["sort"];
  handleSort: ReturnType<typeof useSortIndicator>["handleSort"];
  onFieldValueChange: (
    nodeID: Identifier,
    fid: MappableNodeField,
    value: string,
  ) => void;
};

function GridCard({
  nodeID,
  node,
  editing,
  columnValues,
  valueTable,
  coverImagePlaceholder,
  fullBleedCover,
  coverImageHiddenState,
  nameColumnHiddenState,
  descColumnHiddenState,
  linkColumnHiddenState,
  columns,
  sort,
  handleSort,
  onFieldValueChange,
}: GridCardProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
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

  const [isOpen, setOpen] = useState(false);
  function handleMenuToggle(e: React.MouseEvent) {
    e.stopPropagation();
    setOpen((prev) => !prev);
  }

  const { active } = useDndContext();
  const isDraggingAnything = active !== null;

  const dragStyle = {
    transform: CSS.Translate.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    zIndex: isDragging ? 10 : undefined,
  };

  const dragHandleStyle = {
    cursor: isDragging ? "grabbing" : "grab",
  };

  const visCol = visibilityColour(node.visibility);

  const visibilityStyles = css({
    colorPalette: visCol,
    borderWidth: node.visibility === Visibility.published ? "none" : "thin",
    borderColor:
      node.visibility === Visibility.published
        ? "transparent"
        : "colorPalette.border",
    borderStyle: node.visibility === Visibility.published ? "none" : "dashed",
  });

  return (
    <GridItem
      ref={setNodeRef}
      style={dragStyle}
      className={cx(visibilityStyles)}
      position="relative"
      borderRadius="md"
      bg="bg.muted"
      display="flex"
      flexDirection="column"
      justifyContent="space-between"
      overflow="hidden"
      gap="2"
      pb={fullBleedCover ? "0" : "1"}
    >
      {editing && (
        <Box
          position="absolute"
          top="2"
          right="2"
          {...listeners}
          {...attributes}
        >
          <Tooltip.Root
            openDelay={0}
            closeDelay={0}
            disabled={isDraggingAnything}
            positioning={{
              slide: true,
              gutter: 4,
              placement: "top",
            }}
          >
            <Tooltip.Trigger asChild>
              <Box position="relative">
                <IconButton
                  style={dragHandleStyle}
                  variant="subtle"
                  size="xs"
                  minWidth="5"
                  width="5"
                  height="5"
                  padding="0"
                  onClick={handleMenuToggle}
                >
                  <DragHandleIcon width="4" />
                </IconButton>
              </Box>
            </Tooltip.Trigger>
            <Portal>
              <Tooltip.Positioner>
                <Tooltip.Arrow>
                  <Tooltip.ArrowTip />
                </Tooltip.Arrow>

                <Tooltip.Content p="1" borderRadius="sm">
                  <p>
                    <styled.span fontWeight="semibold">Click</styled.span>
                    &nbsp;
                    <styled.span fontWeight="normal">to open menu</styled.span>
                  </p>
                  <p>
                    <styled.span fontWeight="semibold">Drag</styled.span>
                    &nbsp;
                    <styled.span fontWeight="normal">to move</styled.span>
                  </p>
                </Tooltip.Content>
              </Tooltip.Positioner>
            </Portal>
          </Tooltip.Root>
          <Box
            position="absolute"
            top="0"
            width="6"
            height="6"
            pointerEvents="none"
          >
            <LibraryPageMenu
              variant="ghost"
              node={node}
              parentID={nodeID}
              open={isOpen}
              onOpenChange={(details) => setOpen(details.open)}
              onInteractOutside={() => setOpen(false)}
            >
              <Box position="absolute" width="6" height="6" />
            </LibraryPageMenu>
          </Box>
        </Box>
      )}

      {!coverImageHiddenState ? (
        <Link {...linkDisabledProps(editing)} href={`/l/${node.slug}`}>
          {node.primary_image ? (
            <styled.img
              src={getAssetURL(node.primary_image.path)}
              objectFit="cover"
              w="full"
              h="32"
            />
          ) : (
            <Center h="32" p="2">
              {coverImagePlaceholder ? (
                <styled.span
                  textAlign="center"
                  textWrap="balance"
                  color="fg.muted"
                  lineClamp={3}
                >
                  {node.name}
                </styled.span>
              ) : (
                <EmptyIcon />
              )}
            </Center>
          )}
        </Link>
      ) : (
        <Box />
      )}

      {fullBleedCover ? null : (
        <LStack gap="0" px="2">
          {!nameColumnHiddenState &&
            (editing ? (
              <styled.input
                w="full"
                fontWeight="semibold"
                defaultValue={node.name}
                onChange={(event) => {
                  onFieldValueChange(
                    node.id,
                    "fixed:name" as MappableNodeField,
                    event.target.value,
                  );
                }}
                _focusVisible={{
                  outline: "none",
                }}
              />
            ) : (
              <Link href={`/l/${node.slug}`}>
                <styled.div
                  fontWeight="semibold"
                  lineClamp="1"
                  textOverflow="ellipsis"
                  wordBreak="break-all"
                >
                  {node.name}
                </styled.div>
              </Link>
            ))}

          {!descColumnHiddenState &&
            (editing ? (
              <styled.input
                w="full"
                placeholder="Description..."
                _placeholder={{
                  color: "fg.subtle",
                }}
                defaultValue={node.description}
                onChange={(event) =>
                  onFieldValueChange(
                    node.id,
                    "fixed:description" as MappableNodeField,
                    event.target.value,
                  )
                }
                _focusVisible={{
                  outline: "none",
                }}
              />
            ) : (
              node.description && (
                <styled.div
                  fontSize="sm"
                  color="fg.muted"
                  lineClamp="1"
                  textOverflow="ellipsis"
                  wordBreak="break-all"
                >
                  {node.description}
                </styled.div>
              )
            ))}

          {!linkColumnHiddenState &&
            (editing ? (
              <styled.input
                w="full"
                placeholder="Link..."
                _placeholder={{
                  color: "fg.subtle",
                }}
                defaultValue={node.link?.url}
                onChange={(event) =>
                  onFieldValueChange(
                    node.id,
                    "fixed:link" as MappableNodeField,
                    event.target.value,
                  )
                }
                _focusVisible={{
                  outline: "none",
                }}
              />
            ) : (
              node.link && (
                <styled.a
                  fontSize="sm"
                  color="fg.muted"
                  lineClamp="1"
                  textOverflow="ellipsis"
                  wordBreak="break-all"
                  href={node.link.url}
                >
                  {node.link.title || node.link.url}
                </styled.a>
              )
            ))}
        </LStack>
      )}

      {columns.length > 0 && (
        <styled.dl className={lstack()} gap="0" px="2" pb="2">
          {columns.map((property) => {
            const column = valueTable[property.fid];
            if (!column) {
              console.warn(
                `unable to find property ${property.fid} in value table`,
                valueTable,
              );
              return null;
            }

            const isSortingBy = sort?.property === property.name;
            const isSortingAsc = sort?.order === "asc";
            const isSortingDesc = sort?.order === "desc";
            const isSorting = isSortingBy && (isSortingAsc || isSortingDesc);

            const handleCellChange = (v: ChangeEvent<HTMLInputElement>) => {
              onFieldValueChange(
                node.id,
                column.fid as MappableNodeField,
                v.target.value,
              );
            };

            return (
              <HStack
                w="full"
                key={property.fid}
                cursor="pointer"
                {...(isSorting && {
                  "data-active": "",
                })}
                p="0"
              >
                <HStack w="full" pr="1">
                  <ColumnMenu column={property}>
                    <styled.dt
                      display="inline-flex"
                      minW="0"
                      flexGrow="1"
                      alignItems="center"
                      justifyContent="space-between"
                      gap="1"
                      textWrap="nowrap"
                      flexWrap="nowrap"
                      color="fg.subtle"
                    >
                      {property.name}
                    </styled.dt>
                  </ColumnMenu>

                  <styled.dd>
                    {editing ? (
                      <styled.input
                        w="full"
                        minW="0"
                        textAlign="right"
                        defaultValue={column.value}
                        onChange={handleCellChange}
                        _focusVisible={{
                          outline: "none",
                        }}
                      />
                    ) : column.href ? (
                      <Link href={column.href}>{column.value}</Link>
                    ) : (
                      <>{column.value}</>
                    )}
                  </styled.dd>
                </HStack>
              </HStack>
            );
          })}
        </styled.dl>
      )}
    </GridItem>
  );
}
