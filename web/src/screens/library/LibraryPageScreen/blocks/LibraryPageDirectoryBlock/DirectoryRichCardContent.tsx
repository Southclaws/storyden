import Link from "next/link";
import { ReactNode } from "react";

import { EmptyIcon } from "@/components/ui/icons/Empty";
import { Text } from "@/components/ui/text";
import { HStack, LStack, styled } from "@/styled-system/jsx";

import {
  DirectoryCardProperties,
  DirectoryFieldEditor,
} from "./DirectoryCardFields";
import { DirectoryPropertyBadges } from "./DirectoryPropertyBadges";
import { DirectoryPropertyProps, RichCardShape } from "./directoryCard";

type ImageFallbackProps = Pick<DirectoryPropertyProps, "node"> & {
  showTitle: boolean;
};

export function DirectoryRichCardImageFallback({
  node,
  showTitle,
}: ImageFallbackProps): ReactNode {
  if (showTitle) {
    return <EmptyIcon />;
  }

  return (
    <styled.span data-search-highlight lineClamp={3}>
      {node.name}
    </styled.span>
  );
}

type HeaderProps = DirectoryPropertyProps & {
  shape: RichCardShape;
  showTitle: boolean;
  url: string;
};

export function DirectoryRichCardHeader({
  shape,
  showTitle,
  url,
  ...props
}: HeaderProps) {
  let title: ReactNode;

  if (shape === "fill" && showTitle) {
    title = <OverlayTitle {...props} url={url} />;
  }

  return (
    <HStack
      gap="1"
      flexWrap="wrap"
      pr="0"
      data-editing={props.editing}
      css={{ "&[data-editing=true]": { pr: "7.5" } }}
    >
      {title}
      <DirectoryPropertyBadges {...props} placement="header" />
    </HStack>
  );
}

type FieldProps = Pick<
  DirectoryPropertyProps,
  "node" | "editing" | "onFieldValueChange"
>;

function OverlayTitle({
  node,
  editing,
  onFieldValueChange,
  url,
}: FieldProps & { url: string }) {
  let title: ReactNode = <Link href={url}>{node.name}</Link>;

  if (editing) {
    title = (
      <DirectoryFieldEditor
        node={node}
        onFieldValueChange={onFieldValueChange}
        fid="fixed:name"
        label="Name"
        value={node.name}
        elevated
      />
    );
  }

  return (
    <styled.h3
      fontSize="sm"
      fontWeight="semibold"
      minW="0"
      data-editing={editing}
      css={{ "&[data-editing=true]": { w: "full" } }}
    >
      {title}
    </styled.h3>
  );
}

type BodyProps = DirectoryPropertyProps & {
  shape: RichCardShape;
  showDescription: boolean;
  showLink: boolean;
};

export function DirectoryRichCardBody({
  shape,
  showDescription,
  showLink,
  ...props
}: BodyProps) {
  const minimal = shape === "minimal";
  return (
    <LStack gap="1" w="full" minW="0">
      {showDescription && <RichCardDescription {...props} minimal={minimal} />}
      {showLink && <RichCardLink {...props} minimal={minimal} />}
      <DirectoryCardProperties {...props} linkValues={minimal} elevated />
    </LStack>
  );
}

function RichCardDescription({
  node,
  editing,
  onFieldValueChange,
  minimal,
}: FieldProps & { minimal: boolean }) {
  if (editing) {
    return (
      <DirectoryFieldEditor
        node={node}
        onFieldValueChange={onFieldValueChange}
        fid="fixed:description"
        label="Description"
        value={node.description}
        elevated
      />
    );
  }

  if (minimal) {
    return (
      <styled.p
        fontSize="sm"
        color="text.subtle"
        lineClamp={1}
        data-search-highlight
      >
        {node.description}
      </styled.p>
    );
  }

  return (
    <Text as="p" variant="supporting" lineClamp={2} data-search-highlight>
      {node.description}
    </Text>
  );
}

function RichCardLink({
  node,
  editing,
  onFieldValueChange,
  minimal,
}: FieldProps & { minimal: boolean }) {
  if (editing) {
    return (
      <DirectoryFieldEditor
        node={node}
        onFieldValueChange={onFieldValueChange}
        fid="fixed:link"
        label="Link"
        value={node.link?.url ?? ""}
        elevated
      />
    );
  }

  if (!node.link) {
    return null;
  }

  if (minimal) {
    return (
      <styled.a
        href={node.link.url}
        fontSize="sm"
        color="text.subtle"
        lineClamp={1}
      >
        {node.link.title || node.link.url}
      </styled.a>
    );
  }

  return (
    <styled.a
      href={node.link.url}
      position="relative"
      zIndex="var(--card-content-index)"
      fontSize="xs"
      color="text.muted"
      maxW="full"
      truncate
    >
      {node.link.title || node.link.url}
    </styled.a>
  );
}
