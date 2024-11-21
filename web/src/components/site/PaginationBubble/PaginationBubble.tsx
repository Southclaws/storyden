import { parseAsInteger, useQueryState } from "nuqs";

import { Box, CardBox } from "@/styled-system/jsx";

import { PaginationControls } from "../PaginationControls/PaginationControls";

type Props = {
  path: string;
  totalPages: number;
  pageSize: number;
  onPageChange: (page: number) => void;
};

export function PaginationBubble({
  path,
  totalPages,
  pageSize,
  onPageChange,
}: Props) {
  const [page, setPage] = useQueryState("page", {
    ...parseAsInteger,
    defaultValue: 1,
    clearOnDefault: true,
  });

  function handlePage(page: number) {
    setPage(page);
    onPageChange(page);
  }

  const isFirst = page === 1;

  return (
    <Box
      style={{
        display: isFirst ? "none" : "block",
      }}
      position="fixed"
      bottom={{
        base: "24",
        md: "8",
      }}
    >
      <CardBox borderRadius="lg" p="1">
        <PaginationControls
          path={path}
          currentPage={page}
          pageSize={pageSize}
          totalPages={totalPages}
          onClick={handlePage}
        />
      </CardBox>
    </Box>
  );
}
