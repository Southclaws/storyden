"use client";

import { MagnifyingGlassIcon } from "@heroicons/react/24/outline";
import { XMarkIcon } from "@heroicons/react/24/solid";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { styled } from "@/styled-system/jsx";

import { Props, useSearch } from "./useSearch";

export function Search(props: Props) {
  const { form, data, handlers } = useSearch(props);
  return (
    <styled.form
      display="flex"
      w="full"
      onSubmit={handlers.handleSearch}
      action="/search"
    >
      <Input
        w="full"
        size="sm"
        borderRight="none"
        borderRightRadius="none"
        type="search"
        defaultValue={props.query}
        background="bg.default"
        placeholder={`Search...`}
        _focus={{
          // NOTE: This disables the default focus behaviour styles for inputs.
          boxShadow: "none" as any, // TODO: Fix types at Park-UI or Panda level
          borderColor: "border.default",
        }}
        {...form.register("q")}
      />

      {(props.query || data.q) && (
        <Button
          size="sm"
          variant="outline"
          borderX="none"
          borderRadius="none"
          borderColor="border.default"
          type="reset"
          onClick={handlers.handleReset}
        >
          <XMarkIcon />
        </Button>
      )}
      <Button
        size="sm"
        variant="outline"
        flexShrink="0"
        borderLeft="none"
        borderLeftRadius="none"
        borderColor="border.default"
        type="submit"
        width="min"
      >
        <MagnifyingGlassIcon />
      </Button>
    </styled.form>
  );
}
