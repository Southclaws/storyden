"use client";

import { PropsWithChildren } from "react";

import { Search } from "src/components/site/Navigation/Search/Search";

import { VStack } from "@/styled-system/jsx";

export default function Layout({ children }: PropsWithChildren) {
  return (
    <VStack height="full">
      <Search />

      {children}
    </VStack>
  );
}
