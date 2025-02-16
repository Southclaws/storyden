"use client";

import Image from "next/image";
import { FixedCropper, ImageRestriction } from "react-advanced-cropper";
import { FormProvider } from "react-hook-form";

import { CancelAction } from "src/components/site/Action/Cancel";
import { EditAction } from "src/components/site/Action/Edit";
import { SaveAction } from "src/components/site/Action/Save";

import { useNodeGet } from "@/api/openapi-client/nodes";
import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { LibraryPageAssetList } from "@/components/library/LibraryPageAssetList/LibraryPageAssetList";
import { LibraryPageCoverImageControl } from "@/components/library/LibraryPageCoverImageControl/LibraryPageCoverImageControl";
import { LibraryPageImportFromURL } from "@/components/library/LibraryPageImportFromURL/LibraryPageImportFromURL";
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
import { LibraryPagePropertyTable } from "@/components/library/LibraryPagePropertyTable/LibraryPagePropertyTable";
import { LibraryPageTagsList } from "@/components/library/LibraryPageTagsList/LibraryPageTagsList";
import { IntelligenceAction } from "@/components/site/Action/Intelligence";
import { UnreadyBanner } from "@/components/site/Unready";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import { css } from "@/styled-system/css";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

import "react-advanced-cropper/dist/style.css";

import { ContentInput } from "./ContentInput";
import { TitleInput } from "./TitleInput";
import {
  CROP_STENCIL_HEIGHT,
  CROP_STENCIL_WIDTH,
  Form,
  Props,
  useLibraryPageScreen,
} from "./useLibraryPageScreen";

export function LibraryPageScreen(props: Props) {
  const { data, error } = useNodeGet(props.node.slug, {
    swr: { fallbackData: props.node },
  });
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  // NOTE: There's a bug in SWR here where if the fallback data for an array
  // is passed as empty, it becomes undefined. Maybe cache or mutate related?
  data.tags = data.tags ?? [];

  return <LibraryPage node={data} />;
}

export function LibraryPage(props: Props) {
  const {
    form,
    handlers: {
      handleSubmit,
      handleEditMode,
      handleSuggestTitle,
      handleResetGeneratedTitle,
      handleAssetUpload,
      handleImportFromLink,
    },
    libraryPath,
    editing,
    node,
    generatedContent,
    generatedTitle,
    cropperRef,
    primaryAssetURL,
    primaryAssetEditingURL,
    initialCoverCoordinates,
    isAllowedToEdit,
    isTitleSuggestEnabled,
    isLoadingSuggestTitle,
  } = useLibraryPageScreen(props);

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
      <FormProvider {...form}>
        <LStack h="full">
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
                    <CancelAction type="button" onClick={handleEditMode}>
                      Cancel
                    </CancelAction>
                    <SaveAction type="submit">Save</SaveAction>
                  </>
                ) : (
                  <>
                    <EditAction onClick={handleEditMode}>Edit</EditAction>
                  </>
                )}
                <LibraryPageMenu node={node} />
              </HStack>
            )}
          </WStack>

          {editing && (
            <HStack w="full" justify="end">
              {/* TODO: Icons/emojis custom for pages too */}
              {/* <Button size="xs" variant="outline">
                <SmilePlusIcon /> page icon
              </Button> */}
              <LibraryPageCoverImageControl node={node} />
              {/* TODO: Import from other sources */}
              {/* <Button size="xs" variant="outline">
                <ImportIcon /> import
              </Button> */}
            </HStack>
          )}

          {editing && primaryAssetEditingURL ? (
            <Box width="full" height="64">
              <FixedCropper
                ref={cropperRef}
                className={css({
                  maxWidth: "full",
                  maxHeight: "64",
                  borderRadius: "md",
                  // TODO: Remove black background when empty
                  backgroundColor: "bg.default",
                })}
                defaultPosition={
                  initialCoverCoordinates && {
                    top: initialCoverCoordinates.top,
                    left: initialCoverCoordinates.left,
                  }
                }
                backgroundWrapperProps={{
                  scaleImage: false,
                }}
                stencilProps={{
                  handlers: false,
                  lines: false,
                  movable: false,
                  resizable: false,
                }}
                stencilSize={{
                  width: CROP_STENCIL_WIDTH,
                  height: CROP_STENCIL_HEIGHT,
                }}
                imageRestriction={ImageRestriction.stencil}
                src={primaryAssetEditingURL}
              />
            </Box>
          ) : (
            primaryAssetURL && (
              <Box height="64" width="full">
                <Image
                  className={css({
                    width: "full",
                    height: "full",
                    borderRadius: "md",
                    objectFit: "cover",
                    objectPosition: "center",
                  })}
                  src={primaryAssetURL}
                  alt=""
                  width={CROP_STENCIL_WIDTH}
                  height={CROP_STENCIL_HEIGHT}
                />
              </Box>
            )
          )}

          <LibraryPageAssetList node={node} />

          <LStack gap="2">
            <LStack minW="0">
              <WStack alignItems="end">
                {editing ? (
                  <>
                    <TitleInput
                      imperativeValue={generatedTitle}
                      onResetImperativeValue={handleResetGeneratedTitle}
                    />
                    {isTitleSuggestEnabled && (
                      <IntelligenceAction
                        title="Suggest a title for this page"
                        onClick={handleSuggestTitle}
                        variant="subtle"
                        h="full"
                        loading={isLoadingSuggestTitle}
                      />
                    )}
                  </>
                ) : (
                  <Heading fontSize="heading.2" fontWeight="bold">
                    {node.name || "(untitled)"}
                  </Heading>
                )}
              </WStack>
            </LStack>
          </LStack>

          <HStack
            w="full"
            flexDirection={{
              base: "column-reverse",
              sm: "row",
            }}
            alignItems={{
              base: "start",
              sm: "center",
            }}
          >
            {!editing && node.link?.url && (
              <LinkButton href={node.link?.url} size="xs" variant="subtle">
                {node.link?.domain}
              </LinkButton>
            )}

            <LibraryPageTagsList<Form>
              control={form.control}
              name="tags"
              editing={editing}
              node={node}
            />
          </HStack>

          {editing && (
            <LibraryPageImportFromURL
              control={form.control}
              name="link"
              node={node}
              onImport={handleImportFromLink}
            />
          )}

          <LibraryPagePropertyTable
            control={form.control}
            name="properties"
            editing={editing}
            node={node}
          />

          <ContentInput
            // TODO: Fix this via ref
            // value={form.getValues().content}
            disabled={!editing}
            onAssetUpload={handleAssetUpload}
            initialValue={
              node.content ?? form.formState.defaultValues?.["content"]
            }
            value={generatedContent}
          />
        </LStack>
      </FormProvider>
    </styled.form>
  );
}
