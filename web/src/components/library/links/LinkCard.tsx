import { LinkReference } from "src/api/openapi-schema";

import { LinkButton } from "@/components/ui/link-button";
import { Card } from "@/components/ui/rich-card";
import { RichCardVariantProps } from "@/styled-system/recipes";
import { getAssetURL } from "@/utils/asset";

export type Props = {
  link: LinkReference;
} & RichCardVariantProps;

export function LinkCard({ link, ...rest }: Props) {
  const title = link.title || link.url;
  const asset = link.primary_image;

  const domainSearch = `/links?q=${link.domain}`;
  const linkPagePath = `/links/${link.slug}`;
  const linkURL = link.url;

  return (
    <Card
      id={link.slug}
      title={title}
      url={linkPagePath}
      text={link.description || "(no description)"}
      image={getAssetURL(asset?.path)}
      shape="row"
      {...rest}
    >
      <LinkButton size="xs" variant="subtle" href={linkURL}>
        {link.domain}
      </LinkButton>
    </Card>
  );
}
