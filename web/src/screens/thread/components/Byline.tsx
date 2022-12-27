import { Flex, Text } from "@chakra-ui/react";
import { formatDistanceToNow } from "date-fns";

type Props = {
  author: string;
  time: Date;
};

export function Byline(props: Props) {
  return (
    <Flex gap={2} fontSize="sm" color="blackAlpha.700">
      <Text>@{props.author}</Text>
      <span>•</span>
      <Text>{formatDistanceToNow(props.time)}</Text>
    </Flex>
  );
}
