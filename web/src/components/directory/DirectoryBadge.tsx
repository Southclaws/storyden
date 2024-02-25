import Link from "next/link";

import { styled } from "@/styled-system/jsx";

export function DirectoryBadge() {
  return (
    <styled.span backgroundColor="accent.100" px="1" borderRadius="md">
      <Link href="/directory">directory</Link>
    </styled.span>
  );
}
