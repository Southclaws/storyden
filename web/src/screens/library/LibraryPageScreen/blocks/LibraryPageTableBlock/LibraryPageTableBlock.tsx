import { SortableContext } from "@dnd-kit/sortable";

import { useNodeListChildren } from "@/api/openapi-client/nodes";
import { PropertySchema } from "@/api/openapi-schema";
import {
  SortIndicator,
  useSortIndicator,
} from "@/components/site/SortIndicator";
import { Unready } from "@/components/site/Unready";
import * as Table from "@/components/ui/table";
import { Box } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";

export function LibraryPageTableBlock() {
  const { node } = useLibraryPageContext();
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

  const columns = [
    {
      name: "Name",
      fid: "name",
      sort: "0",
      type: "text",
    } satisfies PropertySchema,
    ...node.child_property_schema,
  ];

  return (
    <Table.Root size="sm" variant="dense">
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
                onClick={handleClickSortAction}
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
              >
                <Box
                  display="inline-flex"
                  alignItems="center"
                  gap="1"
                  flexWrap="nowrap"
                >
                  {property.name}
                  <SortIndicator order={sortState} />
                </Box>
              </Table.Header>
            );
          })}
        </Table.Row>
      </Table.Head>
      <Table.Body>
        <SortableContext items={nodes}>
          {nodes.map((child) => (
            <Table.Row key={child.id}>
              <Table.Cell fontWeight="medium">{child.name}</Table.Cell>
              {node.child_property_schema.map((property) => {
                const value = child.properties.find(
                  (p) => p.fid === property.fid,
                )?.value;

                console.log(child.properties, property);

                return <Table.Header key={property.fid}>{value}</Table.Header>;
              })}
            </Table.Row>
          ))}
        </SortableContext>
      </Table.Body>
    </Table.Root>
  );
}
