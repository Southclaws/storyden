import { Link as LinkSchema } from "src/api/openapi/schemas";
import { Card } from "src/theme/components/Card";
import { Link } from "src/theme/components/Link";

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
        <Link size="xs" href={`/l/${link.slug}`} variant="ghost">
          View in directory
        </Link>
        <Link size="xs" href={domainSearch} variant="ghost">
          More from this site
        </Link>
      </HStack>
    </Card>
  );
}
