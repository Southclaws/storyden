"use client";

import { usePathname } from "next/navigation";

import { Navpill } from "src/components/Navigation/Navpill/Navpill";
import { Sidebar } from "src/components/Navigation/Sidebar/Sidebar";

import { css } from "@/styled-system/css";

import { SIDEBAR_WIDTH } from "./useNavigation";

const ROUTES_WITHOUT_NAVPILL = ["/new"];

const isNavpillShown = (path: string | null) =>
  ROUTES_WITHOUT_NAVPILL.includes(path ?? "");

export function Navigation() {
  const pathname = usePathname();

  return (
    <>
      {/* MOBILE */}
      <div
        id="mobile-nav-container"
        className={css({
          display: {
            base: isNavpillShown(pathname) ? "none" : "unset",
            md: "none",
          },
        })}
      >
        <Navpill />
      </div>

      {/* DESKTOP */}
      <div
        id="desktop-nav-container"
        className={css({
          display: {
            base: "none",
            md: "flex",
          },
          minWidth: SIDEBAR_WIDTH,
          height: "100vh",
        })}
      >
        <Sidebar />
      </div>
    </>
  );
}
