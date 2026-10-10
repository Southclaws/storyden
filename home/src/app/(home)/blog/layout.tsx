import { PropsWithChildren } from "react";

import { Box } from "@/styled-system/jsx";

export default function Layout({ children }: PropsWithChildren) {
  return <Box pb="2">{children}</Box>;
}
