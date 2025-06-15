import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
import { CancelAction } from "@/components/site/Action/Cancel";
import { EditAction } from "@/components/site/Action/Edit";
import { HStack, WStack } from "@/styled-system/jsx";

import { useLibraryPath } from "../useLibraryPath";

import { useLibraryPageContext } from "./Context";
import { useLibraryPagePermissions } from "./permissions";
import { useWatch } from "./store";
import { useEditState } from "./useEditState";

export function LibraryPageControls() {
  const libraryPath = useLibraryPath();
  const { store } = useLibraryPageContext();
  const { draft, setSlug } = store.getState();

  const slug = useWatch((s) => s.draft.slug);
  const visibility = useWatch((s) => s.draft.visibility);

  const { isAllowedToEdit } = useLibraryPagePermissions();
  const { editing } = useEditState();

  return (
    <WStack alignItems="start">
      <Breadcrumbs
        libraryPath={libraryPath}
        visibility={visibility}
        create={editing ? "edit" : "show"}
        defaultValue={slug}
        onChange={(slug) => setSlug(slug.currentTarget.value)}
      />
      <HStack>
        {isAllowedToEdit && <EditControls />}
        <LibraryPageMenu node={draft} />
      </HStack>
    </WStack>
  );
}

function EditControls() {
  const { editing, saving, handleToggleEditMode } = useEditState();

  if (!editing) {
    return (
      <EditAction type="button" onClick={handleToggleEditMode}>
        Edit
      </EditAction>
    );
  }

  return (
    <>
      <CancelAction
        type="button"
        loading={saving}
        disabled={saving}
        onClick={handleToggleEditMode}
      >
        View
      </CancelAction>
    </>
  );
}
