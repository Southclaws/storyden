import { PropsWithChildren } from "react";

import { Center } from "@/styled-system/jsx";

import { Empty } from "./Empty";

export function EmptyBox({ children }: PropsWithChildren) {
  return (
    <Center
      w="full"
      h="full"
      padding="8"
      borderRadius="lg"
      borderStyle="dashed"
      borderWidth="medium"
      borderColor="bg.subtle"
    >
      <Empty>{children}</Empty>
    </Center>
  );
}
