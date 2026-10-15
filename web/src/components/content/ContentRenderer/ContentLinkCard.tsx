import { type ElementType, useEffect } from "react";

import { useLinkCreate } from "@/api/openapi-client/links";
import { LinkCard } from "@/components/library/links/LinkCard";
import { Button } from "@/components/ui/button";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { LinkButton } from "@/components/ui/link-button";
import { Spinner } from "@/components/ui/spinner";
import { Text } from "@/components/ui/text";
import { normalizeLink } from "@/lib/link/validation";
import { css } from "@/styled-system/css";
import { Center, LStack, styled } from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";

export function ContentLinkCard({
  href,
  dataDiff,
  isEditable,
  isSelected,
  Wrapper = "div",
}: {
  href: string;
  dataDiff?: "insertion" | "deletion";
  isEditable?: boolean;
  isSelected?: boolean;
  Wrapper?: ElementType;
}) {
  const { data, error, isMutating, trigger } = useLinkCreate();
  const fallbackHref = normalizeLink(href);

  useEffect(() => {
    trigger({
      url: href,
    });
  }, [href, trigger]);

  const diffClass = dataDiff ? `diff-node-${dataDiff}` : "";
  const showSelection = isSelected && !isMutating;

  return (
    <Wrapper
      data-node-view-wrapper=""
      style={{ whiteSpace: "normal" }}
      data-diff={dataDiff}
      className={`${diffClass} ${css({
        position: "relative",
        display: "inline-block",
        width: "full",
        outlineWidth: isSelected ? "medium" : "none",
        outlineStyle: "solid",
        outlineColor: isSelected ? "blue.a6" : "transparent",
        borderRadius: "lg",
        userSelect: isEditable ? "none" : "auto",
        // subtle saturation bump, combined with...
        saturate: showSelection ? "150%" : "100%",
        filter: "auto",
        // background mix with subtle selection colour
        background: showSelection ? "blue.5" : "transparent",
        mixBlendMode: showSelection ? "screen" : "normal",
      })}`}
    >
      <div data-no-typography>
        {!data ? (
          <Center w="full" minH="12">
            {error ? (
              <ContentLinkCardError
                href={href}
                error={error}
                isEditable={isEditable}
                loading={isMutating}
                onRetry={() => trigger({ url: href })}
              />
            ) : !isEditable && fallbackHref ? (
              <LinkButton variant="subtle" href={fallbackHref} maxW="full">
                <Text truncate>{href}</Text>
              </LinkButton>
            ) : (
              <Spinner />
            )}
          </Center>
        ) : (
          <LinkCard shape="row" link={data} disableAnchors={isEditable} />
        )}
      </div>
    </Wrapper>
  );
}

function ContentLinkCardError({
  href,
  error,
  isEditable,
  loading,
  onRetry,
}: {
  href: string;
  error: unknown;
  isEditable?: boolean;
  loading: boolean;
  onRetry: () => void;
}) {
  if (!isEditable) {
    return (
      <LStack w="full" gap="1" userSelect="none">
        <LinkButton variant="subtle" href={href}>
          {href}
        </LinkButton>
        <Text variant="metadata">
          <WarningIcon w="3" display="inline" />
          &nbsp;<span>Link preview failed to load</span>
        </Text>
      </LStack>
    );
  }
  return (
    <styled.div
      position="absolute"
      inset="0"
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
      backgroundColor="status.danger.surface"
      borderRadius="lg"
      padding="2"
      height="min"
      gap="2"
      userSelect="none"
      contentEditable={false}
      role="alert"
      aria-live="polite"
    >
      <Text
        variant="supporting"
        color="status.danger.content"
        fontWeight="medium"
        maxW="prose"
      >
        Link preview failed: {deriveError(error)}
      </Text>
      <Button
        type="button"
        variant="subtle"
        onClick={onRetry}
        loading={loading}
      >
        Retry
      </Button>
    </styled.div>
  );
}
