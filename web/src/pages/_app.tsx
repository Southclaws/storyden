import { ChakraProvider, VStack } from "@chakra-ui/react";
import type { AppProps } from "next/app";
import { Navigation } from "src/components/Navigation/Navigation";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <ChakraProvider>
      <VStack>
        <Navigation />

        <Component {...pageProps} />
      </VStack>
    </ChakraProvider>
  );
}

export default MyApp;
