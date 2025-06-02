import { SortableContext } from "@dnd-kit/sortable";

import { useNodeListChildren } from "@/api/openapi-client/nodes";
import {
  Identifier,
  Node,
  NodeWithChildren,
  PropertyName,
  PropertySchema,
  PropertySchemaList,
  PropertyType,
  PropertyValue,
} from "@/api/openapi-schema";
import {
  SortIndicator,
  useSortIndicator,
} from "@/components/site/SortIndicator";
import { Unready } from "@/components/site/Unready";
import * as Table from "@/components/ui/table";
import {
  LibraryPageBlockTypeTable,
  LibraryPageBlockTypeTableConfig,
} from "@/lib/library/metadata";
import { Box } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";

type ColumnDefinitionCommon = {
  fid: Identifier;
  name: PropertyName;
  type: PropertyType;
  hidden: boolean;
};

type ColumnDefinitionProperty = ColumnDefinitionCommon & {
  _fixedFieldName?: undefined;
};

type ColumnDefinitionFixed = ColumnDefinitionCommon & {
  fid: `fixed:${MappableNodeField}`;
  name: MappableNodeField; // TODO: Pretty names for fixed fields
  _fixedFieldName: MappableNodeField;
};

type ColumnDefinition = ColumnDefinitionProperty | ColumnDefinitionFixed;

type ColumnValue = {
  fid: Identifier;
  value: PropertyValue;
};

type MappableNodeField = Extract<
  "slug" | "name" | "link" | "description",
  keyof NodeWithChildren
>;

const MappableNodeFields: Array<MappableNodeField> = [
  "slug",
  "name",
  "link",
  "description",
];

type ProcessedConfig = {
  columns: Array<ColumnDefinition>;
};

function getDefaultBlockConfig(ps: PropertySchemaList): ProcessedConfig {
  const fixedCols: ColumnDefinitionFixed[] = MappableNodeFields.map(
    (field) =>
      ({
        fid: `fixed:${field}`,
        name: field,
        type: "text" as const,
        hidden: false,
        _fixedFieldName: field as MappableNodeField,
      }) satisfies ColumnDefinitionFixed,
  );

  const propCols: ColumnDefinitionProperty[] = ps.map((property) => ({
    fid: property.fid,
    name: property.name,
    type: property.type,
    sort: property.sort,
    hidden: false,
  }));

  const columns: ColumnDefinition[] = [...fixedCols, ...propCols];

  return {
    columns,
  };
}

function processBlockConfig(
  config: LibraryPageBlockTypeTableConfig,
): ProcessedConfig {
  const columns: ColumnDefinition[] = config.columns.map((column) => {
    if (column.fid.startsWith("fixed:")) {
      return {
        fid: column.fid as `fixed:${MappableNodeField}`,
        name: column.name as MappableNodeField, // TODO: Pretty names for fixed fields
        type: "text" as const, // Fixed fields are always text for now
        hidden: column.hidden,
        _fixedFieldName: column.fid.replace("fixed:", "") as MappableNodeField,
      } satisfies ColumnDefinitionFixed;
    }

    return column as ColumnDefinitionProperty;
  });

  return {
    ...config,
    columns,
  };
}

function mergeFieldsAndPropertySchema(
  node: NodeWithChildren,
  block: LibraryPageBlockTypeTable,
): ColumnDefinition[] {
  const config =
    block.config === undefined
      ? getDefaultBlockConfig(node.child_property_schema)
      : processBlockConfig(block.config);

  const columns: ColumnDefinition[] = config.columns.map((column) => {
    const r: ColumnDefinition = {
      fid: column.fid,
      name: column.name,
      type: column.type as PropertyType,
      hidden: column.hidden ?? false,
      _fixedFieldName: undefined,
    };

    return r;
  });

  return columns;
}

function mergeFieldsAndProperties(
  schema: PropertySchemaList,
  node: NodeWithChildren,
  block: LibraryPageBlockTypeTable,
): ColumnValue[] {
  const config =
    block.config === undefined
      ? getDefaultBlockConfig(schema)
      : processBlockConfig(block.config);

  const columns: ColumnValue[] = config.columns.map((column) => {
    if (column._fixedFieldName) {
      if (column._fixedFieldName === "link") {
        return {
          fid: column.fid,
          value: node.link?.url ?? "",
        };
      }

      return {
        fid: column.fid,
        value: node[column._fixedFieldName] ?? "",
      };
    } else {
      const value = node.properties.find((p) => p.fid === column.fid)?.value;

      return {
        fid: column.fid,
        value: value ?? "",
      };
    }
  });

  return columns;
}

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

  const block = currentMeta.layout?.blocks.find((b) => b.type === "table");

  if (!block) {
    console.warn(
      "attempting to render a LibraryPageTableBlock without a block in the form metadata",
    );
    return null;
  }

  const columns = mergeFieldsAndPropertySchema(node, block);

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
          {nodes.map((child) => {
            const columns = mergeFieldsAndProperties(
              node.child_property_schema,
              child,
              block,
            );

            console.log(columns);

            return (
              <Table.Row key={child.id}>
                {columns.map((column) => {
                  return (
                    <Table.Header key={column.fid}>{column.value}</Table.Header>
                  );
                })}
              </Table.Row>
            );
          })}
        </SortableContext>
      </Table.Body>
    </Table.Root>
  );
}
