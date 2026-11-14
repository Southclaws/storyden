import type { Link } from "@/api/openapi-schema";
import { AssetThumbnailList } from "@/components/asset/AssetThumbnailList";
import { Badge } from "@/components/ui/badge";
import { Breadcrumbs } from "@/components/ui/breadcrumbs";
import { SearchIcon } from "@/components/ui/icons/Search";
import { LinkButton } from "@/components/ui/link-button";
import { MetaGrid, MetaItem } from "@/components/ui/meta-grid";
import { PageHeader } from "@/components/ui/page-header";
import { RelativeTime } from "@/components/ui/relative-time";
import { SectionHeading } from "@/components/ui/section-heading";
import { StatusBadge } from "@/components/ui/status-badge";
import { Card, CardRows } from "@/components/ui/surface";
import { Text } from "@/components/ui/text";
import { Grid, HStack, LStack, styled } from "@/styled-system/jsx";
import { getAssetImageProps } from "@/utils/asset";

import { getLinkStatus } from "./link-presentation";

const REACHABILITY_LABEL = {
  alive: "Reachable at last check",
  offline: "Couldn't reach site at last check",
  unknown: "No conclusive result",
};

const SCRAPE_LABEL = {
  pending: "Awaiting metadata retrieval",
  success: "Metadata retrieved successfully",
  failed: "Couldn't retrieve metadata",
};

/** Link identity, saved metadata, and the community content that references it. */
export function LinkDetails({ link }: { link: Link }) {
  const title = link.title || link.url;
  const status = getLinkStatus(link);
  const hasReferences = link.nodes.length > 0 || link.posts.length > 0;
  const { width, height, ...imageProps } = getAssetImageProps(
    link.primary_image,
  );
  let savedImages = link.assets;

  if (savedImages.length === 0 && link.primary_image) {
    savedImages = [link.primary_image];
  }

  return (
    <LStack gap="4">
      <LStack gap="2">
        <PageHeader
          title={title}
          gap="3"
          navigation={
            <Breadcrumbs
              index={{ label: "Links", href: "/links" }}
              crumbs={[{ label: title, href: `/links/${link.slug}` }]}
            />
          }
          actions={
            <LinkButton variant="subtle" size="md" href={link.url}>
              Visit site
            </LinkButton>
          }
        />

        <HStack gap="2" flexWrap="wrap" minW="0" w="full">
          <Text variant="metadata" overflowWrap="anywhere" minW="0">
            {link.domain}
          </Text>
          <LinkButton
            variant="ghost"
            href={`/links?q=${encodeURIComponent(link.domain)}`}
          >
            <SearchIcon />
            More from this site
          </LinkButton>
          {status && (
            <StatusBadge tone={status.tone} size="sm">
              {status.label}
            </StatusBadge>
          )}
        </HStack>
      </LStack>

      {(link.description || link.primary_image) && (
        <Grid
          w="full"
          gridTemplateColumns="repeat(auto-fit, minmax(min(100%, 16rem), 1fr))"
          gap="4"
          alignItems="start"
        >
          {link.description && (
            <Text variant="supporting" minW="0" overflowWrap="anywhere">
              {link.description}
            </Text>
          )}

          {link.primary_image && (
            <styled.img
              {...imageProps}
              htmlWidth={width}
              htmlHeight={height}
              alt=""
              width="full"
              maxWidth="full"
              height="auto"
              objectFit="contain"
              objectPosition="left top"
              borderRadius="md"
            />
          )}
        </Grid>
      )}

      <LStack as="section" gap="3" aria-labelledby="link-metadata">
        <SectionHeading id="link-metadata">About this link</SectionHeading>
        <MetaGrid>
          <MetaItem label="Saved">
            <RelativeTime value={link.createdAt} title={link.createdAt} />
          </MetaItem>
          <MetaItem label="Last checked">
            {link.liveness_checked_at ? (
              <RelativeTime
                value={link.liveness_checked_at}
                title={link.liveness_checked_at}
              />
            ) : (
              "Not checked yet"
            )}
            {link.liveness_checked_at && link.liveness_status && (
              <Text variant="metadata">
                {REACHABILITY_LABEL[link.liveness_status]}
              </Text>
            )}
          </MetaItem>
          <MetaItem label="Metadata retrieved">
            {link.scraped_at ? (
              <RelativeTime value={link.scraped_at} title={link.scraped_at} />
            ) : (
              "Not recorded"
            )}
          </MetaItem>
          {link.scrape_attempted_at && (
            <MetaItem label="Latest metadata attempt">
              <RelativeTime
                value={link.scrape_attempted_at}
                title={link.scrape_attempted_at}
              />
              {link.scrape_status && (
                <Text variant="metadata">
                  {SCRAPE_LABEL[link.scrape_status]}
                </Text>
              )}
            </MetaItem>
          )}
        </MetaGrid>
      </LStack>

      {savedImages.length > 0 && (
        <LStack as="section" gap="3" aria-labelledby="link-images">
          <SectionHeading id="link-images">Saved images</SectionHeading>
          <AssetThumbnailList assets={savedImages} />
        </LStack>
      )}

      {hasReferences && (
        <LStack as="section" gap="3" aria-labelledby="link-references">
          <SectionHeading id="link-references">Referenced in</SectionHeading>
          <CardRows>
            {link.nodes.map((node) => (
              <Card
                key={node.id}
                id={`link-node-${node.id}`}
                title={node.name}
                text={node.description}
                url={`/l/${node.slug}`}
                headingAs="h3"
                shape="row"
              >
                <Badge>Library page</Badge>
              </Card>
            ))}
            {link.posts.map((post) => (
              <Card
                key={post.id}
                id={`link-post-${post.id}`}
                title={post.title || "Untitled post"}
                text={post.description}
                url={`/t/locate/${post.id}`}
                headingAs="h3"
                shape="row"
              >
                <Badge>Discussion</Badge>
              </Card>
            ))}
          </CardRows>
        </LStack>
      )}
    </LStack>
  );
}
