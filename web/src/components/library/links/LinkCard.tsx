import { LinkReference } from "src/api/openapi-schema";

import { LinkButton } from "@/components/ui/link-button";
import { Card } from "@/components/ui/rich-card";
import {
  ButtonVariantProps,
  RichCardVariantProps,
} from "@/styled-system/recipes";
import { getAssetURL } from "@/utils/asset";

export type Props = {
  link: LinkReference;
} & RichCardVariantProps;

export function LinkCard({ link, ...rest }: Props) {
  const title = link.title || link.url;
  const asset = link.primary_image;

  const linkPagePath = `/links/${link.slug}`;

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
      <LinkRefButton link={link} />
    </Card>
  );
}

type LinkRefButtonProps = { link: LinkReference } & ButtonVariantProps;

export function LinkRefButton({ link, ...rest }: LinkRefButtonProps) {
  return (
    <LinkButton size="xs" variant="subtle" href={link.url} {...rest}>
      {link.domain}
    </LinkButton>
  );
}
