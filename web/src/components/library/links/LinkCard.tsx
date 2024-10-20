import { LinkReference } from "src/api/openapi-schema";

import { LinkButton } from "@/components/ui/link-button";
import { Card } from "@/components/ui/rich-card";
import { HStack } from "@/styled-system/jsx";
import { RichCardVariantProps } from "@/styled-system/recipes";
import { getAssetURL } from "@/utils/asset";

export type Props = {
  link: LinkReference;
} & RichCardVariantProps;

export function LinkCard({ link, ...rest }: Props) {
  const title = link.title || link.url;
  const asset = link.primary_image;
  const domainSearch = `/links?q=${link.domain}`;

  return (
    <Card
      id={link.slug}
      title={title}
      url={link.url}
      text={link.description}
      image={getAssetURL(asset?.path)}
      shape="row"
      {...rest}
    >
      <HStack>
        <LinkButton size="xs" href={`/links/${link.slug}`} variant="ghost">
          View in library
        </LinkButton>
        <LinkButton size="xs" href={domainSearch} variant="ghost">
          More from this site
        </LinkButton>
      </HStack>
    </Card>
  );
}
