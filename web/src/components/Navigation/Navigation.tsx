import { Box, Flex, SlideFade, VStack } from "@chakra-ui/react";
import { Authenticated } from "./Authenticated";
import { Menu } from "./components/Menu/Menu";
import { Unauthenticated } from "./Unauthenticated";
import { useNavigation } from "./useNavigation";

const inactive = `hsla(180, 2%, 98%, 0.8)`;
const active = `hsla(220, 15%, 95%, 0.75)`;

export function Navigation() {
  const { isAuthenticated, isExpanded, onExpand, category } = useNavigation();

  return (
    <Box
      id="navigation-overlay"
      position="fixed"
      bottom="env(safe-area-inset-bottom)"
      width="100vw"
      height="100vh"
      pointerEvents="none"
      zIndex="overlay"
    >
      <Flex
        height="full"
        gap={3}
        p={2}
        justifyContent="end"
        alignItems="center"
        flexDir="column"
      >
        <Box
          maxW={{ base: "23em", md: "container.sm" }}
          width="full"
          minH="0"
          flex="0 0 1"
        >
          <SlideFade
            in={isExpanded}
            style={{
              maxHeight: "100%",
              display: "flex",
            }}
          >
            <VStack width="full">
              <Menu />
            </VStack>
          </SlideFade>
        </Box>

        <Flex
          p={{ base: 1, md: 2 }}
          borderRadius={{ base: 24, md: 28 }}
          backdropFilter="blur(1em)"
          transitionProperty="background-color"
          transitionDuration="0.5s"
          bgColor={isExpanded ? active : inactive}
          width="full"
          maxW={{ base: "23em", md: "container.sm" }}
          justifyContent="space-between"
          alignItems="center"
          pointerEvents="auto"
        >
          {isAuthenticated ? (
            <Authenticated onExpand={onExpand} category={category} />
          ) : (
            <Unauthenticated onExpand={onExpand} category={category} />
          )}
        </Flex>
      </Flex>
    </Box>
  );
}
