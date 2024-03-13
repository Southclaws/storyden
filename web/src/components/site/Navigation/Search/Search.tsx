import { MagnifyingGlassIcon } from "@heroicons/react/24/outline";
import { XMarkIcon } from "@heroicons/react/24/solid";

import { Button } from "src/theme/components/Button";
import { Input } from "src/theme/components/Input";

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
        borderRight="none"
        borderRightRadius="none"
        type="search"
        defaultValue={props.query}
        background="bg.default"
        placeholder={`Search...`}
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
        <MagnifyingGlassIcon />
      </Button>
    </styled.form>
  );
}
