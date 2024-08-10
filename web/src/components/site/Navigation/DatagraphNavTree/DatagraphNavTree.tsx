"use client";

import { reduce } from "lodash/fp";

import { useNodeList } from "@/api/openapi-client/nodes";
import { NodeWithChildren } from "@/api/openapi-schema";
import { Child, TreeView, TreeViewData } from "@/components/ui/tree-view";

const recursivelyMapChildren = reduce<NodeWithChildren, Child[]>(
  (prev: Child[], curr: NodeWithChildren) => {
    const next = {
      value: curr.slug,
      name: curr.name,
      url: `/directory/${curr.slug}`,
      children: recursivelyMapChildren(curr.children),
    } satisfies Child;

    return [...prev, next];
  },
  [],
);

export function DatagraphNavTree() {
  const { data } = useNodeList();

  if (!data) return null;

  const children: Child[] = recursivelyMapChildren(data.nodes);

  const root = {
    label: "Root",
    children,
  } satisfies TreeViewData;

  return <TreeView data={root} />;
}
