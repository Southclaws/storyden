import { ChakraProvider } from "@chakra-ui/react";
import { NextPage } from "next";
import type { AppProps } from "next/app";
import { ReactElement, ReactNode } from "react";
import { InfoProvider } from "src/api/InfoProvider/InfoProvider";
import { AuthProvider } from "src/auth/AuthProvider";
import { Default } from "src/layouts/Default";
import { extended } from "src/theme";

export type NextPageWithLayout<P = object, IP = P> = NextPage<P, IP> & {
  getLayout?: (page: ReactElement) => ReactNode;
};

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout;
};

function MyApp({ Component, pageProps }: AppPropsWithLayout) {
  const withLayout =
    Component.getLayout || ((page) => <Default>{page}</Default>);

  return (
    <ChakraProvider theme={extended}>
      <InfoProvider>
        <AuthProvider>
          <>{withLayout(<Component {...pageProps} />)}</>
        </AuthProvider>
      </InfoProvider>
    </ChakraProvider>
  );
}

export default MyApp;
