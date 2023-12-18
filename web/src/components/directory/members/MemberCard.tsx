import { formatDistanceToNow } from "date-fns";

import { PublicProfile } from "src/api/openapi/schemas";

import { Box, styled } from "@/styled-system/jsx";

export function MemberCard(props: PublicProfile) {
  const createdAt = new Date(props.createdAt);
  return (
    <styled.tr
      w="full"
      overflow="hidden"
      boxShadow="md"
      borderRadius="lg"
      css={{
        "&[data-selected=true]": {
          outlineStyle: "dashed",
          outlineOffset: "-0.5",
          outlineWidth: "medium",
          outlineColor: "accent.200",
        },
      }}
    >
      <styled.td p="1" maxW="24" textOverflow="ellipsis" overflow="auto">
        <styled.h1 color="fg.default">{props.name}</styled.h1>
        <styled.h2 color="fg.subtle">@{props.handle}</styled.h2>
      </styled.td>

      <styled.td p="1" overflow="auto">
        <styled.p color="fg.default">
          <styled.time wordBreak="keep-all" textWrap="nowrap">
            {formatDistanceToNow(createdAt)}
          </styled.time>{" "}
          ago
        </styled.p>
      </styled.td>

      <styled.td p="1" maxW="64">
        <styled.p overflow="auto" textOverflow="ellipsis">
          {props.id}
        </styled.p>
      </styled.td>
    </styled.tr>
  );
}
