import { range } from "lodash";
import { parseAsInteger, useQueryState } from "nuqs";
import { MouseEvent, useState } from "react";

import { LinkButton } from "@/components/ui/link-button";
import { HStack, styled } from "@/styled-system/jsx";

const MAX_PAGES_SHOWN = 5;
const MID_POINT = MAX_PAGES_SHOWN / 2;

export type Props = {
  path: string;
  params?: Record<string, string>;
  currentPage: number;
  totalPages: number;
  pageSize: number;
  onClick?: (page: number) => void;
};

export function PaginationControls({
  path,
  params,
  currentPage,
  totalPages,
  onClick,
}: Props) {
  if (totalPages <= 1) {
    return null;
  }

  // NOTE: pages are 1-indexed
  const lastPage = totalPages;
  const lastPageIndex = totalPages + 1;
  const allPages = range(1, lastPageIndex);

  const tooMany = allPages.length > MAX_PAGES_SHOWN;
  const nearEnd = lastPage - currentPage < MID_POINT;
  const nearStart = currentPage < MID_POINT;

  const needStartJump = tooMany && currentPage > MID_POINT + 1;
  const needEndJump = tooMany && currentPage < lastPage - MID_POINT;

  const getPages = () => {
    if (!tooMany) {
      return allPages;
    }

    if (nearEnd) {
      return range(lastPageIndex - MAX_PAGES_SHOWN, lastPageIndex);
    }

    if (nearStart) {
      return range(1, MAX_PAGES_SHOWN + 1);
    }

    const start = currentPage - 3;
    const end = currentPage + 2;

    return allPages.slice(start, end);
  };

  const targetPages = getPages();

  const clickHandler =
    (page: number) => (event: MouseEvent<HTMLAnchorElement>) => {
      if (onClick) {
        event.preventDefault();
        onClick(page);
      }
    };

  return (
    <HStack w="min" p="1">
      {needStartJump && (
        <>
          <LinkButton
            variant="ghost"
            size="xs"
            href={`${path}?${new URLSearchParams({
              ...params,
              page: "1",
            }).toString()}`}
            onClick={clickHandler(1)}
          >
            {1}
          </LinkButton>
          <Sep />
        </>
      )}

      {targetPages.map((v) => {
        const pageName = v.toString();
        const withPage = new URLSearchParams({
          ...params,
          page: pageName,
        });

        return (
          <LinkButton
            variant="ghost"
            backgroundColor={v === currentPage ? "bg.muted" : undefined}
            size="xs"
            key={v}
            href={`${path}?${withPage.toString()}`}
            onClick={clickHandler(v)}
          >
            {pageName}
          </LinkButton>
        );
      })}

      {needEndJump && (
        <>
          <Sep />
          <LinkButton
            variant="ghost"
            size="xs"
            href={`${path}?${new URLSearchParams({
              ...params,
              page: lastPage.toString(),
            }).toString()}`}
            onClick={clickHandler(lastPage)}
          >
            {lastPage}
          </LinkButton>
        </>
      )}
    </HStack>
  );
}

function Sep() {
  return (
    <styled.span fontSize="xs" color="fg.disabled">
      •••
    </styled.span>
  );
}
