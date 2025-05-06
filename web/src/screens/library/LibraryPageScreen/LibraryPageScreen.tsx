"use client";

import { useNodeGet } from "@/api/openapi-client/nodes";
import { UnreadyBanner } from "@/components/site/Unready";
import { LStack, styled } from "@/styled-system/jsx";

import "react-advanced-cropper/dist/style.css";

import { LibraryPageProvider, Props } from "./Context";
import { LibraryPageAssetsBlock } from "./blocks/LibraryPageAssetsBlock/LibraryPageAssetsBlock";
import { LibraryPageContentBlock } from "./blocks/LibraryPageContentBlock/LibraryPageContentBlock";
import { LibraryPageCoverBlock } from "./blocks/LibraryPageCoverBlock/LibraryPageCoverBlock";
import { LibraryPageLinkBlock } from "./blocks/LibraryPageLinkBlock/LibraryPageLinkBlock";
import { LibraryPageTagsBlock } from "./blocks/LibraryPageTagsBlock/LibraryPageTagsBlock";
import { LibraryPageTitleBlock } from "./blocks/LibraryPageTitleBlock/LibraryPageTitleBlock";
import { LibraryPagePropertiesBlock } from "./blocks/LibraryPropertiesBlock/LibraryPropertiesBlock";
import { LibraryPageControls } from "./components/LibraryPageBreadcrumbs";
import { useCoverImage } from "./useCoverImage";
import { useSave } from "./useSave";

export function LibraryPageScreen(props: Props) {
  const { data, error } = useNodeGet(props.node.slug, undefined, {
    swr: { fallbackData: props.node },
  });
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  // NOTE: There's a bug in SWR here where if the fallback data for an array
  // is passed as empty, it becomes undefined. Maybe cache or mutate related?
  data.tags = data.tags ?? [];

  return <LibraryPageForm node={data} />;
}

function LibraryPageForm(props: Props) {
  return (
    <LibraryPageProvider node={props.node}>
      <LibraryPage />
    </LibraryPageProvider>
  );
}

export function LibraryPage() {
  const { cropperRef, handleUploadCroppedCover } = useCoverImage();

  const { handleSubmit } = useSave({
    handleUploadCroppedCover: handleUploadCroppedCover,
  });

  return (
    <styled.form
      display="flex"
      flexDir="column"
      w="full"
      h="full"
      gap="3"
      alignItems="start"
      onSubmit={handleSubmit}
    >
      <LStack h="full">
        <LibraryPageControls />
        <LibraryPageCoverBlock ref={cropperRef} />
        <LibraryPageAssetsBlock />
        <LibraryPageTitleBlock />
        <LibraryPageTagsBlock />
        <LibraryPageLinkBlock />
        <LibraryPagePropertiesBlock />
        <LibraryPageContentBlock />
      </LStack>
    </styled.form>
  );
}
