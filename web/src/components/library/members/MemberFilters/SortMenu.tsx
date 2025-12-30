"use client";

import { parseAsString, useQueryState } from "nuqs";

import { Button } from "@/components/ui/button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { SortIcon } from "@/components/ui/icons/Sort";
import * as Menu from "@/components/ui/menu";
import { HStack } from "@/styled-system/jsx";

const SORT_OPTIONS = [
  { value: "name", label: "Display name (A-Z)" },
  { value: "-name", label: "Display name (Z-A)" },
  { value: "handle", label: "Handle (A-Z)" },
  { value: "-handle", label: "Handle (Z-A)" },
  { value: "created_at", label: "Join date (oldest)" },
  { value: "-created_at", label: "Join date (newest)" },
] as const;

export function SortMenu() {
  const [sort, setSort] = useQueryState("sort", parseAsString);

  const handleSortChange = async (value: string) => {
    await setSort(value);
  };

  const currentLabel =
    SORT_OPTIONS.find((opt) => opt.value === sort)?.label || "Sort by...";

  return (
    <Menu.Root positioning={{ placement: "bottom-start" }} lazyMount>
      <Menu.Trigger asChild>
        <Button variant="subtle" size="sm" aria-label="Sort options">
          <HStack gap="1">
            <SortIcon />
            {currentLabel}
          </HStack>
        </Button>
      </Menu.Trigger>

      <Menu.Positioner>
        <Menu.Content minW="56">
          <Menu.ItemGroup id="sort-options">
            {SORT_OPTIONS.map((option) => (
              <Menu.Item
                key={option.value}
                value={option.value}
                onClick={() => handleSortChange(option.value)}
                aria-label={option.label}
              >
                <HStack gap="2" justify="space-between" w="full">
                  <span>{option.label}</span>
                  {sort === option.value && <CheckIcon />}
                </HStack>
              </Menu.Item>
            ))}
          </Menu.ItemGroup>
        </Menu.Content>
      </Menu.Positioner>
    </Menu.Root>
  );
}
