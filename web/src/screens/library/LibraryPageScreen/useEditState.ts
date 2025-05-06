import { parseAsBoolean, useQueryState } from "nuqs";

import { useLibraryPageContext } from "./Context";
import { useLibraryPagePermissions } from "./permissions";

export function useEditState() {
  const [editing, setEditing] = useQueryState("edit", {
    ...parseAsBoolean,
    defaultValue: false,
    clearOnDefault: true,
  });

  const { node, form, defaultFormValues } = useLibraryPageContext();

  const { isAllowedToEdit } = useLibraryPagePermissions(node);

  function handleToggleEditMode() {
    if (editing) {
      setEditing(false);

      form.reset(defaultFormValues);
    } else {
      if (!isAllowedToEdit) return;

      setEditing(true);

      form.reset(defaultFormValues);
    }
  }

  return {
    editing,
    setEditing,
    handleToggleEditMode,
  };
}
