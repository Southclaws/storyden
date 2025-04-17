import { Box } from "@/styled-system/jsx";
import { PropsWithChildren } from "react";

export default function Layout({ children }: PropsWithChildren) {
  return <Box pb="2">{children}</Box>;
}
