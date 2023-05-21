import { Box, Flex } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { Navpill } from "src/components/Navigation/Navpill/Navpill";
import { Sidebar } from "src/components/Navigation/Sidebar/Sidebar";

export function Default(props: PropsWithChildren) {
  return (
    <Flex
      width="full"
      minHeight="100vh"
      alignItems="stretch"
      flexDirection="row"
    >
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
        minWidth={{
          md: "25%",
          lg: "33%",
        }}
        height="100vh"
        px={4}
        bgColor="blackAlpha.50"
      >
        <Sidebar />
      </Box>

      <Box
        as="main"
        width="full"
        height="100vh"
        maxW={{
          base: "full",
          lg: "container.md",
        }}
        overflowY="scroll"
        px={4}
      >
        {props.children}
      </Box>
    </Flex>
  );
}
