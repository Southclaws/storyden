import { Flex, Text } from "@chakra-ui/react";
import { formatDistanceToNow } from "date-fns";
import { ProfileReference } from "src/components/ProfileReference/ProfileReference";

type Props = {
  author: string;
  time: Date;
};

export function Byline(props: Props) {
  return (
    <Flex alignItems="center" gap={2} fontSize="sm" color="blackAlpha.700">
      <ProfileReference handle={props.author} />
      <span>•</span>
      <Text>{formatDistanceToNow(props.time)} ago</Text>
    </Flex>
  );
}
