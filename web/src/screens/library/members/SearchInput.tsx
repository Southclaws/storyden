"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { type SubmitEvent, useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { Input } from "@/components/ui/input";
import { styled } from "@/styled-system/jsx";

type Props = {
  index: string;
  initialQuery: string | undefined;
  placeholder?: string;
};

export function SearchInput({
  index,
  initialQuery,
  placeholder = "Search for members",
}: Props) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [query, setQuery] = useState(initialQuery ?? "");

  useEffect(() => {
    setQuery(initialQuery ?? "");
  }, [initialQuery]);

  function pushSearch(nextQuery: string) {
    const params = new URLSearchParams(searchParams.toString());
    const trimmedQuery = nextQuery.trim();

    params.delete("page");

    if (trimmedQuery) {
      params.set("q", trimmedQuery);
    } else {
      params.delete("q");
    }

    const queryString = params.toString();
    router.push(queryString ? `${index}?${queryString}` : index);
  }

  function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    pushSearch(query);
  }

  function handleReset() {
    setQuery("");
    pushSearch("");
  }

  return (
    <styled.form display="flex" w="full" onSubmit={handleSubmit} action={index}>
      <Input
        w="full"
        borderRight="none"
        borderRightRadius="none"
        type="search"
        placeholder={placeholder}
        value={query}
        onChange={(event) => setQuery(event.target.value)}
      />

      {query && (
        <Button
          borderX="none"
          borderRadius="none"
          type="reset"
          onClick={handleReset}
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
  );
}
