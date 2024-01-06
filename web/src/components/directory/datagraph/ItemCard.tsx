import { Item } from "src/api/openapi/schemas";
import {
  DirectoryPath,
  joinDirectoryPath,
} from "src/screens/directory/datagraph/useDirectoryPath";
import { Heading3 } from "src/theme/components/Heading/Index";

import { Box, LinkBox, LinkOverlay, VStack, styled } from "@/styled-system/jsx";
import { Card, FrostedGlass } from "@/styled-system/patterns";

type Props = {
  item: Item;
  directoryPath: DirectoryPath;
};

export function ItemCard({ item, directoryPath }: Props) {
  const slug = joinDirectoryPath(directoryPath, item.slug);
  return (
    <styled.article containerType="inline-size">
      <LinkBox
        className={Card({ kind: "edge", display: "grid" })}
        w="full"
        h="full"
        aspectRatio="square"
        gridTemplateAreas='"x"'
      >
        <styled.img
          src={item.image_url}
          height="full"
          width="full"
          objectPosition="top"
          objectFit="cover"
          aspectRatio="square"
          gridArea="x"
        />

        <VStack gridArea="x" alignItems="center" justifyContent="end">
          <Box
            className={FrostedGlass()}
            height="min"
            p="2"
            wordBreak="break-all"
          >
            <Heading3 className="fluid-font-size" lineClamp={1}>
              {/* TODO: Next link */}
              <LinkOverlay href={`/directory/${slug}`}>{item.name}</LinkOverlay>
            </Heading3>
            <styled.p lineClamp={1}>{item.description}</styled.p>
          </Box>
        </VStack>
      </LinkBox>
    </styled.article>
  );
}
