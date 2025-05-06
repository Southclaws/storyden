import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
import { CancelAction } from "@/components/site/Action/Cancel";
import { EditAction } from "@/components/site/Action/Edit";
import { SaveAction } from "@/components/site/Action/Save";
import { HStack, WStack } from "@/styled-system/jsx";

import { useLibraryPath } from "../useLibraryPath";

import { useLibraryPageContext } from "./Context";
import { useLibraryPagePermissions } from "./permissions";
import { useEditState } from "./useEditState";

export function LibraryPageControls() {
  const libraryPath = useLibraryPath();
  const { node, form } = useLibraryPageContext();
  const { isAllowedToEdit } = useLibraryPagePermissions(node);
  const { editing } = useEditState();

  return (
    <WStack alignItems="start">
      <Breadcrumbs
        libraryPath={libraryPath}
        visibility={node.visibility}
        create={editing ? "edit" : "show"}
        defaultValue={node.slug}
        {...form.register("slug")}
      />

      <HStack>
        {isAllowedToEdit && <EditControls />}
        <LibraryPageMenu node={node} />
      </HStack>
    </WStack>
  );
}

function EditControls() {
  const { editing, handleToggleEditMode } = useEditState();

  if (!editing) {
    return <EditAction onClick={handleToggleEditMode}>Edit</EditAction>;
  }

  return (
    <>
      <CancelAction type="button" onClick={handleToggleEditMode}>
        Cancel
      </CancelAction>
      <SaveAction type="submit">Save</SaveAction>
    </>
  );
}
