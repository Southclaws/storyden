import { Flex, Image, Text } from "@chakra-ui/react";
import { Anchor } from "../site/Anchor";

import { useProfileReference } from "./useProfileReference";

type Props = {
  handle: string;
};

export function ProfileReference(props: Props) {
  const { fallback, src } = useProfileReference(props.handle);

  return (
    <Anchor href={`/p/${props.handle}`}>
      <Flex gap={1}>
        <Image
          borderRadius="full"
          boxSize={5}
          src={src}
          fallbackSrc={fallback}
          alt="pic"
        />
        <Text>@{props.handle}</Text>
      </Flex>
    </Anchor>
  );
}
