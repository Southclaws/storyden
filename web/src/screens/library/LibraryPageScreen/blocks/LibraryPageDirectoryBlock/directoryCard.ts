import { ReactNode } from "react";

import { Asset, Identifier, NodeWithChildren } from "@/api/openapi-schema";
import { LibraryPageBlockTypeDirectoryConfig } from "@/lib/library/metadata";

import { ColumnDefinition, ColumnValue, MappableNodeField } from "./column";

export type DirectoryCardLayout = "rows" | "grid" | "gallery";

export type DirectoryPropertyProps = {
  node: NodeWithChildren;
  columns: ColumnDefinition[];
  values: Record<string, ColumnValue>;
  editing: boolean;
  onFieldValueChange: (nodeID: Identifier, fid: string, value: string) => void;
  renderPropertyLabel?: (column: ColumnDefinition) => ReactNode;
};

export type DirectoryCardProps = DirectoryPropertyProps & {
  layout: DirectoryCardLayout;
};

export type DirectoryRichCardProps = DirectoryCardProps & {
  config: LibraryPageBlockTypeDirectoryConfig;
};

export function applyDirectoryDraft(
  source: NodeWithChildren,
  draft: NodeWithChildren | undefined,
): NodeWithChildren {
  if (!draft) {
    return source;
  }

  // Draft children can omit media and content from the full directory entry.
  // Only the fields edited by directory cards belong to this overlay.
  return {
    ...source,
    name: draft.name,
    description: draft.description,
    link: draft.link ?? source.link,
    properties: draft.properties,
  };
}

export function getDirectoryBadgeColumns(
  columns: ColumnDefinition[],
  values: Record<string, ColumnValue>,
  editing: boolean,
  placement: ColumnDefinition["placement"],
) {
  return columns.filter((column) => {
    if (column.placement !== placement) {
      return false;
    }

    if (editing) {
      return true;
    }

    return Boolean(values[column.fid]?.value);
  });
}

export type RichCardShape = "minimal" | "row" | "box" | "fill";

export function getRichCardShape(
  layout: DirectoryCardLayout,
  style: LibraryPageBlockTypeDirectoryConfig["card_style"],
  hasImage: boolean,
): RichCardShape {
  if (layout === "rows") {
    return "row";
  }

  if (style === undefined || style === "minimal") {
    return "minimal";
  }

  if (layout === "gallery") {
    return "box";
  }

  if (style === "row") {
    return "row";
  }

  if (style === "fill" && hasImage) {
    return "fill";
  }

  return "box";
}

export function getRichCardImageAspectRatio(
  layout: DirectoryCardLayout,
  shape: RichCardShape,
  image: Pick<Asset, "width" | "height"> | undefined,
): number | undefined {
  if (shape === "minimal" && layout === "gallery") {
    return undefined;
  }

  if (shape !== "box" && shape !== "minimal") {
    return undefined;
  }

  if (!image) {
    return undefined;
  }

  if (layout !== "gallery") {
    return 4 / 3;
  }

  if (image.width <= 0) {
    return 4 / 3;
  }

  if (image.height <= 0) {
    return 4 / 3;
  }

  return image.width / image.height;
}

export function isRichCardFieldVisible(
  config: Pick<LibraryPageBlockTypeDirectoryConfig, "columns">,
  field: MappableNodeField,
  shape: RichCardShape = "box",
): boolean {
  const column = config.columns.find(
    (column) => column.fid === `fixed:${field}`,
  );

  if (column) {
    return !column.hidden;
  }

  if (field === "primary_image") {
    return true;
  }

  return shape !== "minimal";
}
