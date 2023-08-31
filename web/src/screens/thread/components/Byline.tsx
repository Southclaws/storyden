import { differenceInSeconds, formatDistanceToNow } from "date-fns";

import { ProfileReference as ProfileReferenceSchema } from "src/api/openapi/schemas";
import { ProfileReference } from "src/components/ProfileReference/ProfileReference";
import { Timestamp } from "src/components/Timestamp";
import { formatDistanceDefaults } from "src/utils/date";

import { Flex, styled } from "@/styled-system/jsx";

type Props = {
  href: string;
  author: ProfileReferenceSchema;
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
    <Flex alignItems="center" fontSize="sm" color="gray.500" gap={0}>
      <ProfileReference profileReference={props.author} />
      <styled.span pr={2}>•</styled.span>
      <Timestamp created={created} updated={updated} href={props.href} />
    </Flex>
  );
}
