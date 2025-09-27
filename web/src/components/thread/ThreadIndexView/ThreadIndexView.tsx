import { PaginationControls } from "src/components/site/PaginationControls/PaginationControls";
import { Unready } from "src/components/site/Unready";

import { ThreadReferenceList } from "@/components/post/ThreadReferenceList";
import { Button } from "@/components/ui/button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { Input } from "@/components/ui/input";
import { LStack, styled } from "@/styled-system/jsx";

import { Props, useThreadIndexView } from "./useThreadIndexView";

export function ThreadIndexView(props: Props) {
  const { form, data, handlers } = useThreadIndexView(props);

  if (form.formState.isLoading) return <Unready />;

  return (
    <LStack>
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
        path="/t"
        params={{ q: props.query ?? "" }}
        onClick={handlers.handlePage}
        currentPage={props.page ?? 1}
        totalPages={data.threads.total_pages}
        pageSize={data.threads.page_size}
      />

      <ThreadReferenceList threads={data.threads.threads} />
    </LStack>
  );
}
