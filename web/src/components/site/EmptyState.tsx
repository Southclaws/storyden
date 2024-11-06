
import { Center } from "@/styled-system/jsx";
import { EmptyIcon } from "../ui/icons/Empty";

export function EmptyState() {
  return (
    <Center height="96" flexDirection="column" gap="2" color="fg.subtle">
      <EmptyIcon />
      <p>There&apos;s no content here.</p>
    </Center>
  );
}
