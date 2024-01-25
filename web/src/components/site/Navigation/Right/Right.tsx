import { Divider, VStack } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

export function Right() {
  return (
    <VStack className={Floating()} justify="space-between" px="4">
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
