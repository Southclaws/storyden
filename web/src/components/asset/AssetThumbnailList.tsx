import { Asset } from "@/api/openapi-schema";
import { AssetThumbnail } from "@/components/asset/AssetThumbnail";
import { css } from "@/styled-system/css";
import { Box, HStack } from "@/styled-system/jsx";

export type Props = {
  assets: Asset[];
  align?: "start" | "end";
};

export function AssetThumbnailList({ assets, align = "start" }: Props) {
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
      <HStack
        w="fit"
        minW="fit"
        h="20"
        maxW="full"
        ml={align === "end" ? "auto" : undefined}
      >
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
              className={css({
                objectFit: "cover",
                width: "20",
                height: "20",
                minWidth: "20",
              })}
              asset={a}
              set={assets}
              setIndex={i}
            />
          </Box>
        ))}
      </HStack>
    </HStack>
  );
}
