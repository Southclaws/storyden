"use client";

import { ClusterScreen } from "../ClusterScreen/ClusterScreen";

import { Props, useClusterCreateScreen } from "./useClusterCreateScreen";

export function ClusterCreateScreen(props: Props) {
  const {
    handlers: { handleCreate },
    initial,
  } = useClusterCreateScreen(props);

  return (
    <ClusterScreen
      cluster={initial}
      initialEditingState={true}
      onSave={handleCreate}
    />
  );
}
