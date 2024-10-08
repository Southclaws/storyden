import Link from "next/link";

import { styled } from "@/styled-system/jsx";

export function LibraryBadge() {
  return (
    <styled.span
      backgroundColor="bg.accent"
      color="fg.accent"
      px="1"
      borderRadius="md"
    >
      <Link href="/l">library</Link>
    </styled.span>
  );
}

// TODO: Make this a recipe component.
export function NewBadge() {
  return (
    <styled.span
      fontSize="xs"
      fontWeight="bold"
      backgroundColor="bg.accent"
      color="fg.accent"
      px="1"
      py="0.5"
      borderRadius="sm"
    >
      New
    </styled.span>
  );
}
