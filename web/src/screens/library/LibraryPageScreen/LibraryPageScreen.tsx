"use client";

import Image from "next/image";
import { FixedCropper, ImageRestriction } from "react-advanced-cropper";
import { FormProvider } from "react-hook-form";

import { CancelAction } from "src/components/site/Action/Cancel";
import { EditAction } from "src/components/site/Action/Edit";
import { SaveAction } from "src/components/site/Action/Save";

import { useNodeGet } from "@/api/openapi-client/nodes";
import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { LibraryPageCoverImageControl } from "@/components/library/LibraryPageCoverImageControl/LibraryPageCoverImageControl";
import { LibraryPageImportFromURL } from "@/components/library/LibraryPageImportFromURL/LibraryPageImportFromURL";
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
import { NodeCardRows } from "@/components/library/NodeCardList";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import { css } from "@/styled-system/css";
import { Box, HStack, LStack, styled } from "@/styled-system/jsx";

import "react-advanced-cropper/dist/style.css";

import { ContentInput } from "./ContentInput";
import { TitleInput } from "./TitleInput";
import {
  CROP_STENCIL_HEIGHT,
  CROP_STENCIL_WIDTH,
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

  return <LibraryPage node={data} />;
}

export function LibraryPage(props: Props) {
  const {
    form,
    handlers: {
      handleSubmit,
      handleEditMode,
      handleVisibilityChange,
      handleDelete,
      handleAssetUpload,
    },
    libraryPath,
    editing,
    node,
    cropperRef,
    primaryAssetURL,
    primaryAssetEditingURL,
    initialCoverCoordinates,
    isAllowedToEdit,
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
          <HStack w="full" justify="space-between" alignItems="start">
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
                <LibraryPageMenu
                  node={node}
                  onVisibilityChange={handleVisibilityChange}
                  onDelete={handleDelete}
                />
              </HStack>
            )}
          </HStack>

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
                  borderRadius: "lg",
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
                    borderRadius: "lg",
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

          <LStack gap="2">
            <LStack minW="0">
              <HStack w="full" justify="space-between" alignItems="end">
                {editing ? (
                  <TitleInput />
                ) : (
                  <Heading fontSize="heading.2" fontWeight="bold">
                    {node.name || "(untitled)"}
                  </Heading>
                )}
              </HStack>
            </LStack>
          </LStack>

          <LStack gap="2">
            <HStack w="full">
              {editing ? (
                <LibraryPageImportFromURL control={form.control} name="link" />
              ) : (
                node.link?.url && (
                  <LinkButton href={node.link?.url} size="xs" variant="subtle">
                    {node.link?.domain}
                  </LinkButton>
                )
              )}
            </HStack>
          </LStack>

          <ContentInput
            disabled={!editing}
            onAssetUpload={handleAssetUpload}
            initialValue={
              node.content ?? form.formState.defaultValues?.["content"]
            }
          />
        </LStack>

        {node && (node.children.length ?? 0) > 0 && (
          <LStack alignItems="start" w="full">
            <Heading>Child pages</Heading>
            <NodeCardRows
              libraryPath={libraryPath}
              context="library"
              nodes={node.children}
            />
          </LStack>
        )}
      </FormProvider>
    </styled.form>
  );
}
