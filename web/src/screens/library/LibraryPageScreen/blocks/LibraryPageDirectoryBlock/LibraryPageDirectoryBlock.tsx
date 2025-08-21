import { SortableContext } from "@dnd-kit/sortable";
import { keyBy } from "lodash";
import Link from "next/link";
import { ChangeEvent } from "react";

import { useNodeListChildren } from "@/api/openapi-client/nodes";
import {
  Identifier,
  NodeWithChildren,
  PropertySchemaList,
} from "@/api/openapi-schema";
import { CreatePageAction } from "@/components/library/CreatePage";
import {
  SortIndicator,
  useSortIndicator,
} from "@/components/site/SortIndicator";
import { Unready } from "@/components/site/Unready";
import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { MenuIcon } from "@/components/ui/icons/Menu";
import * as Table from "@/components/ui/table";
import { isValidLinkLike } from "@/lib/link/validation";
import {
  Box,
  Center,
  Grid,
  GridItem,
  HStack,
  LStack,
  VStack,
  WStack,
  styled,
} from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";

import { AddPropertyMenu } from "./AddPropertyMenu/AddPropertyMenu";
import { ColumnMenu } from "./ColumnMenu";
import { PropertyListMenu } from "./PropertyListMenu/PropertyListMenu";
import {
  ColumnValue,
  MappableNodeField,
  mergeFieldsAndProperties,
  mergeFieldsAndPropertySchema,
} from "./column";
import { useDirectoryBlock } from "./useDirectoryBlock";

type LibraryPageDirectoryBlockTableProps = {
  nodes: NodeWithChildren[];
  block: ReturnType<typeof useDirectoryBlock>;
  currentChildPropertySchema: PropertySchemaList;
};

export function LibraryPageDirectoryBlock() {
  const { nodeID, initialChildren, store } = useLibraryPageContext();
  const { sort, handleSort } = useSortIndicator();
  const { editing } = useEditState();

  const { setChildPropertyValue } = store.getState();

  // format the sort property as "name" or "-name" for asc/desc
  const childrenSort =
    sort !== null
      ? sort?.order === "asc"
        ? sort.property
        : `-${sort.property}`
      : undefined;

  const { data, error } = useNodeListChildren(
    nodeID,
    {
      children_sort: childrenSort,
    },
    {
      swr: {
        fallbackData: initialChildren,
      },
    },
  );

  const block = useDirectoryBlock();
  const currentChildPropertySchema = useWatch(
    (s) => s.draft.child_property_schema,
  );

  if (!data) {
    return <Unready error={error} />;
  }

  const nodes = data.nodes;

  if (nodes.length === 0) {
    return null;
  }

  if (!block) {
    console.warn(
      "attempting to render a LibraryPageDirectoryBlock without a block in the form metadata",
    );
    return null;
  }

  // Get layout from config, default to table
  const layout = block.config?.layout ?? "table";

  // Switch between different layout views
  switch (layout) {
    case "grid":
      return (
        <LibraryPageDirectoryBlockGrid
          nodes={nodes}
          block={block}
          currentChildPropertySchema={currentChildPropertySchema}
        />
      );
    case "table":
    default:
      return (
        <LibraryPageDirectoryBlockTable
          nodes={nodes}
          block={block}
          currentChildPropertySchema={currentChildPropertySchema}
        />
      );
  }
}

