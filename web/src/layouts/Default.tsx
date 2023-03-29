import { Box, Flex } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { Navpill } from "src/components/Navpill/Navpill";
import { Sidebar } from "src/components/Sidebar/Sidebar";

export function Default(props: PropsWithChildren) {
  return (
    <Flex
      width="full"
      height="full"
      minHeight="100vh"
      alignItems="stretch"
      flexDirection="row"
    >
      <Box
        visibility={{
          base: "unset",
          md: "collapse",
        }}
      >
        <Navpill />
      </Box>

      <Flex
        as="header"
        width={{ md: "33%", lg: "40%" }}
        bgColor="blackAlpha.50"
        px={4}
        visibility={{
          base: "collapse",
          md: "unset",
        }}
        justifyContent="end"
      >
        <Box width="xs">
          <Sidebar />
        </Box>
      </Flex>

      <Flex as="main" px={4}>
        <Box
          maxW={{
            base: "full",
            md: "container.md",
          }}
        >
          {props.children}
        </Box>
      </Flex>
    </Flex>
  );
}
