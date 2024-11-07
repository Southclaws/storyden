import { Spinner } from "@/components/ui/Spinner";
import { Center } from "@/styled-system/jsx";

export function LoadingBanner() {
  return (
    <Center w="full" height="96">
      <Spinner />
    </Center>
  );
}
