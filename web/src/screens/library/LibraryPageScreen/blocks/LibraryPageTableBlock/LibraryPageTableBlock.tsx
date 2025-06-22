import { SortableContext } from "@dnd-kit/sortable";
import Link from "next/link";

import { useNodeListChildren } from "@/api/openapi-client/nodes";
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
import { Box, Center, HStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";

import { AddPropertyMenu } from "./AddPropertyMenu/AddPropertyMenu";
import { ColumnMenu } from "./ColumnMenu";
import { PropertyListMenu } from "./PropertyListMenu/PropertyListMenu";
import {
  mergeFieldsAndProperties,
  mergeFieldsAndPropertySchema,
} from "./column";
import { useTableBlock } from "./useTableBlock";

export function LibraryPageTableBlock() {
  const { nodeID, initialChildren } = useLibraryPageContext();
  const { sort, handleSort } = useSortIndicator();
  const { editing } = useEditState();

  const hideChildTree = useWatch((s) => s.draft.child_property_schema);

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

  const block = useTableBlock();
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

  if (!hideChildTree) {
    return null;
  }

  if (!block) {
    console.warn(
      "attempting to render a LibraryPageTableBlock without a block in the form metadata",
    );
    return null;
  }

  const columns = mergeFieldsAndPropertySchema(
    currentChildPropertySchema,
    block,
  );

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
                    return (
                      <Table.Cell
                        key={column.fid}
                        // NOTE: This doesn't work in edit mode because "group"
                        // class is also used in the page edit level, need to
                        // create a second level of hover grouping or something.
                        // _groupHover={{
                        //   bg: "bg.muted",
                        // }}
                      >
                        {column.href ? (
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
