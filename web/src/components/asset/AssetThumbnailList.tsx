import { Asset } from "@/api/openapi-schema";
import { AssetThumbnail } from "@/components/asset/AssetThumbnail";
import { Box, HStack } from "@/styled-system/jsx";

export type Props = {
  assets: Asset[];
};

export function AssetThumbnailList({ assets }: Props) {
  if (assets.length === 0) {
    return null;
  }

  return (
    <HStack
      w="full"
      overflowX="scroll"
      overflowY="hidden"
      mb="-scrollGutter"
      scrollSnapType="x"
      scrollSnapStrictness="mandatory"
    >
      <HStack w="full" h="20" maxW="full">
        {assets.map((a, i) => (
          // Sizing for next/image is measured in px, size tokens are basically
          // 4X, so size token 20 used above is equal to 80px, so we pass 80 here.
          <Box
            key={a.id}
            position="relative"
            scrollSnapAlign="start"
            scrollSnapStop="always"
          >
            <AssetThumbnail
              asset={a}
              set={assets}
              setIndex={i}
              width={120}
              height={120}
            />
          </Box>
        ))}
      </HStack>
    </HStack>
  );
}