function LibraryPageDirectoryBlockTable({
  nodes,
  block,
  currentChildPropertySchema,
}: LibraryPageDirectoryBlockTableProps) {
  const { nodeID, store } = useLibraryPageContext();
  const { sort, handleSort } = useSortIndicator();
  const { editing } = useEditState();

  const { setChildPropertyValue } = store.getState();

  const columns = mergeFieldsAndPropertySchema(
    currentChildPropertySchema,
    block,
  );

  function handleChildFieldValueChange(
    nodeID: Identifier,
    fid: MappableNodeField,
    value: string,
  ) {
    console.log("Child field value changed:", nodeID, fid, value);
    setChildPropertyValue(nodeID, fid, value);
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

            {editing && (
              <Table.Header width="0">
                <Center>
                  <AddPropertyMenu>
                    <IconButton size="xs" variant="ghost">
                      <AddIcon />
                    </IconButton>
                  </AddPropertyMenu>
                  <PropertyListMenu>
                    <IconButton size="xs" variant="ghost">
                      <MenuIcon />
                    </IconButton>
                  </PropertyListMenu>
                </Center>
              </Table.Header>
            )}
          </Table.Row>
        </Table.Head>
        <Table.Body>
          <SortableContext items={nodes}>
            {nodes.map((child) => {
              const columns = mergeFieldsAndProperties(
                currentChildPropertySchema,
                child,
                block,
              );

              return (
                <Table.Row key={child.id} className="group">
                  {columns.map((column) => {
                    function handleCellChange(
                      v: ChangeEvent<HTMLInputElement>,
                    ) {
                      handleChildFieldValueChange(
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
                      >
                        {editing ? (
                          <styled.input
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
                      </Table.Cell>
                    );
                  })}
                </Table.Row>
              );
            })}
          </SortableContext>
        </Table.Body>
        <Table.Foot
          borderBottomStyle="solid"
          borderBottomWidth="thin"
          borderBlockColor="border.subtle"
        >
          <Table.Row>
            <Table.Cell colSpan={columns.length} display="flex">
              <CreatePageAction variant="ghost" size="xs" parentSlug={nodeID} />
            </Table.Cell>
          </Table.Row>
        </Table.Foot>
      </Table.Root>
    </Box>
  );
}

type LibraryPageDirectoryBlockGridProps = {
  nodes: NodeWithChildren[];
  block: ReturnType<typeof useDirectoryBlock>;
  currentChildPropertySchema: PropertySchemaList;
};

function LibraryPageDirectoryBlockGrid({
  nodes,
  block,
  currentChildPropertySchema,
}: LibraryPageDirectoryBlockGridProps) {
  const { store } = useLibraryPageContext();
  const { sort, handleSort } = useSortIndicator();
  const { editing } = useEditState();

  const { setChildPropertyValue } = store.getState();

  const columns = mergeFieldsAndPropertySchema(
    currentChildPropertySchema,
    block,
  );

  function handleChildFieldValueChange(
    nodeID: Identifier,
    fid: MappableNodeField,
    value: string,
  ) {
    console.log("Child field value changed:", nodeID, fid, value);
    setChildPropertyValue(nodeID, fid, value);
  }

  return (
    <LStack w="full" gap="2">
      {editing && (
        <WStack
          justifyContent="end"
          bgColor="bg.subtle"
          borderRadius="sm"
          p="1"
        >
          <AddPropertyMenu>
            <IconButton size="xs" variant="ghost">
              <AddIcon />
            </IconButton>
          </AddPropertyMenu>

          <PropertyListMenu hideFixedFields>
            <IconButton size="xs" variant="ghost">
              <MenuIcon />
            </IconButton>
          </PropertyListMenu>
        </WStack>
      )}
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

          return (
            <GridItem
              key={node.id}
              p="2"
              borderRadius="md"
              bg="bg.subtle"
              display="flex"
              flexDirection="column"
              justifyContent="space-between"
              gap="2"
            >
              <LStack gap="0">
                {editing ? (
                  <styled.input
                    w="full"
                    fontWeight="semibold"
                    defaultValue={node.name}
                    onChange={(event) =>
                      handleChildFieldValueChange(
                        node.id,
                        "name",
                        event.target.value,
                      )
                    }
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
                )}
                {editing ? (
                  <styled.input
                    w="full"
                    placeholder="Description..."
                    _placeholder={{
                      color: "fg.subtle",
                    }}
                    defaultValue={node.description}
                    onChange={(event) =>
                      handleChildFieldValueChange(
                        node.id,
                        "description",
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
                )}
              </LStack>

              <styled.dl className={lstack()} gap="0">
                {columns.map((property) => {
                  if (property._fixedFieldName) {
                    // We don't show fixed fields in the grid view. They are laid
                    // out by the component itself and can be toggled on or off.
                    return null;
                  }

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
                  const isSorting =
                    isSortingBy && (isSortingAsc || isSortingDesc);
                  const sortState = isSorting
                    ? isSortingAsc
                      ? "asc"
                      : "desc"
                    : "none";

                  function handleClickSortAction() {
                    handleSort(property.name);
                  }

                  const handleCellChange = (
                    v: ChangeEvent<HTMLInputElement>,
                  ) => {
                    handleChildFieldValueChange(
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
            </GridItem>
          );
        })}
      </Grid>
    </LStack>
  );
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
