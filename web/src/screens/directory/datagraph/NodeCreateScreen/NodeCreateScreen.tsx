"use client";

import { DatagraphNodeScreen } from "../DatagraphNodeScreen/DatagraphNodeScreen";

import { Props, useNodeCreateScreen } from "./useNodeCreateScreen";

export function NodeCreateScreen(props: Props) {
  const {
    handlers: { handleCreate },
    initial,
  } = useNodeCreateScreen(props);

  return (
    <DatagraphNodeScreen
      node={initial}
      initialEditingState={true}
      onSave={handleCreate}
    />
  );
}
