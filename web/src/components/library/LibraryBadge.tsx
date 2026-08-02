import Link from "next/link";

import { Node } from "@/api/openapi-schema";
import { Text } from "@/components/ui/text";
import { styled } from "@/styled-system/jsx";

export function LibraryBadge() {
  return (
    <styled.span
      position="relative"
      backgroundColor="accent.default"
      color="accent.text"
      px="1"
      borderRadius="md"
    >
      <Link href="/l">library</Link>
    </styled.span>
  );
}

export function LibraryPageBadge(props: Node) {
  const url = `/l/${props.slug}`;
  return (
    <styled.span
      position="relative"
      backgroundColor="accent.default"
      color="accent.text"
      px="1"
      borderRadius="md"
      lineClamp="1"
    >
      <Link href={url}>{props.name}</Link>
    </styled.span>
  );
}

// TODO: Make this a recipe component.
export function NewBadge() {
  return (
    <Text
      as="span"
      variant="metadata"
      fontWeight="bold"
      backgroundColor="accent.default"
      color="accent.text"
      px="1"
      py="0.5"
      borderRadius="sm"
    >
      New
    </Text>
  );
}
