import { Flex, Text } from "@chakra-ui/react";
import { differenceInSeconds, formatDistanceToNow } from "date-fns";

import { ProfileReference } from "src/components/ProfileReference/ProfileReference";
import { Timestamp } from "src/components/Timestamp";
import { formatDistanceDefaults } from "src/utils/date";

type Props = {
  href: string;
  author: string;
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
      <ProfileReference handle={props.author} />
      <Text as="span" pr={2}>
        •
      </Text>
      <Timestamp created={created} updated={updated} href={props.href} />
    </Flex>
  );
}
