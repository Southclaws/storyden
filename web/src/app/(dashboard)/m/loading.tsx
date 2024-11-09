import { Spinner } from "@/components/ui/Spinner";
import { Center } from "@/styled-system/jsx";

export default async function Loading() {
  return (
    <Center id="authenticated-loading" w="full" height="96">
      <Spinner />
    </Center>
  );
}
