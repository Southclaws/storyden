import { keyBy } from "lodash";
import Link from "next/link";
import { ChangeEvent } from "react";

import {
  Identifier,
  NodeWithChildren,
  PropertySchemaList,
} from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";
import { useSortIndicator } from "@/components/site/SortIndicator";
import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { EmptyIcon } from "@/components/ui/icons/Empty";
import { MenuIcon } from "@/components/ui/icons/Menu";
import { Input } from "@/components/ui/input";
import {
  Box,
  Center,
  Grid,
  GridItem,
  HStack,
  LStack,
  WStack,
  styled,
} from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { linkDisabledProps } from "@/utils/anchor";
import { getAssetURL } from "@/utils/asset";

import { useLibraryPageContext } from "../../Context";
import { useEditState } from "../../useEditState";

import { AddPropertyMenu } from "./AddPropertyMenu/AddPropertyMenu";
import { ColumnMenu } from "./ColumnMenu";
import { useDirectoryBlockContext } from "./Context";
import { PropertyListMenu } from "./PropertyListMenu/PropertyListMenu";
import {
  MappableNodeField,
  mergeFieldsAndProperties,
  mergeFieldsAndPropertySchema,
} from "./column";
import { useDirectoryBlock } from "./useDirectoryBlock";

type Props = {
  nodes: NodeWithChildren[];
  block: ReturnType<typeof useDirectoryBlock>;
  currentChildPropertySchema: PropertySchemaList;
};

