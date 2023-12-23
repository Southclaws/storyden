import { formatDistanceToNow } from "date-fns";

import { PublicProfile } from "src/api/openapi/schemas";
import { MemberOptionsTrigger } from "src/components/member/MemberOptions/MemberOptionsTrigger";
import { Avatar } from "src/components/site/Avatar/Avatar";

import { LinkBox, LinkOverlay, VStack, styled } from "@/styled-system/jsx";

type Props = PublicProfile & {
  onChange?: () => void;
};

export function MemberCard(props: Props) {
  const createdAt = new Date(props.createdAt);
  return (
    <styled.tr
      display="contents"
      w="full"
      overflow="hidden"
      background="bg.default"
      boxShadow="md"
      borderRadius="lg"
      _hover={{
        boxShadow: "lg",
      }}
    >
      <styled.td overflow="auto" flexShrink="1" flexGrow="0">
        <LinkBox display="flex" alignItems="center" gap="2">
          <Avatar flexShrink="0" handle={props.handle} />
          <VStack alignItems="start" gap="0">
            <styled.h1
              color="fg.default"
              textOverflow="ellipsis"
              overflow="auto"
            >
              {props.name}
            </styled.h1>
            <styled.h2
              color="fg.subtle"
              textOverflow="ellipsis"
              overflow="auto"
            >
              <LinkOverlay href={`/p/${props.handle}`}>
                @{props.handle}
              </LinkOverlay>
            </styled.h2>
          </VStack>
        </LinkBox>
      </styled.td>

      <styled.td overflow="auto" flexShrink="0" flexGrow="1">
        <VStack alignItems="end" gap="0">
          <styled.p color="fg.default">
            <styled.time wordBreak="keep-all" textWrap="nowrap">
              {formatDistanceToNow(createdAt)}
            </styled.time>
            &nbsp;ago
          </styled.p>

          {props.deletedAt && (
            <styled.p color="fg.destructive" wordBreak="keep-all">
              (suspended)
            </styled.p>
          )}
        </VStack>
      </styled.td>

      <styled.td display="flex" alignItems="center">
        <MemberOptionsTrigger {...props} />
      </styled.td>
    </styled.tr>
  );
}
