import {
  Identifier,
  NodeWithChildren,
  PropertyName,
  PropertySchemaList,
  PropertyType,
  PropertyValue,
} from "@/api/openapi-schema";
import {
  LibraryPageBlockTypeTable,
  LibraryPageBlockTypeTableConfig,
} from "@/lib/library/metadata";

export type ColumnDefinitionCommon = {
  fid: Identifier;
  name: PropertyName;
  type: PropertyType;
  hidden: boolean;
};

export type ColumnDefinitionProperty = ColumnDefinitionCommon & {
  _fixedFieldName?: undefined;
};

export type ColumnDefinitionFixed = ColumnDefinitionCommon & {
  fid: `fixed:${MappableNodeField}`;
  name: MappableNodeField; // TODO: Pretty names for fixed fields
  _fixedFieldName: MappableNodeField;
};

export type ColumnDefinition = ColumnDefinitionProperty | ColumnDefinitionFixed;

export type ColumnValue = {
  fid: Identifier;
  value: PropertyValue;
};

export type MappableNodeField = Extract<
  "slug" | "name" | "link" | "description",
  keyof NodeWithChildren
>;

export const MappableNodeFields: Array<MappableNodeField> = [
  "slug",
  "name",
  "link",
  "description",
];

export type ProcessedConfig = {
  columns: Array<ColumnDefinition>;
};

export function getDefaultBlockConfig(ps: PropertySchemaList): ProcessedConfig {
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

export function processBlockConfig(
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

export function mergeFieldsAndPropertySchema(
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

export function mergeFieldsAndProperties(
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
