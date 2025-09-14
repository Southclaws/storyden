import { MenuSelectionDetails } from "@ark-ui/react";

import { MediaRemoveIcon } from "@/components/ui/icons/Media";
import * as Menu from "@/components/ui/menu";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";

export function LibraryPageCoverBlockMenuItems() {
  const { store } = useLibraryPageContext();
  const { removePrimaryImage } = store.getState();
  const primary_image = useWatch((s) => s.draft.primary_image);

  const hasCoverImage = !!primary_image;

  function handleRemoveCover() {
    removePrimaryImage();
  }

  if (!hasCoverImage) {
    return null;
  }

  return (
    <Menu.Item value="remove-cover" onSelect={handleRemoveCover}>
      <MediaRemoveIcon />
      &nbsp;Remove cover image
    </Menu.Item>
  );
}