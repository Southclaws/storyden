"use client";

import { Unready } from "src/components/site/Unready";

import { ClusterScreen } from "../ClusterScreen/ClusterScreen";

import { Props, useClusterViewerScreen } from "./useClusterViewerScreen";

export function ClusterViewerScreen(props: Props) {
  const { ready, data, handlers, error } = useClusterViewerScreen(props);

  if (!ready) return <Unready {...error} />;

  return <ClusterScreen cluster={data} onSave={handlers.handleSave} />;
}
