"use client";

import Link from "next/link";
import { FormProvider } from "react-hook-form";

import { CancelAction } from "src/components/site/Action/Cancel";
import { EditAction } from "src/components/site/Action/Edit";
import { SaveAction } from "src/components/site/Action/Save";

import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { LibraryPageImportFromURL } from "@/components/library/LibraryPageImportFromURL/LibraryPageImportFromURL";
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
import { NodeCardRows } from "@/components/library/NodeCardList";
import { Anchor } from "@/components/site/Anchor";
import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import * as Popover from "@/components/ui/popover";
import { HStack, LStack, styled } from "@/styled-system/jsx";

import { ContentInput } from "./ContentInput";
import { TitleInput } from "./TitleInput";
import { Props, useLibraryPageScreen } from "./useLibraryPageScreen";

export function LibraryPageScreen(props: Props) {
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
    isAllowedToEdit,
    isSaving,
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
              <Popover.Root open={isSaving} lazyMount>
                <Popover.Anchor>
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
                </Popover.Anchor>

                <Popover.Positioner>
                  <Popover.Content p="2">Saved!</Popover.Content>
                </Popover.Positioner>
              </Popover.Root>
            )}
          </HStack>

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
