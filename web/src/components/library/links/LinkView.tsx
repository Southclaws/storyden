import { ArrowTopRightOnSquareIcon } from "@heroicons/react/24/outline";

import { Link } from "src/api/openapi-schema";

import { PostRefList } from "@/components/feed/PostRef/PostRefList";
import { LinkButton } from "@/components/ui/link-button";
import { Box, Flex, HStack, LinkOverlay, styled } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

type Props = {
  link: Link;
};

export function LinkView({ link }: Props) {
  const mainImage = getAssetURL(link.assets?.[0]?.filename);
  const images = mainImage
    ? link.assets.slice(1).map((v) => getAssetURL(v.filename))
    : undefined;

  const domainSearch = `/links?q=${link.domain}`;

  return (
    <Flex flexDir="column" gap="2">
      <Box position="relative">
        <HStack justify="space-between">
          <LinkOverlay
            display="flex"
            alignItems="center"
            color="fg.subtle"
            href={link.url}
          >
            {link.domain}&nbsp;
            <ArrowTopRightOnSquareIcon height="1rem" />
          </LinkOverlay>

          <LinkButton w="min" size="xs" href={domainSearch}>
            More from this site
          </LinkButton>
        </HStack>

        {link.title ? (
          <styled.h1 fontSize="heading.variable.2">{link.title}</styled.h1>
        ) : (
          <styled.h1 fontSize="heading.variable.2" lineClamp={1}>
            (no title) {link.slug}
          </styled.h1>
        )}
      </Box>

      <styled.p color="fg.muted">
        {link.description || "No description was found at this link's site."}
      </styled.p>

      <Flex flexDir={{ base: "column", md: "row" }}>
        {mainImage && (
          <styled.img
            width="full"
            maxWidth="96"
            maxHeight="48"
            objectFit="cover"
            aspectRatio="wide"
            borderRadius="lg"
            src={mainImage}
          />
        )}

        {images?.map((v) => <>{v}</>)}
      </Flex>

      <styled.h2 fontSize="heading.variable.3">Mentioned in replies</styled.h2>
      <PostRefList items={link.posts} emptyText="not mentioned in any posts" />
    </Flex>
  );
}
