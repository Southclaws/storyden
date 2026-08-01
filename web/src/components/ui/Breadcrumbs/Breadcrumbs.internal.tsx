import { ChevronRightIcon } from "@heroicons/react/24/outline";
import { Fragment, PropsWithChildren } from "react";

import { LinkButton } from "@/components/ui/link-button";
import { Box, HStack, styled } from "@/styled-system/jsx";

type Props = {
  index: Breadcrumb;
  crumbs: Breadcrumb[];
};

type Breadcrumb = {
  label: string;
  href: string;
};

export function Breadcrumbs({
  index,
  crumbs,
  children,
}: PropsWithChildren<Props>) {
  return (
    <HStack
      w="full"
      gap="1"
      justify="space-between"
      // Align top because generally we just want same-height children and the
      // HStack below is using an offset at the bottom for the scroll gutter.
      alignItems="start"
    >
      <HStack
        w="full"
        color="fg.subtle"
        overflowX="scroll"
        pt="scrollGutter"
        mt="-scrollGutter"
      >
        <LinkButton
          size="xs"
          variant="subtle"
          flexShrink="0"
          minW="min"
          href={index.href}
        >
          {index.label}
        </LinkButton>
        {crumbs.map((c) => {
          return (
            <Fragment key={c.href}>
              <Box flexShrink="0">
                <ChevronRightIcon width="1rem" />
              </Box>

              <BreadcrumbButton crumb={c} />
            </Fragment>
          );
        })}
      </HStack>

      {children}
    </HStack>
  );
}

function BreadcrumbButton({ crumb }: { crumb: Breadcrumb }) {
  return (
    <LinkButton
      size="xs"
      variant="subtle"
      flexShrink="0"
      maxW="64"
      overflow="hidden"
      href={crumb.href}
    >
      <styled.span overflowX="hidden" textWrap="nowrap" textOverflow="ellipsis">
        {crumb.label}
      </styled.span>
    </LinkButton>
  );
}
