import { Link as LinkSchema } from "src/api/openapi/schemas";
import { Card } from "src/theme/components/Card";
import { Link } from "src/theme/components/Link";

import { CardVariantProps } from "@/styled-system/recipes";

export type Props = {
  link: LinkSchema;
} & CardVariantProps;

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
      <Link size="sm" href={domainSearch}>
        More from {link.domain}
      </Link>
    </Card>
  );
}
