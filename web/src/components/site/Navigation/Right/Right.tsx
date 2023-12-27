import { navbarStyles } from "../common";

import { Divider, VStack } from "@/styled-system/jsx";

export function Right() {
  return (
    <VStack className={navbarStyles} justify="space-between" px="4">
      <p>Threadbase</p>
      <p>Southclaws</p>

      <Divider />

      <p>sidebar stuff</p>
      <p>sidebar stuff</p>
      <p>sidebar stuff</p>
      <p>sidebar stuff</p>
      <p>sidebar stuff</p>
      <p>sidebar stuff</p>
      <p>sidebar stuff</p>
      <p>sidebar stuff</p>
      <p>sidebar stuff</p>
      <p>sidebar stuff</p>
    </VStack>
  );
}
