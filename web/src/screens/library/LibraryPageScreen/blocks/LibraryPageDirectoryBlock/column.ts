import { keyBy } from "lodash";

import {
  Identifier,
  NodeWithChildren,
  PropertyName,
  PropertySchemaList,
  PropertyType,
  PropertyValue,
} from "@/api/openapi-schema";
import {
  LibraryPageBlockTypeDirectory,
  LibraryPageBlockTypeDirectoryConfig,
  LibraryPageBlockTypeDirectoryLayout,
} from "@/lib/library/metadata";

export type ColumnDefinitionCommon = {
  fid: Identifier;
  name: PropertyName;
  type: PropertyType;
  hidden: boolean;
};

export type ColumnDefinitionProperty = ColumnDefinitionCommon & {
  fixed: false;
  _fixedFieldName?: undefined;
};

export type ColumnDefinitionFixed = ColumnDefinitionCommon & {
  fixed: true;
  fid: `fixed:${MappableNodeField}`;
  name: MappableNodeField; // TODO: Pretty names for fixed fields
  _fixedFieldName: MappableNodeField;
};

export type ColumnDefinition = ColumnDefinitionProperty | ColumnDefinitionFixed;

export type ColumnValue = {
  fid: Identifier;
  value: PropertyValue;
  href?: string;
};

export type MappableNodeField = Extract<
  "name" | "link" | "description" | "primary_image",
  keyof NodeWithChildren
>;

export const MappableNodeFields: Array<MappableNodeField> = [
  "name",
  "link",
  "description",
  "primary_image",
];

export type ProcessedConfig = {
  layout: LibraryPageBlockTypeDirectoryLayout;
  columns: Array<ColumnDefinition>;
};

export function getDefaultBlockConfig(ps: PropertySchemaList): ProcessedConfig {
  const fixedCols: ColumnDefinitionFixed[] = MappableNodeFields.map((field) => {
    return {
      fixed: true,
      fid: `fixed:${field}`,
      name: field,
      type: "text" as const,
      hidden: false,
      _fixedFieldName: field as MappableNodeField,
    } satisfies ColumnDefinitionFixed;
  });

  const propCols: ColumnDefinitionProperty[] = ps.map((property) => ({
    fixed: false,
    fid: property.fid,
    name: property.name,
    type: property.type,
    sort: property.sort,
    hidden: false,
  }));

  const columns: ColumnDefinition[] = [...fixedCols, ...propCols];

  return {
    layout: "table",
    columns,
  };
}

export function processBlockConfig(
  ps: PropertySchemaList,
  config: LibraryPageBlockTypeDirectoryConfig,
  showHidden = false,
): ProcessedConfig {
  const schemaMap = keyBy(ps, "fid");

  const columns: ColumnDefinition[] = config.columns.reduce((prev, column) => {
    if (showHidden === false && column.hidden) {
      return prev;
    }

    if (column.fid.startsWith("fixed:")) {
      const fixedColumn = {
        fixed: true,
        fid: column.fid as `fixed:${MappableNodeField}`,
        name: column.fid.replace("fixed:", "") as MappableNodeField,
        type: "text" as const, // Fixed fields are always text for now
        hidden: column.hidden,
        _fixedFieldName: column.fid.replace("fixed:", "") as MappableNodeField,
      } satisfies ColumnDefinitionFixed;

      return [...prev, fixedColumn];
    }

    const schema = schemaMap[column.fid];
    if (!schema) {
      return prev;
    }

    const next = {
      fixed: false,
      fid: schema.fid,
      name: schema.name,
      type: schema.type as PropertyType,
      hidden: column.hidden ?? false,
    } satisfies ColumnDefinitionProperty;

    return [...prev, next];
  }, [] as ColumnDefinition[]);

  return {
    ...config,
    columns,
  };
}

export function mergeFieldsAndPropertySchema(
  ps: PropertySchemaList,
  block: LibraryPageBlockTypeDirectory,
  showHidden = false,
): ColumnDefinition[] {
  const config =
    block.config === undefined
      ? getDefaultBlockConfig(ps)
      : processBlockConfig(ps, block.config, showHidden);

  return config.columns;
}

export function mergeFieldsAndProperties(
  schema: PropertySchemaList,
  node: NodeWithChildren,
  block: LibraryPageBlockTypeDirectory,
): ColumnValue[] {
  const config =
    block.config === undefined
      ? getDefaultBlockConfig(schema)
      : processBlockConfig(schema, block.config);

  const columns: ColumnValue[] = config.columns.reduce((prev, column) => {
    if (column.hidden) {
      return prev;
    }

    if (column._fixedFieldName) {
      if (column._fixedFieldName === "link") {
        return [
          ...prev,
          {
            fid: column.fid,
            value: node.link?.url ?? "",
            href: node.link?.url ?? "",
          },
        ];
      }

      if (column._fixedFieldName === "name") {
        return [
          ...prev,
          {
            fid: column.fid,
            value: node.name ?? "",
            href: `/l/${node.slug}`,
          },
        ];
      }

      if (column._fixedFieldName === "primary_image") {
        return [
          ...prev,
          {
            fid: column.fid,
            value: node.primary_image?.path ?? "",
          },
        ];
      }

      const next = {
        fid: column.fid,
        value: node[column._fixedFieldName] ?? "",
      };

      return [...prev, next];
    } else {
      const value = node.properties.find((p) => p.fid === column.fid)?.value;

      const next = {
        fid: column.fid,
        value: value ?? "",
      };

      return [...prev, next];
    }
  }, [] as ColumnValue[]);

  return columns;
}
