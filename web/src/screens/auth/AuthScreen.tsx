import { Flex } from "@chakra-ui/react";
import { AuthMethod } from "./AuthMethod/AuthMethod";
import { AuthSelection } from "./AuthSelection/AuthSelection";
import { AuthBox } from "./components/AuthBox";

type Props = {
  method?: string | undefined | null;
};
export function AuthScreen({ method }: Props) {
  return (
    <Flex
      height="100vh"
      width="full"
      justifyContent="center"
      flexDirection="column"
      alignItems="center"
      backgroundImage="/blobs.png"
      backgroundPosition="center"
      backgroundSize="cover"
      gap={4}
    >
      <AuthBox>
        {method === null ? <AuthSelection /> : <AuthMethod method={method} />}
      </AuthBox>
    </Flex>
  );
}
