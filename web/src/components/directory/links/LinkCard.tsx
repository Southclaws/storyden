import { LinkReference } from "src/api/openapi-schema";

import { LinkButton } from "@/components/ui/link-button";
import { Card } from "@/components/ui/rich-card";
import { HStack } from "@/styled-system/jsx";
import { RichCardVariantProps } from "@/styled-system/recipes";

export type Props = {
  link: LinkReference;
} & RichCardVariantProps;

export function LinkCard({ link, ...rest }: Props) {
  const asset = link.primary_image;
  const domainSearch = `/l?q=${link.domain}`;

  return (
    <Card
      id={link.slug}
      title={link.title ?? link.url}
      url={link.url}
      text={link.description}
      image={asset?.url}
      shape="row"
      {...rest}
    >
      <HStack>
        <LinkButton size="xs" href={`/l/${link.slug}`} variant="ghost">
          View in directory
        </LinkButton>
        <LinkButton size="xs" href={domainSearch} variant="ghost">
          More from this site
        </LinkButton>
      </HStack>
    </Card>
  );
}
