import { Asset, Link } from "src/api/openapi/schemas";

import { Box, LinkBox, LinkOverlay, VStack, styled } from "@/styled-system/jsx";

type Props = {
  link: Link;
  asset?: Asset;
};

export function LinkView({ link, asset }: Props) {
  return (
    <LinkBox
      display="flex"
      w="full"
      borderRadius="xl"
      bgColor="accent.100"
      overflow="hidden"
      height="24"
      shadow="sm"
    >
      {asset && (
        <Box flexGrow="1" flexShrink="0" width="32">
          <styled.img
            src={asset.url}
            height="full"
            objectPosition="left"
            objectFit="cover"
          />
        </Box>
      )}
      <VStack
        w="full"
        alignItems="start"
        justifyContent="space-evenly"
        gap="0"
        p="2"
      >
        <styled.h2 fontWeight="bold">
          <LinkOverlay target="blank" href={link.url}>
            {link.title}
          </LinkOverlay>
        </styled.h2>
        <Box maxLines={1} overflow="hidden">
          <styled.p maxLines={2} lineClamp={2}>
            {link.description}
          </styled.p>
          <br />
        </Box>
      </VStack>
    </LinkBox>
  );
}
