import { parseAsBoolean, useQueryState } from "nuqs";

import { useLibraryPagePermissions } from "./permissions";

export function useEditState() {
  const [editing, setEditing] = useQueryState("edit", {
    ...parseAsBoolean,
    defaultValue: false,
    clearOnDefault: true,
  });

  const { isAllowedToEdit } = useLibraryPagePermissions();

  function handleToggleEditMode() {
    if (editing) {
      setEditing(false);
    } else {
      if (!isAllowedToEdit) return;

      setEditing(true);
    }
  }

  return {
    editing,
    setEditing,
    handleToggleEditMode,
  };
}
