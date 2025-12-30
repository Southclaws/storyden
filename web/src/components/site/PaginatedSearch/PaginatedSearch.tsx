"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter, useSearchParams } from "next/navigation";
import { parseAsInteger, useQueryState } from "nuqs";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { PaginationControls } from "src/components/site/PaginationControls/PaginationControls";

import { Button } from "@/components/ui/button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { Input } from "@/components/ui/input";
import { VStack, styled } from "@/styled-system/jsx";

export type Props = {
  index: string;
  initialQuery: string | undefined;
  initialPage: number | undefined;
  totalPages: number;
  pageSize: number;
};

export const FormSchema = z.object({
  q: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function usePaginatedSearch({ initialQuery, initialPage, index }: Props) {
  const router = useRouter();
  const searchParams = useSearchParams();

  const [page, setPage] = useQueryState("page", {
    ...parseAsInteger,
    defaultValue: initialPage ?? 1,
  });

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      q: initialQuery,
    },
  });

  const { q } = form.watch();

  // Get all current params as an object
  const currentParams = Object.fromEntries(searchParams.entries());

  const handleSubmit = form.handleSubmit(async (data) => {
    const { page, ...paramsWithoutPage } = currentParams;
    const params = new URLSearchParams({
      ...paramsWithoutPage,
      q: data.q,
      page: "1",
    });
    router.push(`${index}?${params.toString()}`);
  });

  const handleReset = async () => {
    form.reset();
    // Preserve all params except 'q' and 'page'
    const { q, page, ...restParams } = currentParams;
    const params = new URLSearchParams(restParams);
    const url = params.toString() ? `${index}?${params.toString()}` : index;
    router.push(url);
  };

  const handlePage = async (nextPage: number) => {
    setPage(nextPage);
  };

  return {
    form,
    query: q,
    page,
    currentParams,
    handlers: {
      handleSubmit,
      handleReset,
      handlePage,
    },
  };
}

export function PaginatedSearch(props: Props) {
  const { form, query, page, currentParams, handlers } =
    usePaginatedSearch(props);

  return (
    <VStack w="full">
      <styled.form
        display="flex"
        w="full"
        onSubmit={handlers.handleSubmit}
        action="/m"
      >
        <Input
          w="full"
          borderRight="none"
          borderRightRadius="none"
          type="search"
          placeholder="Search for members"
          defaultValue={props.initialQuery}
          {...form.register("q")}
        />

        {query && (
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
        path={props.index}
        params={currentParams}
        currentPage={page}
        totalPages={props.totalPages}
        pageSize={props.pageSize}
      />
    </VStack>
  );
}
