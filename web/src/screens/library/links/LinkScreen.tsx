"use client";

import { useLinkGet } from "@/api/openapi-client/links";
import { Link } from "@/api/openapi-schema";
import { AssetThumbnailList } from "@/components/asset/AssetThumbnailList";
import { Unready } from "@/components/site/Unready";
import { Breadcrumbs } from "@/components/ui/Breadcrumbs";
import { Heading } from "@/components/ui/heading";
import { SearchIcon } from "@/components/ui/icons/Search";
import { LinkButton } from "@/components/ui/link-button";
import { Flex, HStack, LStack, styled } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

export type Props = {
  initialLink: Link;
  slug: string;
};

export function LinkScreen(props: Props) {
  const { data, error } = useLinkGet(props.slug, {
    swr: { fallbackData: props.initialLink },
  });

  if (!data) return <Unready error={error} />;

  const link = data;

  const titleLabel = link.title ? link.title : link.url;
  const descriptionLabel =
    link.description || "No description was found at this link's site.";

  const mainImage = getAssetURL(link.primary_image?.path);

  const domainSearch = `/links?q=${link.domain}`;

  const assetsForThumbnails =
    link.assets && link.assets.length > 0
      ? link.assets
      : link.primary_image
        ? [link.primary_image]
        : [];

  return (
    <LStack>
      <Breadcrumbs
        index={{
          label: "Links",
          href: "/links",
        }}
        crumbs={[
          {
            label: titleLabel,
            href: `/links/${link.slug}`,
          },
        ]}
      >
        <LinkButton
          flexShrink="0"
          size="xs"
          variant="subtle"
          href={domainSearch}
        >
          <SearchIcon />
          More from this site
        </LinkButton>
      </Breadcrumbs>

      <Flex
        w="full"
        gap="3"
        flexDirection={{
          base: "column",
          md: "row",
        }}
      >
        <LStack>
          <Heading size="lg">{titleLabel}</Heading>

          <styled.p color="fg.muted">{descriptionLabel}</styled.p>

          <HStack>
            <LinkButton w="min" size="xs" variant="subtle" href={link.url}>
              {link.domain}
            </LinkButton>
          </HStack>
        </LStack>

        <LStack>
          {mainImage && (
            <styled.img
              width="auto"
              maxWidth="full"
              // maxHeight="64"
              borderRadius="lg"
              src={mainImage}
            />
          )}
        </LStack>
      </Flex>

      <AssetThumbnailList assets={assetsForThumbnails} />
    </LStack>
  );
}
