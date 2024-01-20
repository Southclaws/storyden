"use client";

import { isEmpty } from "lodash";
import { FormProvider } from "react-hook-form";

import { ContentViewer } from "src/components/content/ContentViewer/ContentViewer";
import { Breadcrumbs } from "src/components/directory/datagraph/Breadcrumbs";
import { ClusterList } from "src/components/directory/datagraph/ClusterList";
import { ItemGrid } from "src/components/directory/datagraph/ItemGrid";
import { CancelAction } from "src/components/site/Action/Cancel";
import { EditAction } from "src/components/site/Action/Edit";
import { SaveAction } from "src/components/site/Action/Save";
import { Empty } from "src/components/site/Empty";
import { Admonition } from "src/theme/components/Admonition";
import { Heading1, Heading2 } from "src/theme/components/Heading/Index";
import { Input } from "src/theme/components/Input";
import {
  Popover,
  PopoverAnchor,
  PopoverContent,
  PopoverPositioner,
} from "src/theme/components/Popover";

import { HStack, VStack, styled } from "@/styled-system/jsx";

import { AssetsInput } from "./AssetsInput";
import { ContentInput } from "./ContentInput";
import { TitleInput } from "./TitleInput";
import { Props, useDatagraphNodeScreen } from "./useDatagraphNodeScreen";

export function DatagraphNodeScreen(props: Props) {
  const {
    form,
    handlers: {
      handleSubmit,
      handleEditMode,
      handleAssetUpload,
      handleAssetRemove,
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
        <HStack w="full" justify="space-between">
          <Breadcrumbs
            directoryPath={directoryPath}
            create={editing ? "edit" : "show"}
            defaultValue={node.slug}
            {...form.register("slug")}
          />
          {isAllowedToEdit && (
            <Popover open={isSaving} lazyMount>
              <PopoverAnchor>
                <HStack>
                  {editing ? (
                    <>
                      <CancelAction onClick={handleEditMode}>
                        Cancel
                      </CancelAction>
                      <SaveAction type="submit">Save</SaveAction>
                    </>
                  ) : (
                    <EditAction onClick={handleEditMode}>Edit</EditAction>
                  )}
                </HStack>
              </PopoverAnchor>

              <PopoverPositioner>
                <PopoverContent p="2">Saved!</PopoverContent>
              </PopoverPositioner>
            </Popover>
          )}
        </HStack>

        <VStack w="full" alignItems="start" gap="2">
          <AssetsInput
            editing={editing}
            initialAssets={node.assets}
            handleAssetUpload={handleAssetUpload}
            handleAssetRemove={handleAssetRemove}
          />

          <VStack alignItems="start" w="full" minW="0">
            <HStack w="full" justify="space-between" alignItems="end">
              {editing ? (
                <TitleInput />
              ) : (
                <Heading1>{node.name || "(untitled)"}</Heading1>
              )}
            </HStack>

            {/* TODO: Links and link editing for clusters
            {cluster.link && (
              <Box w="full">
                <Link href={cluster.link?.url} w="full" size="sm">
                  {cluster.link?.url}
                </Link>
              </Box>
            )} */}

            {editing ? (
              <Input
                placeholder="Description"
                defaultValue={node.description}
                {...form.register("description")}
              />
            ) : (
              <styled.p>{node.description}</styled.p>
            )}
          </VStack>
        </VStack>

        {editing ? (
          <ContentInput onAssetUpload={handleAssetUpload} />
        ) : (
          <ContentViewer value={node.content ?? ""} />
        )}

        {node.type === "cluster" && (
          <VStack alignItems="start" w="full">
            {node.clusters.length === 0 && node.items.length === 0 && (
              <Empty>Nothing inside</Empty>
            )}

            {node && (node.clusters.length ?? 0) > 0 && (
              <ClusterList
                directoryPath={directoryPath}
                clusters={node.clusters}
              />
            )}

            {node && node.items.length > 0 && (
              <ItemGrid directoryPath={directoryPath} items={node.items} />
            )}
          </VStack>
        )}

        {node.type === "item" && (
          <VStack alignItems="start" w="full">
            <Heading2>Member of</Heading2>

            {node.clusters.length ? (
              <ClusterList
                directoryPath={directoryPath}
                clusters={node.clusters}
              />
            ) : (
              <Empty>No Items</Empty>
            )}
          </VStack>
        )}
      </FormProvider>
    </styled.form>
  );
}
