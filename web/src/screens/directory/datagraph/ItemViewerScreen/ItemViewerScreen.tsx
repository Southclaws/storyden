"use client";

import { Unready } from "src/components/site/Unready";

import { DatagraphNodeScreen } from "../DatagraphNodeScreen/DatagraphNodeScreen";

import { Props, useItemViewerScreen } from "./useItemViewerScreen";

export function ItemViewerScreen(props: Props) {
  const { ready, data, handlers, error } = useItemViewerScreen(props);

  if (!ready) return <Unready {...error} />;

  return (
    <DatagraphNodeScreen
      node={{ type: "item", ...data }}
      onSave={handlers.handleSave}
      onDelete={handlers.handleDelete}
    />
  );
}
