import { Box, Flex } from "@chakra-ui/react";
import { Navpill } from "./Navpill/Navpill";
import { Sidebar } from "./Sidebar/Sidebar";

// Navigation displays either the sidebar on desktop or the navpill on mobile.
export function Navigation() {
  return (
    <>
      {/* MOBILE */}
      <Box
        display={{
          base: "unset",
          md: "none",
        }}
      >
        <Navpill />
      </Box>

      {/* DESKTOP */}
      <Flex
        display={{
          base: "none",
          md: "flex",
        }}
        border="1px solid red"
        as="header"
        minWidth={{
          md: "25%",
          lg: "33%",
        }}
        px={4}
        justifyContent="end"
        bgColor="blackAlpha.50"
      >
        <Box
          border="1px solid blue"
          maxWidth="2xs"
          //
        >
          <Sidebar />
        </Box>
      </Flex>
    </>
  );
}
