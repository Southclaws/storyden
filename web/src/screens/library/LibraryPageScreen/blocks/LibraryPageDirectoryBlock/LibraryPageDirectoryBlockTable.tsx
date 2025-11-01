import { Portal } from "@ark-ui/react";
import { useDndContext } from "@dnd-kit/core";
import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import Link from "next/link";
import { ChangeEvent, useEffect, useRef, useState } from "react";
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
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
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

  const items = nodes.map((child) => `${nodeID}.${child.id}`);

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
          <SortableContext items={items} strategy={verticalListSortingStrategy}>
            {nodes.map((child) => {
              const columns = mergeFieldsAndProperties(
                currentChildPropertySchema,
                child,
                block,
              ).filter(isAlwaysFilteredForTableViews);

              return (
                <Row
                  key={child.id}
                  nodeID={nodeID}
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
  nodeID,
  columns,
  onFieldValueChange,
  editing,
}: {
  child: NodeWithChildren;
  nodeID: Identifier;
  columns: ColumnValue[];
  onFieldValueChange: (
    nodeID: Identifier,
    fid: MappableNodeField,
    value: string,
  ) => void;
  editing: boolean;
}) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: `${nodeID}.${child.id}`,
    data: {
      type: "node",
      node: child,
      parentID: nodeID,
      context: "node-children",
    } satisfies DragItemNode,
  });

  // Manage the menu state manually due to the complexity of the menu trigger
  // also being a drag handle for the row.
  const [isOpen, setOpen] = useState(false);
  function handleMenuToggle() {
    setOpen((prev) => !prev);
  }

  // Manually handle click-away behaviour - the default menu behaviour has been
  // overridden by making it a controlled component in order to allow for the
  // drag handle to be used as a menu open trigger.
  const elementRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!isOpen) return;

    function handleClickAway(event: MouseEvent) {
      if (
        elementRef.current &&
        !elementRef.current.contains(event.target as Node)
      ) {
        setOpen(false);
      }
    }

    document.addEventListener("click", handleClickAway);
    return () => document.removeEventListener("click", handleClickAway);
  }, [isOpen]);

  // Check if we're dragging anything at all, to hide the tooltip.
  const { active } = useDndContext();
  const isDraggingAnything = active !== null;

  const dragStyle = {
    transform: CSS.Translate.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    flexShrink: 0,
    // willChange: "transform",
  };

  const dragHandleStyle = {
    cursor: isDragging ? "grabbing" : "grab",
  };

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

  return (
    <Table.Row
      className={cx("group", visibilityStyles)}
      ref={setNodeRef}
      style={dragStyle}
      key={child.id}
    >
      {columns.map((column, idx) => {
        const isFirst = idx === 0;

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
            position="relative"
          >
            {editing ? (
              <HStack gap="1">
                {isFirst && (
                  <Box
                    id={`node-${child.id}_gutter-drag-handle`}
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
                        placement: "bottom-start",
                      }}
                    >
                      <Tooltip.Trigger asChild>
                        <Box position="relative">
                          <Box style={dragHandleStyle}>
                            <IconButton
                              style={dragHandleStyle}
                              id={`node-${child.id}_gutter-drag-handle-button`}
                              variant={{
                                base: "subtle",
                                md: "ghost",
                              }}
                              size="xs"
                              minWidth="5"
                              width="5"
                              height="5"
                              padding="0"
                              color="fg.muted"
                              onClick={handleMenuToggle}
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
                    <Box
                      position="absolute"
                      top="0"
                      width="6"
                      height="6"
                      pointerEvents="none"
                      // NOTE: ClickAway Does not work currently, need to do some shenanigans with portal
                      ref={elementRef}
                    >
                      <LibraryPageMenu
                        variant="ghost"
                        node={child}
                        parentID={nodeID}
                        open={isOpen}
                      >
                        <Box position="absolute" width="6" height="6" />
                      </LibraryPageMenu>
                    </Box>
                  </Box>
                )}

                <styled.input
                  w="full"
                  defaultValue={column.value}
                  onChange={handleCellChange}
                  _focusVisible={{
                    outline: "none",
                  }}
                />
              </HStack>
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
