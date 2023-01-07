import { Flex } from "@chakra-ui/react";
import { StorydenLogo } from "src/components/StorydenLogo";
import { AuthMethod } from "./AuthMethod/AuthMethod";
import { AuthSelection } from "./AuthSelection/AuthSelection";
import { AuthBox } from "./components/AuthBox";

type Props = {
  method?: string | undefined | null;
};
export function AuthScreen({ method }: Props) {
  return (
    <Flex
      id="AuthScreen"
      height="100vh"
      width="full"
      justifyContent="center"
      flexDirection="column"
      alignItems="center"
      backgroundImage="/blobs.png"
      backgroundPosition="center"
      backgroundSize="cover"
      gap={4}
      padding={6}
    >
      <StorydenLogo height="3em" />
      {method === null ? ( //
        <AuthBox bg="linear-gradient(141.91deg, #B7CEF1 0%, #2FD596 99.55%)">
          <AuthSelection />
        </AuthBox>
      ) : (
        <AuthBox
          bg="light"
          boxShadow="0 10px 30px rgba(0, 0, 0, 0.05)"
          borderRadius="md"
        >
          <AuthMethod method={method} />
        </AuthBox>
      )}
    </Flex>
  );
}
