import { Link as LinkSchema } from "src/api/openapi/schemas";

import { LinkButton } from "@/components/ui/link-button";
import { Card } from "@/components/ui/rich-card";
import { HStack } from "@/styled-system/jsx";
import { RichCardVariantProps } from "@/styled-system/recipes";

export type Props = {
  link: LinkSchema;
} & RichCardVariantProps;

export function LinkCard({ link, ...rest }: Props) {
  const asset = link.assets?.[0] ?? undefined;
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
