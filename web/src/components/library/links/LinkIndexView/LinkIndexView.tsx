import { PaginationControls } from "src/components/site/PaginationControls/PaginationControls";
import { Unready } from "src/components/site/Unready";

import { Button } from "@/components/ui/button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { Input } from "@/components/ui/input";
import { VStack, styled } from "@/styled-system/jsx";

import { LinkCard } from "../LinkCard";

import { LinkResultList } from "./LinkResultList";
import { IndexingState, Props, useLinkIndexView } from "./useLinkIndexView";

export function LinkIndexView(props: Props) {
  const { form, data, handlers } = useLinkIndexView(props);

  if (form.formState.isLoading) return <Unready />;

  return (
    <VStack>
      <styled.form
        display="flex"
        w="full"
        onSubmit={handlers.handleSubmission}
        action="/l"
      >
        <Input
          w="full"
          borderRight="none"
          borderRightRadius="none"
          type="search"
          placeholder="Search or paste a new link"
          defaultValue={props.query}
          {...form.register("q")}
        />

        {(props.query || data.q) && (
          <Button
            borderX="none"
            borderRadius="none"
            type="reset"
            onClick={handlers.handleReset}
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
        >
          Search
        </Button>
      </styled.form>

      <PaginationControls
        path="/l"
        params={{ q: data.q }}
        onClick={handlers.handlePage}
        currentPage={props.page ?? 1}
        totalPages={data.links.total_pages}
        pageSize={data.links.page_size}
      />

      {data.indexing.state !== "not-indexing" ? (
        <IndexingStateBadge {...data.indexing} />
      ) : (
        <LinkResultList links={data.links} />
      )}
    </VStack>
  );
}

function IndexingStateBadge(props: IndexingState) {
  switch (props.state) {
    case "not-indexing":
      return <></>;
    case "indexing":
      return <>Indexing {props.url}...</>;
    case "indexed":
      return <LinkCard shape="row" link={props.link} />;
    case "error":
      return <>{props.error}</>;
  }
}
