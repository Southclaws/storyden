import { InputEvent, useState } from "react";

import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
import { CancelAction } from "@/components/site/Action/Cancel";
import { EditAction } from "@/components/site/Action/Edit";
import { isSlugReady, processMarkInput } from "@/lib/mark/mark";
import { HStack, WStack } from "@/styled-system/jsx";

import { useLibraryPath } from "../useLibraryPath";

import { useLibraryPageContext } from "./Context";
import { useLibraryPagePermissions } from "./permissions";
import { useWatch } from "./store";
import { useEditState } from "./useEditState";

function useLibraryPageControls() {
  const libraryPath = useLibraryPath();
  const { store } = useLibraryPageContext();
  const { draft, setSlug } = store.getState();

  const slug = useWatch((s) => s.draft.slug);
  const visibility = useWatch((s) => s.draft.visibility);

  // Ensure the final item is the real slug, not the cached copy
  const updatedLibraryPath = [...libraryPath.slice(0, -1), slug];

  const { isAllowedToEdit } = useLibraryPagePermissions();
  const { editing } = useEditState();

  const [isSlugInvalid, setSlugInvalid] = useState(true);

  function handleSlugChange(event: InputEvent<HTMLInputElement>) {
    const raw = event.currentTarget.value;
    const slug = processMarkInput(raw);
    const valid = isSlugReady(slug);
    setSlug(slug);
    setSlugInvalid(!valid);
  }

  return {
    libraryPath: updatedLibraryPath,
    draft,
    slug,
    isSlugInvalid,
    visibility,
    setSlug,
    isAllowedToEdit,
    editing,
    handleSlugChange,
  };
}

export function LibraryPageControls() {
  const {
    libraryPath,
    draft,
    slug,
    visibility,
    isSlugInvalid,
    isAllowedToEdit,
    editing,
    handleSlugChange,
  } = useLibraryPageControls();

  return (
    <WStack alignItems="start">
      <Breadcrumbs
        libraryPath={libraryPath}
        visibility={visibility}
        create={editing ? "edit" : "show"}
        defaultValue={slug}
        value={slug}
        invalid={isSlugInvalid}
        onChange={handleSlugChange}
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
