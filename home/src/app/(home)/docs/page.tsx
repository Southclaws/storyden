import { VStack } from "@/styled-system/jsx";
import Link from "next/link";

export default function Page() {
  return (
    <VStack pt="16">
      <h1>docs home</h1>
      <Link href="/docs/introduction">Get started</Link>
    </VStack>
  );
}
