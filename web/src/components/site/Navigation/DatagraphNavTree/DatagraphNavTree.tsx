"use client";

import { useNodeList } from "@/api/openapi-client/nodes";
import { DatagraphNodeTree } from "@/components/directory/datagraph/DatagraphNodeTree/DatagraphNodeTree";

type Props = {
  currentNode: string | undefined;
};

export function DatagraphNavTree({ currentNode }: Props) {
  const { data } = useNodeList();

  if (!data) return null;

  return (
    <DatagraphNodeTree
      currentNode={currentNode}
      data={{
        label: "Directory",
        children: data.nodes,
      }}
    />
  );
}
