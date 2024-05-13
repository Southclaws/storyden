"use client";

import { Unready } from "src/components/site/Unready";

import { DatagraphNodeScreen } from "../DatagraphNodeScreen/DatagraphNodeScreen";

import { Props, useNodeViewerScreen } from "./useNodeViewerScreen";

export function NodeViewerScreen(props: Props) {
  const { ready, data, handlers, error } = useNodeViewerScreen(props);

  if (!ready) return <Unready {...error} />;

  return (
    <DatagraphNodeScreen
      node={data}
      onSave={handlers.handleSave}
      onVisibilityChange={handlers.handleVisibilityChange}
      onDelete={handlers.handleDelete}
    />
  );
}
