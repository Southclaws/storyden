import slugify from "@sindresorhus/slugify";

import { LinkIcon } from "@/components/ui/icons/Link";
import * as Menu from "@/components/ui/menu";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";

export function LibraryPageTitleBlockMenuItems() {
  const { store } = useLibraryPageContext();
  const { setSlug } = store.getState();

  const name = useWatch((s) => s.draft.name);

  const canUpdateSlug = Boolean(name?.trim());

  function handleUpdateSlug() {
    if (name) {
      const newSlug = slugify(name);
      setSlug(newSlug);
    }
  }

  if (!canUpdateSlug) {
    return null;
  }

  return (
    <Menu.Item value="update-slug" onClick={handleUpdateSlug}>
      <LinkIcon />
      &nbsp;Update URL slug
    </Menu.Item>
  );
}
