import { InputEvent, useState } from "react";

import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
import { ButtonGroup } from "@/components/ui/button";
import { isSlugReady, processMarkInput } from "@/lib/mark/mark";
import { HStack, WStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "./Context";
import { LibraryPageEditMenu } from "./LibraryPageEditMenu";
import { useWatch } from "./store";
import { useEditState } from "./useEditState";

function useLibraryPageControls() {
  const { ancestors, store } = useLibraryPageContext();
  const { draft, setSlug } = store.getState();

  const current = useWatch((s) => ({
    id: s.draft.id,
    name: s.draft.name,
    slug: s.draft.slug,
  }));
  const visibility = useWatch((s) => s.draft.visibility);

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
    ancestors,
    current,
    draft,
    isSlugInvalid,
    visibility,
    setSlug,
    editing,
    handleSlugChange,
  };
}

export function LibraryPageControls() {
  const {
    ancestors,
    current,
    draft,
    visibility,
    isSlugInvalid,
    editing,
    handleSlugChange,
  } = useLibraryPageControls();

  return (
    <WStack alignItems="start">
      <Breadcrumbs
        ancestors={ancestors}
        current={current}
        visibility={visibility}
        create={editing ? "edit" : "show"}
        defaultValue={current.slug}
        value={current.slug}
        invalid={isSlugInvalid}
        onChange={handleSlugChange}
      />

      <HStack>
        <ButtonGroup variant="subtle" attached>
          <EditMenuControls />
          <LibraryPageMenu node={draft} />
        </ButtonGroup>
      </HStack>
    </WStack>
  );
}

function EditMenuControls() {
  const { initialNode } = useLibraryPageContext();

  return <LibraryPageEditMenu node={initialNode} />;
}
