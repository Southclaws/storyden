import { Box, Flex, HStack, Image, Input, SlideFade } from "@chakra-ui/react";
import {
  Bell,
  Close,
  Create,
  Dashboard,
  Home,
  Login,
  Settings,
} from "src/components/Action/Action";
import { Menu } from "./components/Menu/Menu";
import { useNavpill } from "./useNavpill";
import { ProfileReference } from "src/components/ProfileReference/ProfileReference";

export function Navpill() {
  const { overlayRef, isExpanded, onExpand, account } = useNavpill();
  return (
    <Box
      id="navpill-overlay"
      position="fixed"
      bottom="env(safe-area-inset-bottom)"
      pb={2}
      width="100vw"
      height="100vh"
      pointerEvents="none"
      zIndex="overlay"
    >
      <Flex
        id="flex-container"
        height="full"
        gap={3}
        p="min(3vh, 1em)"
        justifyContent="end"
        alignItems="center"
        flexDir="column"
        ref={overlayRef}
      >
        <Flex
          pl={{ base: isExpanded ? 2 : 2 }}
          pr={{ base: isExpanded ? 2 : 4 }}
          py={1}
          gap={2}
          flexDirection="column"
          borderRadius={20}
          backdropFilter="blur(4px)"
          transitionProperty="background-color"
          transitionDuration="0.5s"
          backgroundColor="hsla(210, 38.5%, 94.9%, 0.8)"
          border="1px solid hsla(209, 100%, 20%, 0.02)"
          width="full"
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

          {account ? (
            <HStack gap={4} w="full" justifyContent="space-between">
              {isExpanded ? (
                <>
                  <ProfileReference
                    handle={account.handle}
                    showHandle={false}
                  />

                  <Input
                    variant="outline"
                    border="none"
                    size="sm"
                    placeholder="Search anything..."
                  />
                  <Close onClick={onExpand} />
                </>
              ) : (
                <>
                  <ProfileReference
                    handle={account.handle}
                    showHandle={false}
                  />
                  <Home />
                  <Create />
                  <Bell />
                  <Dashboard onClick={onExpand} />
                </>
              )}
            </HStack>
          ) : (
            <HStack gap={4} w="full" justifyContent="space-between">
              <Login />

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
          )}
        </Flex>
      </Flex>
    </Box>
  );
}
