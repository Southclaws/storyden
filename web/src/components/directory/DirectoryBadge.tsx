import Link from "next/link";

import { styled } from "@/styled-system/jsx";

export function DirectoryBadge() {
  return (
    <styled.span backgroundColor="accent.100" px="1" borderRadius="md">
      <Link href="/directory">directory</Link>
    </styled.span>
  );
}

// TODO: Make this a recipe component.
export function NewBadge() {
  return (
    <styled.span
      fontSize="xs"
      fontWeight="bold"
      backgroundColor="accent.100"
      color="accent.800"
      px="1"
      py="0.5"
      borderRadius="sm"
    >
      New
    </styled.span>
  );
}
