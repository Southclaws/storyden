"use client";

import { Unready } from "src/components/site/Unready";

import { DatagraphNodeScreen } from "../DatagraphNodeScreen/DatagraphNodeScreen";

import { Props, useClusterViewerScreen } from "./useClusterViewerScreen";

export function ClusterViewerScreen(props: Props) {
  const { ready, data, handlers, error } = useClusterViewerScreen(props);

  if (!ready) return <Unready {...error} />;

  return (
    <DatagraphNodeScreen
      node={{ type: "cluster", ...data }}
      onSave={handlers.handleSave}
      onDelete={handlers.handleDelete}
    />
  );
}
