import { Box, Flex } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { Navpill } from "src/components/Navigation/Navpill/Navpill";
import { Sidebar } from "src/components/Navigation/Sidebar/Sidebar";
import { SIDEBAR_WIDTH } from "src/components/Navigation/useNavigation";

export function Default(props: PropsWithChildren) {
  return (
    <Flex width="full" flexDirection="row">
      {/* MOBILE */}
      <Box
        id="mobile-nav-container"
        display={{
          base: "unset",
          md: "none",
        }}
      >
        <Navpill />
      </Box>

      {/* DESKTOP */}
      <Box
        display={{
          base: "none",
          md: "flex",
        }}
        minWidth={SIDEBAR_WIDTH}
        height="100vh"
      >
        <Sidebar />
      </Box>

      <Box
        as="main"
        width="full"
        maxW={{
          base: "full",
          lg: "container.md",
        }}
        px={4}
      >
        {props.children}
        <Box height="6rem"></Box>
      </Box>
    </Flex>
  );
}
