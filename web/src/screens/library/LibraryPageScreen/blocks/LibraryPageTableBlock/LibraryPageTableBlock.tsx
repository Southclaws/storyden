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
import * as Table from "@/components/ui/table";
import { Box } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";

import { ColumnMenu } from "./ColumnMenu";
import {
  mergeFieldsAndProperties,
  mergeFieldsAndPropertySchema,
} from "./column";

export function LibraryPageTableBlock() {
  const { node, form } = useLibraryPageContext();
  const { sort, handleSort } = useSortIndicator();

  // format the sort property as "name" or "-name" for asc/desc
  const childrenSort =
    sort !== null
      ? sort?.order === "asc"
        ? sort.property
        : `-${sort.property}`
      : undefined;

  const { data, error } = useNodeListChildren(node.slug, {
    children_sort: childrenSort,
  });

  if (!data) {
    return <Unready error={error} />;
  }

  const nodes = data.nodes;

  if (nodes.length === 0) {
    return null;
  }

  if (!node.hide_child_tree) {
    return null;
  }

  const currentMeta = form.watch("meta");
  const currentChildPropertySchema = form.watch(
    "childPropertySchema",
    node.child_property_schema,
  );

  const block = currentMeta.layout?.blocks.find((b) => b.type === "table");

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
                  <ColumnMenu column={property}>
                    <Box
                      p="2"
                      display="inline-flex"
                      w="full"
                      alignItems="center"
                      justifyContent="space-between"
                      gap="1"
                      flexWrap="nowrap"
                      fontWeight="semibold"
                    >
                      {property.name}
                      <IconButton
                        type="button"
                        variant="ghost"
                        size="xs"
                        onClick={handleClickSortAction}
                      >
                        <SortIndicator order={sortState} />
                      </IconButton>
                    </Box>
                  </ColumnMenu>
                </Table.Header>
              );
            })}

            {/* NOTE: Not available until we move edit mode to immediate mutations */}
            {/* <Table.Header p="1">
              <AddPropertyMenu>
                <IconButton size="xs" variant="ghost">
                  <AddIcon />
                </IconButton>
              </AddPropertyMenu>
            </Table.Header> */}
          </Table.Row>
        </Table.Head>
        <Table.Body>
          <SortableContext items={nodes}>
            {nodes.map((child) => {
              const columns = mergeFieldsAndProperties(
                node.child_property_schema,
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
              <CreatePageAction
                variant="ghost"
                size="xs"
                parentSlug={node.slug}
              />
            </Table.Cell>
          </Table.Row>
        </Table.Foot>
      </Table.Root>
    </Box>
  );
}
