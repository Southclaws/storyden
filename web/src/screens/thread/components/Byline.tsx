import { Flex, HStack, Text } from "@chakra-ui/react";
import { differenceInSeconds, formatDistanceToNow } from "date-fns";
import { ProfileReference } from "src/components/ProfileReference/ProfileReference";
import { Anchor } from "src/components/site/Anchor";
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
    differenceInSeconds(props.time, props.updated) > 0 &&
    formatDistanceToNow(props.updated, formatDistanceDefaults);

  return (
    <HStack justifyContent="space-between">
      <Flex
        alignItems={{
          // base: "start",
          md: "center",
        }}
        gap={{
          base: 1,
          md: 2,
        }}
        fontSize="sm"
        color="blackAlpha.700"
        flexDir={{
          // base: "column",
          md: "row",
        }}
      >
        <ProfileReference handle={props.author} />

        <Text
          as="span"
          display={{
            // base: "none",
            md: "inline",
          }}
        >
          •
        </Text>

        <Text>
          <Anchor href={props.href}>{created} ago</Anchor>
          {updated && <> (updated {updated} ago)</>}
        </Text>
      </Flex>

      {props.more}
    </HStack>
  );
}
