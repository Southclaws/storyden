import { DatagraphItem, DatagraphItemKind } from "@/api/openapi-schema";

export function getCommonProperties(item: DatagraphItem) {
  switch (item.kind) {
    case DatagraphItemKind.post:
      return {
        name: item.ref.title,
        description: item.ref.description,
        slug: item.ref.slug,
      };
    case DatagraphItemKind.thread:
      return {
        name: item.ref.title,
        description: item.ref.description,
        slug: item.ref.slug,
      };
    case DatagraphItemKind.reply:
      return {
        name: item.ref.title,
        description: item.ref.description,
        slug: item.ref.slug,
      };
    case DatagraphItemKind.node:
      return {
        name: item.ref.name,
        description: item.ref.description,
        slug: item.ref.slug,
      };
    case DatagraphItemKind.profile:
      return {
        name: item.ref.name,
        description: item.ref.bio,
        slug: item.ref.handle,
      };
  }
}
