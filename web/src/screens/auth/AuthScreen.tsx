import { Grid, GridItem, Image, VStack } from "@chakra-ui/react";
import { AuthMethod } from "./AuthMethod/AuthMethod";
import { AuthSelection } from "./AuthSelection/AuthSelection";
import { AuthBox } from "./components/AuthBox";
import { Blue, Green } from "./components/Patterns";

type Props = {
  method?: string | undefined | null;
};
export function AuthScreen({ method }: Props) {
  return (
    <Grid
      templateRows="4fr 3fr 2fr auto 2fr 3fr 4fr"
      templateColumns="20vw 3fr 2fr auto 2fr 3fr 20vw"
      height="100vh"
      width="100vw"
      overflow="clip"
    >
      <GridItem
        gridRow="2/6"
        gridColumn="2/6"
        alignSelf="end"
        justifySelf="end"
      >
        <Blue />
      </GridItem>

      <GridItem
        gridRow="3/7"
        gridColumn="3/7"
        alignSelf="start"
        justifySelf="start"
      >
        <Green />
      </GridItem>

      <GridItem
        gridRow="3/4"
        gridColumn="4/5"
        alignSelf="center"
        justifySelf="center"
        p={6}
      >
        <Image height="50px" width="50px" src="/logo_200x200.png" alt="Logo" />
      </GridItem>

      <GridItem gridRow="4/5" gridColumn="4/5">
        <VStack
          width="full"
          height="full"
          justifyContent="center"
          flexDirection="column"
          alignItems="center"
          gap={4}
        >
          <AuthBox>
            {method === null ? (
              <AuthSelection />
            ) : (
              <AuthMethod method={method} />
            )}
          </AuthBox>
        </VStack>
      </GridItem>
    </Grid>
  );
}
