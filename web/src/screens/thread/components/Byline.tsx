import { Flex, HStack, Text } from "@chakra-ui/react";
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
    <Flex
      alignItems={{
        // base: "start",
        md: "center",
      }}
      gap={1}
      fontSize="sm"
      color="blackAlpha.700"
      flexDir={{
        // base: "column",
        md: "row",
      }}
    >
      <ProfileReference handle={props.author} />

      <HStack>
        <Text as="span">•</Text>
      </HStack>

      <HStack>
        <Timestamp created={created} updated={updated} href={props.href} />
      </HStack>
    </Flex>
  );
}
