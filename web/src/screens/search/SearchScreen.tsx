"use client";

import { DatagraphSearchResults } from "src/components/search/DatagraphSearchResults";
import { UnreadyBanner } from "src/components/site/Unready";

import { DatagraphItemKind } from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { MultiSelectPicker } from "@/components/ui/MultiSelectPicker";
import { Button } from "@/components/ui/button";
import { DatagraphKindFilterField } from "@/components/ui/form/DatagraphKindFilterField";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { ReplyIcon } from "@/components/ui/icons/Reply";
import { SearchIcon } from "@/components/ui/icons/Search";
import { Input } from "@/components/ui/input";
import { useI18n } from "@/i18n/provider";
import { Flex, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";

import { Props, useSearchScreen } from "./useSearch";

export function SearchScreen(props: Props) {
  const { ready, form, error, isLoading, data, handlers, filters } =
    useSearchScreen(props);
  const { t } = useI18n();

  const { query, page, results } = data;

  return (
    <styled.form
      className={vstack()}
      display="flex"
      w="full"
      onSubmit={handlers.handleSearch}
      action="/search"
    >
      <WStack gap="0">
        <Input
          w="full"
          size="md"
          borderRight="none"
          borderRightRadius="none"
          type="search"
          background="bg.default"
          placeholder={t("Search...")}
          _focus={{
            // NOTE: This disables the default focus behaviour styles for inputs.
            boxShadow: "none" as any, // TODO: Fix types at Park-UI or Panda level
            borderColor: "border.default",
          }}
          {...form.register("q")}
        />

        {query && (
          <Button
            size="md"
            variant="outline"
            borderX="none"
            borderRadius="none"
            borderColor="border.default"
            type="reset"
            onClick={handlers.handleReset}
          >
            <CancelIcon />
          </Button>
        )}
        <Button
          size="md"
          variant="outline"
          flexShrink="0"
          borderLeft="none"
          borderLeftRadius="none"
          borderColor="border.default"
          type="submit"
          width="min"
          loading={isLoading}
        >
          <SearchIcon />
        </Button>
      </WStack>

      <LStack w="full" gap="2">
        <DatagraphKindFilterField
          control={form.control}
          name="kind"
          items={[
            {
              label: t("Threads"),
              description: t("Include discussion threads in the search."),
              icon: <DiscussionIcon />,
              value: DatagraphItemKind.thread,
            },
            {
              label: t("Replies"),
              description: t("Include replies to discussion threads in the search."),
              icon: <ReplyIcon />,
              value: DatagraphItemKind.reply,
            },
            {
              label: t("Library"),
              description: t("Include library pages in the search."),
              icon: <LibraryIcon />,
              value: DatagraphItemKind.node,
            },
          ]}
        />

        <Flex
          w="full"
          gap="2"
          flexDirection={{
            base: "column",
            md: "row",
          }}
        >
          <MultiSelectPicker
            value={filters.authorsValue}
            onChange={handlers.handleAuthorsChange}
            onQuery={handlers.handleQueryAuthors}
            queryResults={filters.authorsResults}
            queryError={filters.authorsError}
            inputPlaceholder={t("Authors...")}
            size="sm"
            triggerProps={{
              width: "full",
              minW: "32",
              flexShrink: "1",
            }}
          />

          {filters.showCategories && (
            <MultiSelectPicker
              value={filters.categoriesValue}
              onChange={handlers.handleCategoriesChange}
              onQuery={handlers.handleQueryCategories}
              queryResults={filters.categoriesResults}
              queryError={filters.categoriesError}
              inputPlaceholder={t("Categories...")}
              size="sm"
              triggerProps={{
                width: "full",
                minW: "32",
                flexShrink: "1",
              }}
            />
          )}

          {filters.showTags && (
            <MultiSelectPicker
              value={filters.tagsValue}
              onChange={handlers.handleTagsChange}
              onQuery={handlers.handleQueryTags}
              queryResults={filters.tagsResults}
              queryError={filters.tagsError}
              inputPlaceholder={t("Tags...")}
              size="sm"
              triggerProps={{
                width: "full",
                minW: "32",
                flexShrink: "1",
              }}
            />
          )}
        </Flex>
      </LStack>

      {isLoading || error !== undefined ? (
        <UnreadyBanner error={error} />
      ) : results?.items.length ? (
        <>
          <PaginationControls
            path="/search"
            params={{ q: query }}
            currentPage={page}
            totalPages={results.total_pages}
            pageSize={results.page_size}
          />
          <DatagraphSearchResults result={results} />
        </>
      ) : (
        <EmptyState hideContributionLabel>
          {query
            ? results && page > results?.total_pages
              ? t("You've gone past the last page! Nothing to see here.")
              : t("No search results.")
            : t("Go forth, seek far and wide.")}
        </EmptyState>
      )}
    </styled.form>
  );
}
