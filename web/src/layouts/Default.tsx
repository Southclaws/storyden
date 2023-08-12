"use client";

import { Box, Flex } from "@chakra-ui/react";
import { usePathname } from "next/navigation";
import { PropsWithChildren } from "react";

import { Navpill } from "src/components/Navigation/Navpill/Navpill";
import { Sidebar } from "src/components/Navigation/Sidebar/Sidebar";
import { SIDEBAR_WIDTH } from "src/components/Navigation/useNavigation";

const ROUTES_WITHOUT_NAVPILL = ["/new"];

const isNavpillShown = (path: string | null) =>
  ROUTES_WITHOUT_NAVPILL.includes(path ?? "");

export function Default(props: PropsWithChildren) {
  const pathname = usePathname();
  return (
    <Flex
      width="full"
      flexDirection="row"
      bgColor="white"
      vaul-drawer-wrapper=""
    >
      {/* MOBILE */}
      <Box
        id="mobile-nav-container"
        display={{
          base: isNavpillShown(pathname) ? "none" : "unset",
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
        backgroundColor="white"
      >
        {props.children}
        <Box height="6rem"></Box>
      </Box>
    </Flex>
  );
}
