import { XMarkIcon } from "@heroicons/react/24/outline";

import { PaginationControls } from "src/components/site/PaginationControls/PaginationControls";
import { Unready } from "src/components/site/Unready";

import { MemberList } from "../MemberList";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { VStack, styled } from "@/styled-system/jsx";

import { Props, useMemberIndexView } from "./useMemberIndexView";

export function MemberIndexView(props: Props) {
  const { form, data, handlers } = useMemberIndexView(props);

  if (form.formState.isLoading) return <Unready />;

  return (
    <VStack>
      <styled.form
        display="flex"
        w="full"
        onSubmit={handlers.handleSubmission}
        action="/p"
      >
        <Input
          w="full"
          borderRight="none"
          borderRightRadius="none"
          type="search"
          placeholder="Search for members"
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
        path="/p"
        params={{ q: data.q }}
        onClick={handlers.handlePage}
        currentPage={props.page ?? 1}
        totalPages={data.results.total_pages}
        pageSize={data.results.page_size}
      />

      <MemberList
        onChange={handlers.handleMutate}
        profiles={data.results.profiles}
      />
    </VStack>
  );
}
