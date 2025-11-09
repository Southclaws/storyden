import Link from "next/link";

import { getSettings } from "@/lib/settings/settings-server";
import { styled } from "@/styled-system/jsx";

export async function Title() {
  const { title } = await getSettings();

  return (
    <styled.h1
      fontSize="lg"
      fontWeight="bold"
      textWrap="nowrap"
      overflow="hidden"
      textOverflow="ellipsis"
    >
      <Link href="/">{title}</Link>
    </styled.h1>
  );
}
