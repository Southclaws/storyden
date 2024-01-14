"use client";

import { isEmpty } from "lodash";
import { FormProvider } from "react-hook-form";

import { ContentViewer } from "src/components/content/ContentViewer/ContentViewer";
import { FileDrop } from "src/components/content/FileDrop/FileDrop";
import { Breadcrumbs } from "src/components/directory/datagraph/Breadcrumbs";
import { ClusterList } from "src/components/directory/datagraph/ClusterList";
import { ItemGrid } from "src/components/directory/datagraph/ItemGrid";
import { EditAction } from "src/components/site/Action/Edit";
import { SaveAction } from "src/components/site/Action/Save";
import { Empty } from "src/components/site/Empty";
import { EmptyBox } from "src/components/site/EmptyBox";
import { Admonition } from "src/theme/components/Admonition";
import { Heading1 } from "src/theme/components/Heading/Index";
import { Input } from "src/theme/components/Input";

import { HStack, VStack, styled } from "@/styled-system/jsx";

import { ContentInput } from "./ContentInput";
import { TitleInput } from "./TitleInput";
import { Props, useClusterScreen } from "./useClusterScreen";

export function ClusterScreen(props: Props) {
  const {
    form,
    handlers: { handleSave, handleEditMode, handleAssetUpload },
    directoryPath,
    editing,
    cluster,
    isAllowedToEdit,
  } = useClusterScreen(props);

  return (
    <styled.form
      display="flex"
      flexDir="column"
      w="full"
      gap="3"
      alignItems="start"
      onSubmit={form.handleSubmit(handleSave)}
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
            defaultValue={cluster.slug}
            {...form.register("slug")}
          />
          {isAllowedToEdit && (
            <HStack>
              {editing ? (
                <SaveAction type="submit">Save</SaveAction>
              ) : (
                <EditAction onClick={handleEditMode}>Edit</EditAction>
              )}
            </HStack>
          )}
        </HStack>

        <VStack w="full" gap="2">
          <HStack w="full" h="full" maxH="64" justify="center" minW="0">
            {cluster.image_url ? (
              <styled.img
                maxHeight="64"
                borderRadius="lg"
                src={cluster.image_url}
              />
            ) : (
              <FileDrop>
                <EmptyBox>add images</EmptyBox>
              </FileDrop>
            )}
          </HStack>

          <VStack alignItems="start" w="full" minW="0">
            <HStack w="full" justify="space-between" alignItems="end">
              {editing ? (
                <TitleInput />
              ) : (
                <Heading1>{cluster.name || "(untitled)"}</Heading1>
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
                defaultValue={cluster.description}
                {...form.register("description")}
              />
            ) : (
              <styled.p>{cluster.description}</styled.p>
            )}
          </VStack>
        </VStack>

        {editing ? (
          <ContentInput onAssetUpload={handleAssetUpload} />
        ) : (
          <ContentViewer value={cluster.content ?? ""} />
        )}

        <VStack alignItems="start" w="full">
          {cluster.clusters.length === 0 && cluster.items.length === 0 && (
            <Empty>Nothing inside</Empty>
          )}

          {cluster && (cluster.clusters.length ?? 0) > 0 && (
            <ClusterList
              directoryPath={directoryPath}
              clusters={cluster.clusters}
            />
          )}

          {cluster && cluster.items.length > 0 && (
            <ItemGrid directoryPath={directoryPath} items={cluster.items} />
          )}
        </VStack>
      </FormProvider>
    </styled.form>
  );
}
