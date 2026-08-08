import { Spinner } from "@/components/ui/spinner";
import { Center } from "@/styled-system/jsx";

export function LoadingBanner() {
  return (
    <Center w="full" height="96">
      <Spinner />
    </Center>
  );
}