export function LibraryPageDirectoryBlockGrid({
  nodes,
  block,
  currentChildPropertySchema,
}: Props) {
  const { store } = useLibraryPageContext();
  const { sort, handleSort } = useSortIndicator();
  const { editing } = useEditState();
  const { handleSearch } = useDirectoryBlockContext();

  const { setChildPropertyValue } = store.getState();

  const fullSchema = mergeFieldsAndPropertySchema(
    currentChildPropertySchema,
    block,
  );

  const columns = fullSchema.filter(
    // We don't show fixed fields in the grid view. They are laid
    // out by the component itself and can be toggled on or off.
    (c) => !c._fixedFieldName,
  );

  const nameColumn = fullSchema.find((c) => c.fid === "fixed:name");
  const nameColumnHiddenState = nameColumn?.hidden ?? true;
  const linkColumn = fullSchema.find((c) => c.fid === "fixed:link");
  const linkColumnHiddenState = linkColumn?.hidden ?? true;
  const descColumn = fullSchema.find((c) => c.fid === "fixed:description");
  const descColumnHiddenState = descColumn?.hidden ?? true;

  const coverImageHiddenState =
    block.config?.columns.find((c) => c.fid === "fixed:primary_image")
      ?.hidden ?? false;

  const fullBleedCover =
    coverImageHiddenState === false &&
    nameColumnHiddenState === true &&
    linkColumnHiddenState === true &&
    descColumnHiddenState === true;

  function handleChildFieldValueChange(
    nodeID: Identifier,
    fid: MappableNodeField,
    value: string,
  ) {
    console.log("Child field value changed:", nodeID, fid, value);
    setChildPropertyValue(nodeID, fid, value);
  }

  function handleSearchChange(event: ChangeEvent<HTMLInputElement>) {
    handleSearch(event.target.value);
  }

  return (
    <LStack w="full" gap="2">
      <WStack justifyContent="end" bgColor="bg.subtle" borderRadius="sm" p="1">
        <Input
          variant="ghost"
          placeholder="Search..."
          size="xs"
          onChange={handleSearchChange}
        />

        {editing && (
          <>
            <AddPropertyMenu>
              <IconButton size="xs" variant="ghost" title="Add a new property.">
                <AddIcon />
              </IconButton>
            </AddPropertyMenu>

            <PropertyListMenu>
              <IconButton size="xs" variant="ghost">
                <MenuIcon />
              </IconButton>
            </PropertyListMenu>
          </>
        )}
      </WStack>

      <Grid
        w="full"
        gap="2"
        gridTemplateColumns="repeat(auto-fill, minmax(200px, 1fr))"
      >
        {nodes.map((node) => {
          const columnValues = mergeFieldsAndProperties(
            currentChildPropertySchema,
            node,
            block,
          );

          const valueTable = keyBy(columnValues, "fid");

          // When the cover image is displayed and the node has no primary image
          // and the name column is hidden, show the node name as a placeholder.
          const coverImagePlaceholder =
            node.primary_image === undefined &&
            !coverImageHiddenState &&
            nameColumnHiddenState;

          return (
            <GridItem
              key={node.id}
              position="relative"
              borderRadius="md"
              bg="bg.muted"
              display="flex"
              flexDirection="column"
              justifyContent="space-between"
              overflow="hidden"
              gap="2"
              // If no other fields are displayed, but the cover image is then
              // treat the card as a "full bleed" and remove the bottom padding.
              pb={fullBleedCover ? "0" : "1"}
            >
              {!coverImageHiddenState ? (
                <Link {...linkDisabledProps(editing)} href={`/l/${node.slug}`}>
                  {node.primary_image ? (
                    <styled.img
                      src={getAssetURL(node.primary_image.path)}
                      objectFit="cover"
                      w="full"
                      h="32"
                    />
                  ) : (
                    <Center h="32" p="2">
                      {coverImagePlaceholder ? (
                        <styled.span
                          textAlign="center"
                          textWrap="balance"
                          color="fg.muted"
                          lineClamp={3}
                        >
                          {node.name}
                        </styled.span>
                      ) : (
                        <EmptyIcon />
                      )}
                    </Center>
                  )}
                </Link>
              ) : (
                <Box />
              )}

              {fullBleedCover ? null : (
                <LStack gap="0" px="2">
                  {!nameColumnHiddenState &&
                    (editing ? (
                      <styled.input
                        w="full"
                        fontWeight="semibold"
                        defaultValue={node.name}
                        onChange={(event) => {
                          console.log(event);
                          handleChildFieldValueChange(
                            node.id,
                            "fixed:name" as MappableNodeField,
                            event.target.value,
                          );
                        }}
                        _focusVisible={{
                          outline: "none",
                        }}
                      />
                    ) : (
                      <Link href={`/l/${node.slug}`}>
                        <styled.div
                          fontWeight="semibold"
                          lineClamp="1"
                          textOverflow="ellipsis"
                          wordBreak="break-all"
                        >
                          {node.name}
                        </styled.div>
                      </Link>
                    ))}

                  {!descColumnHiddenState &&
                    (editing ? (
                      <styled.input
                        w="full"
                        placeholder="Description..."
                        _placeholder={{
                          color: "fg.subtle",
                        }}
                        defaultValue={node.description}
                        onChange={(event) =>
                          handleChildFieldValueChange(
                            node.id,
                            "fixed:description" as MappableNodeField,
                            event.target.value,
                          )
                        }
                        _focusVisible={{
                          outline: "none",
                        }}
                      />
                    ) : (
                      node.description && (
                        <styled.div
                          fontSize="sm"
                          color="fg.muted"
                          lineClamp="1"
                          textOverflow="ellipsis"
                          wordBreak="break-all"
                        >
                          {node.description}
                        </styled.div>
                      )
                    ))}

                  {!linkColumnHiddenState &&
                    (editing ? (
                      <styled.input
                        w="full"
                        placeholder="Link..."
                        _placeholder={{
                          color: "fg.subtle",
                        }}
                        defaultValue={node.link?.url}
                        onChange={(event) =>
                          handleChildFieldValueChange(
                            node.id,
                            "fixed:link" as MappableNodeField,
                            event.target.value,
                          )
                        }
                        _focusVisible={{
                          outline: "none",
                        }}
                      />
                    ) : (
                      node.link && (
                        <styled.a
                          fontSize="sm"
                          color="fg.muted"
                          lineClamp="1"
                          textOverflow="ellipsis"
                          wordBreak="break-all"
                          href={node.link.url}
                        >
                          {node.link.title || node.link.url}
                        </styled.a>
                      )
                    ))}
                </LStack>
              )}

              {columns.length > 0 && (
                <styled.dl className={lstack()} gap="0" px="2" pb="2">
                  {columns.map((property) => {
                    const column = valueTable[property.fid];
                    if (!column) {
                      console.warn(
                        `unable to find property ${property.fid} in value table`,
                        valueTable,
                      );
                      return null;
                    }

                    const isSortingBy = sort?.property === property.name;
                    const isSortingAsc = sort?.order === "asc";
                    const isSortingDesc = sort?.order === "desc";
                    const isSorting =
                      isSortingBy && (isSortingAsc || isSortingDesc);
                    const sortState = isSorting
                      ? isSortingAsc
                        ? "asc"
                        : "desc"
                      : "none";

                    function handleClickSortAction() {
                      handleSort(property.name);
                    }

                    const handleCellChange = (
                      v: ChangeEvent<HTMLInputElement>,
                    ) => {
                      handleChildFieldValueChange(
                        node.id,
                        column.fid as MappableNodeField,
                        v.target.value,
                      );
                    };

                    return (
                      <HStack
                        w="full"
                        key={property.fid}
                        cursor="pointer"
                        {...(isSorting && {
                          "data-active": "",
                        })}
                        p="0"
                      >
                        <HStack w="full" pr="1">
                          <ColumnMenu column={property}>
                            <styled.dt
                              display="inline-flex"
                              minW="0"
                              flexGrow="1"
                              alignItems="center"
                              justifyContent="space-between"
                              gap="1"
                              textWrap="nowrap"
                              flexWrap="nowrap"
                              color="fg.subtle"
                            >
                              {property.name}
                            </styled.dt>
                          </ColumnMenu>

                          <styled.dd>
                            {editing ? (
                              <styled.input
                                w="full"
                                minW="0"
                                textAlign="right"
                                defaultValue={column.value}
                                onChange={handleCellChange}
                                _focusVisible={{
                                  outline: "none",
                                }}
                              />
                            ) : column.href ? (
                              <Link href={column.href}>{column.value}</Link>
                            ) : (
                              <>{column.value}</>
                            )}
                          </styled.dd>
                        </HStack>
                      </HStack>
                    );
                  })}
                </styled.dl>
              )}
            </GridItem>
          );
        })}
      </Grid>

      {nodes.length === 0 && (
        <Center w="full">
          <EmptyState hideContributionLabel>
            There are no pages here.
          </EmptyState>
        </Center>
      )}
    </LStack>
  );
}
