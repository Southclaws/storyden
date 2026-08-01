import type { CSSProperties, ReactNode } from "react";

import { css } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";
import { token } from "@/styled-system/tokens";

type FoundationPageProps = {
  title: string;
  description: string;
  children: ReactNode;
};

export function FoundationPage({
  title,
  description,
  children,
}: FoundationPageProps) {
  return (
    <styled.main
      boxSizing="border-box"
      color="content.default"
      display="flex"
      flexDirection="column"
      gap="8"
      marginX="auto"
      maxW="7xl"
      padding={{ base: "4", md: "8" }}
      width="full"
    >
      <styled.header display="flex" flexDirection="column" gap="2" maxW="3xl">
        <styled.h1
          color="content.default"
          fontSize="2xl"
          fontWeight="semibold"
          lineHeight="tight"
        >
          {title}
        </styled.h1>
        <styled.p color="content.subtle" fontSize="md" lineHeight="relaxed">
          {description}
        </styled.p>
      </styled.header>
      {children}
    </styled.main>
  );
}

type FoundationSectionProps = {
  title: string;
  description?: string;
  children: ReactNode;
};

export function FoundationSection({
  title,
  description,
  children,
}: FoundationSectionProps) {
  return (
    <styled.section display="flex" flexDirection="column" gap="4" width="full">
      <styled.header display="flex" flexDirection="column" gap="1" maxW="3xl">
        <styled.h2
          color="content.default"
          fontSize="xl"
          fontWeight="semibold"
          lineHeight="tight"
        >
          {title}
        </styled.h2>
        {description ? (
          <styled.p color="content.subtle" fontSize="sm" lineHeight="relaxed">
            {description}
          </styled.p>
        ) : null}
      </styled.header>
      {children}
    </styled.section>
  );
}

type TokenCodeProps = {
  children: ReactNode;
};

export function TokenCode({ children }: TokenCodeProps) {
  return (
    <styled.code
      backgroundColor="surface.subtle"
      borderColor="border.subtle"
      borderRadius="sm"
      borderWidth="thin"
      color="content.default"
      fontFamily="mono"
      fontSize="xs"
      maxWidth="full"
      paddingX="1.5"
      paddingY="0.5"
      style={{ overflowWrap: "anywhere" }}
      whiteSpace="normal"
    >
      {children}
    </styled.code>
  );
}

function fallbackTokenVar(tokenPath: string) {
  return `var(--${tokenPath
    .replace(/([a-z])([A-Z])/g, "$1-$2")
    .replace(/\./g, "-")
    .toLowerCase()})`;
}

export function tokenValue(tokenPath: string) {
  return token(tokenPath as never, "");
}

export function tokenVar(tokenPath: string) {
  return token.var(tokenPath as never, fallbackTokenVar(tokenPath));
}

export function tokenStyle(
  customProperty: string,
  tokenPath: string,
): CSSProperties {
  return {
    [customProperty]: tokenVar(tokenPath),
  } as CSSProperties;
}

export const tokenGrid = css({
  display: "grid",
  gap: "3",
  gridTemplateColumns: "[repeat(auto-fit,minmax(13rem,1fr))]",
});

export const compactTokenGrid = css({
  display: "grid",
  gap: "3",
  gridTemplateColumns: "[repeat(auto-fit,minmax(14rem,1fr))]",
});
