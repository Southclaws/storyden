import { Asset, Link } from "src/api/openapi-schema";

import { Box, LinkOverlay, VStack, styled } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

type Props = {
  link: Link;
  asset?: Asset;
};

export function LinkView({ link, asset }: Props) {
  return (
    <Box
      position="relative"
      display="flex"
      w="full"
      borderRadius="xl"
      bgColor="bg.subtle"
      overflow="hidden"
      height="24"
      shadow="sm"
    >
      {asset && (
        <Box flexGrow="1" flexShrink="0" width="32">
          <styled.img
            src={getAssetURL(asset.filename)}
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
        <Box lineClamp={1} overflow="hidden">
          <styled.p lineClamp={2}>{link.description}</styled.p>
          <br />
        </Box>
      </VStack>
    </Box>
  );
}
