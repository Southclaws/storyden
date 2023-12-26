import { XMarkIcon } from "@heroicons/react/24/outline";

import { LinkRefList } from "src/components/directory/links/LinkRefList";
import { PaginationControls } from "src/components/site/PaginationControls/PaginationControls";
import { Unready } from "src/components/site/Unready";
import { Button } from "src/theme/components/Button";
import { Input } from "src/theme/components/Input";

import { VStack, styled } from "@/styled-system/jsx";

import { Props, useLinkIndexView } from "./useLinkIndexView";

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
          placeholder="Search for addresses, keywords or domains"
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
        path="/l"
        params={{ q: data.q }}
        onClick={handlers.handlePage}
        currentPage={props.page ?? 1}
        totalPages={data.links.total_pages}
        pageSize={data.links.page_size}
      />

      <LinkRefList result={data.links} />
    </VStack>
  );
}
