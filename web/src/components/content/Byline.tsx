import { differenceInSeconds, formatDistanceToNow } from "date-fns";

import { ProfileReference } from "src/api/openapi/schemas";
import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";
import { Timestamp } from "src/components/site/Timestamp";
import { formatDistanceDefaults } from "src/utils/date";

import { Flex, styled } from "@/styled-system/jsx";

type Props = {
  href: string;
  author: ProfileReference;
  time: Date;
  updated: Date;
  more?: React.ReactElement;
};

export function Byline(props: Props) {
  const created = formatDistanceToNow(props.time, formatDistanceDefaults);
  const updated =
    differenceInSeconds(props.time, props.updated) > 0
      ? formatDistanceToNow(props.updated, formatDistanceDefaults)
      : undefined;

  return (
    <Flex alignItems="center" justify="space-between" minWidth="0">
      <Flex
        alignItems="center"
        gap="0"
        minWidth="0"
        fontSize="sm"
        color="gray.500"
      >
        <ProfilePill profileReference={props.author} />
        <styled.span pr="2">•</styled.span>
        <Timestamp created={created} updated={updated} href={props.href} />
      </Flex>

      {props.more}
    </Flex>
  );
}
