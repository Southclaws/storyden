import { ArrowTopRightOnSquareIcon } from "@heroicons/react/24/outline";

import { LinkWithRefs } from "src/api/openapi/schemas";
import { PostRefList } from "src/components/feed/common/PostRef/PostRefList";
import { Link } from "src/theme/components/Link";

import { Flex, LinkBox, LinkOverlay, styled } from "@/styled-system/jsx";

type Props = {
  link: LinkWithRefs;
};

export function LinkView({ link }: Props) {
  const mainImage = link.assets[0]?.url;
  const images = mainImage ? link.assets.slice(1).map((v) => v.url) : undefined;

  const domainSearch = `/l?q=${link.domain}`;

  return (
    <Flex flexDir="column" gap="2">
      <LinkBox>
        <LinkOverlay
          display="flex"
          alignItems="center"
          color="fg.subtle"
          href={link.url}
        >
          {link.domain}&nbsp;
          <ArrowTopRightOnSquareIcon height="1rem" />
        </LinkOverlay>

        <styled.h1 fontSize="heading.variable.2">{link.title}</styled.h1>
      </LinkBox>

      <styled.p color="fg.muted">{link.description}</styled.p>

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

      <Link w="min" size="xs" href={domainSearch}>
        More from this site
      </Link>

      <styled.h2 fontSize="heading.variable.3">Shared in</styled.h2>
      <PostRefList items={link.threads} />

      <styled.h2 fontSize="heading.variable.3">Mentioned in replies</styled.h2>
      <PostRefList items={link.posts} />
    </Flex>
  );
}
