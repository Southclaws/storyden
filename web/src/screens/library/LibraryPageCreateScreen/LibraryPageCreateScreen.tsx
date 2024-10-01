"use client";

import { LibraryPageScreen } from "../LibraryPageScreen/LibraryPageScreen";

import {
  Props,
  useLibraryPageCreateScreen,
} from "./useLibraryPageCreateScreen";

export function LibraryPageCreateScreen(props: Props) {
  const {
    handlers: { handleCreate },
    initial,
  } = useLibraryPageCreateScreen(props);

  return (
    <LibraryPageScreen
      node={initial}
      initialEditingState={true}
      onSave={handleCreate}
    />
  );
}
