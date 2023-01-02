import { Flex, Image, Text } from "@chakra-ui/react";
import { formatDistanceToNow } from "date-fns";

type Props = {
  author: string;
  time: Date;
};

export function Byline(props: Props) {
  return (
    <Flex alignItems="center" gap={2} fontSize="sm" color="blackAlpha.700">
      <Image
        borderRadius="full"
        boxSize={8}
        src={`/api/v1/accounts/${props.author}/avatar`}
        fallbackSrc="/logo_50x50.png"
        alt="pic"
      />
      <Text>@{props.author}</Text>
      <span>•</span>
      <Text>{formatDistanceToNow(props.time)} ago</Text>
    </Flex>
  );
}
