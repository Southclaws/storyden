"use client";

import { ClusterScreen } from "../ClusterScreen/ClusterScreen";

import { useClusterCreateScreen } from "./useClusterCreateScreen";

export function ClusterCreateScreen() {
  const {
    handlers: { handleCreate },
    initial,
  } = useClusterCreateScreen();

  return (
    <ClusterScreen
      cluster={initial}
      initialEditingState={true}
      onSave={handleCreate}
    />
  );
}
