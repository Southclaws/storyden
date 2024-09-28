import { XMarkIcon } from "@heroicons/react/24/solid";

import { TextPostList } from "src/components/feed/text/TextPostList";
import { PaginationControls } from "src/components/site/PaginationControls/PaginationControls";
import { Unready } from "src/components/site/Unready";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { VStack, styled } from "@/styled-system/jsx";

import { Props, useThreadIndexView } from "./useThreadIndexView";

export function ThreadIndexView(props: Props) {
  const { form, data, handlers } = useThreadIndexView(props);

  if (form.formState.isLoading) return <Unready />;

  return (
    <VStack>
      <styled.form
        display="flex"
        w="full"
        onSubmit={handlers.handleSubmission}
        action="/t"
      >
        <Input
          w="full"
          borderRight="none"
          borderRightRadius="none"
          type="search"
          placeholder="Search discussions"
          defaultValue={props.query}
          {...form.register("q")}
        />

        {props.query && (
          <Button
            borderX="none"
            borderRadius="none"
            type="reset"
            onClick={handlers.handleReset}
          >
            <XMarkIcon />
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
        path="/t"
        params={{ q: props.query ?? "" }}
        onClick={handlers.handlePage}
        currentPage={props.page ?? 1}
        totalPages={data.threads.total_pages}
        pageSize={data.threads.page_size}
      />

      <TextPostList threads={data.threads.threads} />
    </VStack>
  );
}
