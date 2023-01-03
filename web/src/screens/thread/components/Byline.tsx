import { Flex, Image, Text } from "@chakra-ui/react";
import { formatDistanceToNow } from "date-fns";
import { useByLine } from "./useByLine";

type Props = {
  author: string;
  time: Date;
};

export function Byline(props: Props) {
  const { fallback, src } = useByLine(props.author);
  return (
    <Flex alignItems="center" gap={2} fontSize="sm" color="blackAlpha.700">
      <Image
        borderRadius="full"
        boxSize={5}
        src={src}
        fallbackSrc={fallback}
        alt="pic"
      />
      <Text>@{props.author}</Text>
      <span>•</span>
      <Text>{formatDistanceToNow(props.time)} ago</Text>
    </Flex>
  );
}
