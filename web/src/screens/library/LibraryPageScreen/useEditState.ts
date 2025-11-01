import { parseAsBoolean, useQueryState } from "nuqs";

import { useLibraryPageContext } from "./Context";
import { useLibraryPagePermissions } from "./permissions";

export function useEditState() {
  const [editing, setEditing] = useQueryState("edit", {
    ...parseAsBoolean,
    defaultValue: false,
    clearOnDefault: true,
  });

  const { saving, revalidate } = useLibraryPageContext();

  const { isAllowedToEdit } = useLibraryPagePermissions();

  function handleToggleEditMode() {
    if (editing) {
      setEditing(false);
      revalidate();
    } else {
      if (!isAllowedToEdit) return;

      setEditing(true);
    }
  }

  return {
    editing,
    saving,
    setEditing,
    handleToggleEditMode,
  };
}
