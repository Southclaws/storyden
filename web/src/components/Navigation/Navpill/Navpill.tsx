import { Box, Flex, HStack, Input, SlideFade } from "@chakra-ui/react";
import {
  Bell,
  Close,
  Create,
  Dashboard,
  Home,
  Login,
  Settings,
} from "src/components/Action/Action";
import { useNavpill } from "./useNavpill";
import { Menu } from "./components/Menu/Menu";

export function Navpill() {
  const { overlayRef, isExpanded, onExpand, isAuthenticated } = useNavpill();
  return (
    <Box
      id="navpill-overlay"
      position="fixed"
      bottom="env(safe-area-inset-bottom)"
      width="100vw"
      height="100vh"
      pointerEvents="none"
      zIndex="overlay"
    >
      <Flex
        id="flex-container"
        height="full"
        gap={3}
        p={4}
        justifyContent="end"
        alignItems="center"
        flexDir="column"
        ref={overlayRef}
      >
        <Flex
          px={{ base: isExpanded ? 2 : 4 }}
          py={{ base: 2 }}
          gap={2}
          flexDirection="column"
          borderRadius={25}
          backdropFilter="blur(4px)"
          transitionProperty="background-color"
          transitionDuration="0.5s"
          backgroundColor="hsla(210, 38.5%, 94.9%, 0.8)"
          border="1px solid hsla(209, 100%, 20%, 0.02)"
          width={isExpanded ? "100%" : "min-content"}
          maxW={{ base: "23em", md: "container.sm" }}
          justifyContent="space-between"
          alignItems="center"
          pointerEvents="auto"
        >
          <SlideFade
            in={isExpanded}
            style={{
              maxHeight: "100%",
              width: "100%",
              display: isExpanded ? "flex" : "none",
              flexDirection: "column",
            }}
          >
            <HStack justify="space-between">
              <Home />
              <Settings />
            </HStack>

            <Menu />
          </SlideFade>

          <HStack gap={4} w="full" justifyContent="space-between">
            {isAuthenticated ? <Bell /> : <Login />}

            {isExpanded ? (
              <>
                <Input
                  variant="outline"
                  border="none"
                  placeholder="Search anything..."
                />
                <Close onClick={onExpand} />
              </>
            ) : (
              <>
                <Create />
                <Dashboard onClick={onExpand} />
              </>
            )}
          </HStack>
        </Flex>
      </Flex>
    </Box>
  );
}
