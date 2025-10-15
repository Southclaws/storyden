import { Portal } from "@ark-ui/react";
import { SortableContext, useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import {
  FloatingPortal,
  autoUpdate,
  flip,
  offset,
  shift,
  useFloating,
  useMergeRefs,
} from "@floating-ui/react";
import Link from "next/link";
import { ChangeEvent, useLayoutEffect, useRef, useState } from "react";
import { match } from "ts-pattern";

import {
  Identifier,
  NodeWithChildren,
  Permission,
  PropertySchemaList,
  Visibility,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { CreatePageAction } from "@/components/library/CreatePage";
import { CreatePageFromURLAction } from "@/components/library/CreatePageFromURL/CreatePageFromURL";
import { LinkRefButton } from "@/components/library/links/LinkCard";
import { SortIndicator } from "@/components/site/SortIndicator";
import { IconButton } from "@/components/ui/icon-button";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import * as Table from "@/components/ui/table";
import * as Tooltip from "@/components/ui/tooltip";
import { DragItemNode } from "@/lib/dragdrop/provider";
import { visibilityColour } from "@/lib/library/visibilityColours";
import { isValidLinkLike } from "@/lib/link/validation";
import { css, cx } from "@/styled-system/css";
import { Box, HStack, styled } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

import { useLibraryPageContext } from "../../Context";
import { useEditState } from "../../useEditState";

import { ColumnMenu } from "./ColumnMenu";
import { useDirectoryBlockContext } from "./Context";
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

export function LibraryPageDirectoryBlockTable({
  nodes,
  block,
  currentChildPropertySchema,
}: Props) {
  const { nodeID, store } = useLibraryPageContext();
  const { editing } = useEditState();
  const { sort, handleSort, handleMutateChildren } = useDirectoryBlockContext();
  const session = useSession();

  const { setChildPropertyValue } = store.getState();

  const canManageLibrary = hasPermission(session, Permission.MANAGE_LIBRARY);

  const columns = mergeFieldsAndPropertySchema(
    currentChildPropertySchema,
    block,
  ).filter(isAlwaysFilteredForTableViews);

  function handleChildFieldValueChange(
    nodeID: Identifier,
    fid: MappableNodeField,
    value: string,
  ) {
    setChildPropertyValue(nodeID, fid, value);
  }

  function handleCreatePageComplete() {
    handleMutateChildren();
  }

  return (
    <Box w="full" overflowX="scroll">
      <Table.Root size="sm" variant="dense" borderStyle="none" minW="0">
        <Table.Head>
          <Table.Row position="relative">
            {columns.map((property) => {
              const isSortingBy = sort?.property === property.name;
              const isSortingAsc = sort?.order === "asc";
              const isSortingDesc = sort?.order === "desc";
              const isSorting = isSortingBy && (isSortingAsc || isSortingDesc);
              const sortState = isSorting
                ? isSortingAsc
                  ? "asc"
                  : "desc"
                : "none";

              function handleClickSortAction() {
                handleSort(property.name);
              }

              return (
                <Table.Header
                  key={property.fid}
                  cursor="pointer"
                  {...(isSorting && {
                    "data-active": "",
                  })}
                  _hover={{
                    bg: "bg.muted",
                  }}
                  _active={{
                    bg: "bg.muted",
                  }}
                  p="0"
                >
                  <HStack minW="0" w="full" pr="1">
                    <ColumnMenu column={property}>
                      <Box
                        p="2"
                        display="inline-flex"
                        minW="0"
                        flexGrow="1"
                        alignItems="center"
                        justifyContent="space-between"
                        gap="1"
                        textWrap="nowrap"
                        flexWrap="nowrap"
                        fontWeight="semibold"
                      >
                        {property.name}
                      </Box>
                    </ColumnMenu>
                    <IconButton
                      type="button"
                      variant="ghost"
                      size="xs"
                      onClick={handleClickSortAction}
                    >
                      <SortIndicator order={sortState} />
                    </IconButton>
                  </HStack>
                </Table.Header>
              );
            })}
          </Table.Row>
        </Table.Head>
        <Table.Body>
          <SortableContext items={nodes.map((n) => n.id)}>
            {nodes.map((child) => {
              const columns = mergeFieldsAndProperties(
                currentChildPropertySchema,
                child,
                block,
              ).filter(isAlwaysFilteredForTableViews);

              return (
                <Row
                  key={child.id}
                  child={child}
                  columns={columns}
                  onFieldValueChange={handleChildFieldValueChange}
                  editing={editing}
                />
              );
            })}
          </SortableContext>
        </Table.Body>
        {canManageLibrary && (
          <Table.Foot
            borderBottomStyle="solid"
            borderBottomWidth="thin"
            borderBlockColor="border.subtle"
          >
            <Table.Row>
              <Table.Cell colSpan={columns.length}>
                <HStack gap="2">
                  <CreatePageAction
                    variant="ghost"
                    size="xs"
                    parentSlug={nodeID}
                    disableRedirect
                    onComplete={handleCreatePageComplete}
                  />
                  <CreatePageFromURLAction
                    variant="ghost"
                    size="xs"
                    parentSlug={nodeID}
                    onComplete={handleCreatePageComplete}
                  />
                </HStack>
              </Table.Cell>
            </Table.Row>
          </Table.Foot>
        )}
      </Table.Root>
    </Box>
  );
}

// NOTE: Primary image field is always filtered out for table views.
// TODO: Design a nice way to display this without destroying the layout?
function isAlwaysFilteredForTableViews(c: { fid: string }) {
  return c.fid !== "fixed:primary_image";
}

function checkValidColumnValue(column: ColumnValue) {
  if (column.fid === "fixed:link") {
    if (!column.value) {
      return true;
    }

    return isValidLinkLike(column.value);
  }

  return true;
}

function Row({
  child,
  columns,
  onFieldValueChange,
  editing,
}: {
  child: NodeWithChildren;
  columns: ColumnValue[];
  onFieldValueChange: (
    nodeID: Identifier,
    fid: MappableNodeField,
    value: string,
  ) => void;
  editing: boolean;
}) {
  const { refs, floatingStyles } = useFloating({
    nodeId: child.id,
    whileElementsMounted: autoUpdate,
    placement: "left",
    middleware: [offset(7), flip(), shift()],
  });

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: child.id,
    data: {
      type: "node",
      node: child,
    } as DragItemNode,
  });

  const visCol = visibilityColour(child.visibility);

  const visibilityStyles = css({
    colorPalette: visCol,
    boxSizing: "content-box",
    borderLeftWidth:
      child.visibility === Visibility.published ? "none" : "medium",
    borderLeftColor:
      child.visibility === Visibility.published
        ? "transparent"
        : "colorPalette.border",
    borderLeftStyle:
      child.visibility === Visibility.published ? "solid" : "dashed",
  });

  const dragStyle = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    flexShrink: 0,
  };

  const dragHandleStyle = {
    cursor: isDragging ? "grabbing" : "grab",
  };

  const rowRef = useMergeRefs([setNodeRef, refs.setReference]);

  return (
    <Table.Row
      ref={rowRef}
      className={cx("group", visibilityStyles)}
      key={child.id}
      style={dragStyle}
    >
      {columns.map((column, idx) => {
        function handleCellChange(v: ChangeEvent<HTMLInputElement>) {
          onFieldValueChange(
            child.id,
            column.fid as MappableNodeField,
            v.target.value,
          );
        }

        // NOTE: Does not work because this is uncontrolled at the
        // moment. Making it controlled would be a bit of a pain.
        // But it's here and ready in case we ever actually do that.
        const isValid = checkValidColumnValue(column);

        // TODO: Pass in as prop.
        const isFirstColumn = idx === 0;

        return (
          <Table.Cell
            key={column.fid}
            // NOTE: This doesn't work in edit mode because "group"
            // class is also used in the page edit level, need to
            // create a second level of hover grouping or something.
            // _groupHover={{
            //   bg: "bg.muted",
            // }}
            _hover={{ bg: "bg.subtle" }}
            // position="relative"
          >
            {editing ? (
              <>
                {isFirstColumn && (
                  <FloatingPortal>
                    <Box
                      ref={refs.setFloating}
                      style={floatingStyles}
                      {...listeners}
                      {...attributes}
                    >
                      <Tooltip.Root
                        openDelay={0}
                        closeDelay={0}
                        // disabled={isDraggingAnything}
                        positioning={{
                          slide: true,
                          gutter: 4,
                          placement: "bottom-start",
                        }}
                      >
                        <Tooltip.Trigger asChild>
                          <Box position="relative">
                            <Box style={dragHandleStyle}>
                              <IconButton
                                style={dragHandleStyle}
                                variant="ghost"
                                size="xs"
                                minWidth="5"
                                width="5"
                                height="5"
                                padding="0"
                                color="fg.muted"
                                // onClick={handleMenuToggle}
                              >
                                <DragHandleIcon width="4" />
                              </IconButton>
                            </Box>
                          </Box>
                        </Tooltip.Trigger>
                        <Portal>
                          <Tooltip.Positioner>
                            <Tooltip.Arrow>
                              <Tooltip.ArrowTip />
                            </Tooltip.Arrow>

                            <Tooltip.Content p="1" borderRadius="sm">
                              <p>
                                <styled.span fontWeight="semibold">
                                  Click
                                </styled.span>
                                &nbsp;
                                <styled.span fontWeight="normal">
                                  to open menu
                                </styled.span>
                              </p>
                              <p>
                                <styled.span fontWeight="semibold">
                                  Drag
                                </styled.span>
                                &nbsp;
                                <styled.span fontWeight="normal">
                                  to move
                                </styled.span>
                              </p>
                            </Tooltip.Content>
                          </Tooltip.Positioner>
                        </Portal>
                      </Tooltip.Root>
                    </Box>
                  </FloatingPortal>
                )}
                <styled.input
                  w="full"
                  defaultValue={column.value}
                  onChange={handleCellChange}
                  _focusVisible={{
                    outline: "none",
                  }}
                />
              </>
            ) : (
              <Box minH="4">
                {match(column.fid)
                  .with("fixed:name", () => (
                    <Link href={column.href ?? "#"}>
                      {column.value || (
                        <styled.em color="fg.muted">(untitled page)</styled.em>
                      )}
                    </Link>
                  ))
                  .with("fixed:link", () =>
                    child.link ? (
                      <LinkRefButton link={child.link} variant="link" />
                    ) : (
                      <>{column.value}</>
                    ),
                  )
                  .otherwise(() => (
                    <>{column.value}</>
                  ))}
              </Box>
            )}
          </Table.Cell>
        );
      })}
    </Table.Row>
  );
}
