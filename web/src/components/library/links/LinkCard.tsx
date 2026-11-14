import { LinkReference } from "@/api/openapi-schema";
import { EmptyIcon } from "@/components/ui/icons/Empty";
import { LinkButton } from "@/components/ui/link-button";
import { RelativeTime } from "@/components/ui/relative-time";
import { StatusBadge } from "@/components/ui/status-badge";
import { Card } from "@/components/ui/surface";
import { Text } from "@/components/ui/text";
import { css } from "@/styled-system/css";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import {
  ButtonVariantProps,
  RichCardVariantProps,
} from "@/styled-system/recipes";
import { getAssetImageProps, getAssetURL } from "@/utils/asset";

import { getLinkStatus } from "./link-presentation";

export type Props = {
  link: LinkReference;
  disableAnchors?: boolean;
} & RichCardVariantProps;

export function LinkCard({
  link,
  disableAnchors,
  shape = "row",
  ...rest
}: Props) {
  const title = link.title || link.url;
  const asset = link.primary_image;
  const status = getLinkStatus(link);

  const linkPagePath = `/links/${link.slug}`;
  const timestamp = (
    <Text variant="metadata" minW="0" truncate>
      Saved <RelativeTime value={link.createdAt} title={link.createdAt} />
    </Text>
  );
  const statusBadge = status && (
    <StatusBadge tone={status.tone} size="sm" flexShrink="0">
      {status.label}
    </StatusBadge>
  );

  let footer = (
    <WStack gap="2" flexWrap="wrap" minW="0">
      <HStack minW="0" maxW="full" flexWrap="wrap">
        <LinkRefButton link={link} size="sm" />
        {timestamp}
      </HStack>
      {statusBadge}
    </WStack>
  );

  if (shape === "box") {
    footer = (
      <LStack gap="0" minW="0" w="full">
        <LinkRefButton link={link} size="sm" />
        <WStack gap="2" minW="0" w="full" minH="6" alignItems="end">
          {timestamp}
          {statusBadge}
        </WStack>
      </LStack>
    );
  }

  return (
    <Card
      id={link.slug}
      title={title}
      url={linkPagePath}
      text={link.description}
      image={getAssetURL(asset?.path, "m")}
      imageFallback={shape === "box" && <EmptyIcon aria-hidden />}
      imageProps={{ ...getAssetImageProps(asset), loading: "lazy" }}
      titleIcon={
        link.favicon_image && (
          <styled.img
            src={getAssetURL(link.favicon_image.path, "s")}
            alt=""
            width="4"
            height="4"
            objectFit="contain"
            loading="lazy"
          />
        )
      }
      shape={shape}
      disableAnchors={disableAnchors}
      {...rest}
    >
      {footer}
    </Card>
  );
}

type LinkRefButtonProps = { link: LinkReference } & ButtonVariantProps;

export function LinkRefButton({ link, ...rest }: LinkRefButtonProps) {
  return (
    <LinkButton
      variant="subtle"
      href={link.url}
      minW="0"
      maxW="full"
      className={css({ "& > span:first-child": { minWidth: "0" } })}
      {...rest}
    >
      <styled.span truncate>{link.domain}</styled.span>
    </LinkButton>
  );
}
