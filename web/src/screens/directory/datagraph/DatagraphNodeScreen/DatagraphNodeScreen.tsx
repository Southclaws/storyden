"use client";

import { isEmpty } from "lodash";
import { FormProvider } from "react-hook-form";

import { Breadcrumbs } from "src/components/directory/datagraph/Breadcrumbs";
import { DatagraphNodeMenu } from "src/components/directory/datagraph/DatagraphNodeMenu/DatagraphNodeMenu";
import { NodeCardRows } from "src/components/directory/datagraph/NodeCardList";
import { CancelAction } from "src/components/site/Action/Cancel";
import { EditAction } from "src/components/site/Action/Edit";
import { SaveAction } from "src/components/site/Action/Save";

import { Admonition } from "@/components/ui/admonition";
import { Heading } from "@/components/ui/heading";
import * as Popover from "@/components/ui/popover";
import { HStack, LStack, VStack, styled } from "@/styled-system/jsx";

import { ContentInput } from "./ContentInput";
import { TitleInput } from "./TitleInput";
import { Props, useDatagraphNodeScreen } from "./useDatagraphNodeScreen";

export function DatagraphNodeScreen(props: Props) {
  const {
    form,
    handlers: {
      handleSubmit,
      handleEditMode,
      handleVisibilityChange,
      handleDelete,
      handleAssetUpload,
    },
    directoryPath,
    editing,
    node,
    isAllowedToEdit,
    isSaving,
  } = useDatagraphNodeScreen(props);

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
      <Admonition
        value={!isEmpty(form.formState.errors)}
        title="Errors"
        kind="failure"
      >
        {Object.values(form.formState.errors).map((error, i) => (
          <p key={i}>{error.message}</p>
        ))}
      </Admonition>

      <FormProvider {...form}>
        <LStack h="full">
          <HStack w="full" justify="space-between">
            <Breadcrumbs
              directoryPath={directoryPath}
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
                        <CancelAction onClick={handleEditMode}>
                          Cancel
                        </CancelAction>
                        <SaveAction type="submit">Save</SaveAction>
                      </>
                    ) : (
                      <>
                        <EditAction onClick={handleEditMode}>Edit</EditAction>
                      </>
                    )}
                    <DatagraphNodeMenu
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

          <VStack w="full" alignItems="start" gap="2">
            <VStack alignItems="start" w="full" minW="0">
              <HStack w="full" justify="space-between" alignItems="end">
                {editing ? (
                  <TitleInput />
                ) : (
                  <Heading fontSize="heading.2" fontWeight="bold">
                    {node.name || "(untitled)"}
                  </Heading>
                )}
              </HStack>
            </VStack>
          </VStack>

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
              directoryPath={directoryPath}
              context="directory"
              nodes={node.children}
            />
          </LStack>
        )}
      </FormProvider>
    </styled.form>
  );
}
