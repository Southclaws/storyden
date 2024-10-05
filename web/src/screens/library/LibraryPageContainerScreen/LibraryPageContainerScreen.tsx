"use client";

import { Unready } from "src/components/site/Unready";

import { LibraryPageScreen } from "../LibraryPageScreen/LibraryPageScreen";

import {
  Props,
  useLibraryPageContainerScreen,
} from "./useLibraryPageContainerScreen";

export function LibraryPageContainerScreen(props: Props) {
  const { ready, data, handlers, error } = useLibraryPageContainerScreen(props);

  if (!ready) return <Unready error={error} />;

  return (
    <LibraryPageScreen
      node={data}
      onSave={handlers.handleSave}
      onVisibilityChange={handlers.handleVisibilityChange}
      onDelete={handlers.handleDelete}
    />
  );
}
