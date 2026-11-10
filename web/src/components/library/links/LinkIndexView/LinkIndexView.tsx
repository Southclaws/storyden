import { ReactNode, useRef } from "react";

import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { Unready } from "@/components/site/Unready";
import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { Input } from "@/components/ui/input";
import { LStack, styled } from "@/styled-system/jsx";

import { LinkResultList } from "./LinkResultList";
import { Props, useLinkIndexView } from "./useLinkIndexView";

export function LinkIndexView(props: Props) {
  const { data, handlers } = useLinkIndexView(props);
  const inputRef = useRef<HTMLInputElement>(null);

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
        <LinkResultList links={data.links} />
        <PaginationControls
          path="/links"
          params={{ q: data.query }}
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
      <styled.form display="flex" w="full" onSubmit={handlers.handleSubmission}>
        <Input
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
