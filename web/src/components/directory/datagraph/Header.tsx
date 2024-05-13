import { Node } from "src/api/openapi/schemas";
import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";
import { Heading1 } from "src/theme/components/Heading/Index";
import { Link } from "src/theme/components/Link";

import { Box, HStack, Stack, VStack, styled } from "@/styled-system/jsx";

type Props = Node;

export function DatagraphHeader(props: Props) {
  const asset = props.assets?.[0];
  return (
    <Stack
      w="full"
      direction={{
        base: "column-reverse",
        sm: "row",
      }}
      gap="2"
    >
      <VStack alignItems="start" w="full" minW="0">
        <Heading1>{props.name}</Heading1>

        <HStack>
          <styled.p color="fg.subtle">Maintained by</styled.p>
          <ProfilePill profileReference={props.owner} />
        </HStack>

        {props.link && (
          <Box w="full">
            <Link href={props.link?.url} w="full" size="sm">
              {props.link?.url}
            </Link>
          </Box>
        )}

        <styled.p>{props.description}</styled.p>
      </VStack>

      {asset && (
        <HStack w="full" h="full" maxH="64" justify="center" minW="0">
          <styled.img maxHeight="64" borderRadius="lg" src={asset.url} />
        </HStack>
      )}
    </Stack>
  );
}
