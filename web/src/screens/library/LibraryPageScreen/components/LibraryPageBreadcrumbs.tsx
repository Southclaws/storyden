import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
import { CancelAction } from "@/components/site/Action/Cancel";
import { EditAction } from "@/components/site/Action/Edit";
import { SaveAction } from "@/components/site/Action/Save";
import { HStack, WStack } from "@/styled-system/jsx";

import { useLibraryPath } from "../../useLibraryPath";
import { useLibraryPageContext } from "../Context";
import { useLibraryPagePermissions } from "../permissions";
import { useEditState } from "../useEditState";

export function LibraryPageControls() {
  const libraryPath = useLibraryPath();
  const { node, form } = useLibraryPageContext();
  const { editing, handleToggleEditMode } = useEditState();
  const { isAllowedToEdit } = useLibraryPagePermissions(node);

  return (
    <WStack alignItems="start">
      <Breadcrumbs
        libraryPath={libraryPath}
        visibility={node.visibility}
        create={editing ? "edit" : "show"}
        defaultValue={node.slug}
        {...form.register("slug")}
      />

      {isAllowedToEdit && (
        <HStack>
          {editing ? (
            <>
              <CancelAction type="button" onClick={handleToggleEditMode}>
                Cancel
              </CancelAction>
              <SaveAction type="submit">Save</SaveAction>
            </>
          ) : (
            <>
              <EditAction onClick={handleToggleEditMode}>Edit</EditAction>
            </>
          )}
          <LibraryPageMenu node={node} />
        </HStack>
      )}
    </WStack>
  );
}
