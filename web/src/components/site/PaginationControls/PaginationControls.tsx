import { range } from "lodash";

import { Link } from "src/theme/components/Link";

import { HStack, styled } from "@/styled-system/jsx";

const MAX_PAGES_SHOWN = 5;
const MID_POINT = MAX_PAGES_SHOWN / 2;

export type Props = {
  path: string;
  params?: Record<string, string>;
  currentPage: number;
  totalPages: number;
  pageSize: number;
  onClick: (page: number) => void;
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

  return (
    <HStack w="min" p="1">
      {needStartJump && (
        <>
          <Link
            kind="ghost"
            size="xs"
            href={`${path}?${new URLSearchParams({
              ...params,
              page: "1",
            }).toString()}`}
            onClick={() => onClick(1)}
          >
            {1}
          </Link>
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
          <Link
            kind="ghost"
            backgroundColor={v === currentPage ? "bg.muted" : undefined}
            size="xs"
            key={v}
            href={`${path}?${withPage.toString()}`}
            onClick={() => onClick(v)}
          >
            {pageName}
          </Link>
        );
      })}

      {needEndJump && (
        <>
          <Sep />
          <Link
            kind="ghost"
            size="xs"
            href={`${path}?${new URLSearchParams({
              ...params,
              page: lastPage.toString(),
            }).toString()}`}
            onClick={() => onClick(lastPage)}
          >
            {lastPage}
          </Link>
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
