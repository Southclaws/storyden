"use client";

import { usePathname } from "next/navigation";

import { Navpill } from "src/components/Navigation/Navpill/Navpill";
import { Sidebar } from "src/components/Navigation/Sidebar/Sidebar";

import { Box } from "@/styled-system/jsx";

const ROUTES_WITHOUT_NAVPILL = ["/new"];

const isNavpillShown = (path: string | null) =>
  ROUTES_WITHOUT_NAVPILL.includes(path ?? "");

export function Navigation() {
  const pathname = usePathname();

  return (
    <>
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
        id="desktop-nav-container"
        display={{
          base: "none",
          md: "flex",
        }}
        height="100vh"
        // The sidebar width is identical in both this container and the sidebar
        // itself. The reason for this is the sidebar is "position: fixed" which
        // means it cannot inherit the width from a parent since its true parent
        // is the viewport, and to get around this, the default layout positions
        // an empty box to the left of the viewport in order to push the content
        // right and then the actual sidebar is rendered on top of this with the
        // same width sizing configuration.
        minWidth={{
          md: "1/4",
          lg: "1/3",
        }}
      >
        <Sidebar />
      </Box>
    </>
  );
}
