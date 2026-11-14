import { parseAsStringLiteral, useQueryState } from "nuqs";
import { ReactNode, useRef } from "react";

import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { Unready } from "@/components/site/Unready";
import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { LayoutGridIcon } from "@/components/ui/icons/LayoutGrid";
import { LayoutListIcon } from "@/components/ui/icons/LayoutList";
import { Input } from "@/components/ui/input";
import * as ToggleGroup from "@/components/ui/toggle-group";
import { HStack, LStack, styled } from "@/styled-system/jsx";

import { LinkResultList } from "./LinkResultList";
import { Props, useLinkIndexView } from "./useLinkIndexView";

export function LinkIndexView(props: Props) {
  const { data, handlers } = useLinkIndexView(props);
  const inputRef = useRef<HTMLInputElement>(null);
  const [layout, setLayout] = useQueryState(
    "view",
    parseAsStringLiteral(["rows", "grid"])
      .withDefault("rows")
      .withOptions({ clearOnDefault: true }),
  );
  const paginationParams: Record<string, string> = { q: data.query };

  if (layout === "grid") {
    paginationParams["view"] = "grid";
  }

  let submitLabel = "Search";

  if (data.isURL && data.canUseLinks) {
    submitLabel = "Add link";
  }

  if (data.isURL && data.canUseLinks && data.isSaved) {
    submitLabel = "Saved";
  }

  let results: ReactNode;

  if (data.isLoading) {
    results = <Unready />;
  } else if (data.error) {
    results = (
      <>
        <Unready error="Could not search links." />
        <Button variant="outline" onClick={handlers.handleRefresh}>
          Try again
        </Button>
      </>
    );
  } else if (data.links) {
    results = (
      <>
        <LinkResultList links={data.links} layout={layout} />
        <PaginationControls
          path="/links"
          params={paginationParams}
          onClick={handlers.handlePage}
          currentPage={data.page}
          totalPages={data.links.total_pages}
          pageSize={data.links.page_size}
        />
      </>
    );
  }

  return (
    <LStack w="full" gap="4">
      <HStack w="full" gap="2" flexWrap="wrap">
        <styled.form
          display="flex"
          flexGrow="1"
          flexShrink="1"
          flexBasis="64"
          minW="0"
          onSubmit={handlers.handleSubmission}
        >
          <Input
            size="md"
            w="full"
            minW="0"
            borderRight="none"
            borderRightRadius="none"
            type="search"
            aria-label="Search links or paste a URL"
            placeholder="Search or paste a new link"
            ref={inputRef}
            value={data.query}
            onChange={(event) => handlers.handleQuery(event.target.value)}
          />

          {data.query && (
            <Button
              size="md"
              borderX="none"
              borderRadius="none"
              type="button"
              aria-label="Clear search"
              onClick={() => {
                handlers.handleReset();
                inputRef.current?.focus();
              }}
            >
              <CancelIcon />
            </Button>
          )}
          <Button
            size="md"
            flexShrink="0"
            borderLeft="none"
            borderLeftRadius="none"
            type="submit"
            width="min"
            disabled={data.canUseLinks && data.isURL && data.isSaved}
            loading={data.adding.state === "adding"}
            loadingText="Adding…"
          >
            {submitLabel}
          </Button>
        </styled.form>

        <ToggleGroup.Root
          aria-label="Link layout"
          value={[layout]}
          size="sm"
          flexShrink="0"
          marginLeft="auto"
          variant="outline"
          onValueChange={({ value }) => {
            const next = value[0];

            if (next === "rows" || next === "grid") {
              void setLayout(next);
            }
          }}
        >
          <ToggleGroup.Item value="rows" aria-label="List view">
            <LayoutListIcon />
          </ToggleGroup.Item>
          <ToggleGroup.Item value="grid" aria-label="Grid view">
            <LayoutGridIcon />
          </ToggleGroup.Item>
        </ToggleGroup.Root>
      </HStack>

      {data.adding.state === "error" && (
        <Admonition
          value
          kind="failure"
          title="Could not add link"
          role="alert"
          onChange={handlers.handleDismissError}
        >
          {data.adding.error} Try adding the link again.
        </Admonition>
      )}

      <LStack
        as="section"
        aria-label="Link results"
        aria-busy={data.isLoading}
        gap="4"
      >
        {results}
      </LStack>
    </LStack>
  );
}
