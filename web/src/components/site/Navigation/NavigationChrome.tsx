"use client";

import { usePathname } from "next/navigation";
import type { KeyboardEvent, MouseEvent, ReactNode } from "react";
import { useCallback, useEffect, useRef, useState } from "react";

import * as Tooltip from "@/components/ui/tooltip";

import { HomeAnchor } from "./Anchors/Home";
import { SearchAnchor } from "./Anchors/Search";
import { MobileCommandBar } from "./MobileCommandBar/MobileCommandBar";

type Props = {
  desktopSidebar: ReactNode;
  siteNavigation: ReactNode;
  layoutEditMode: ReactNode;
  sidebarBottom: ReactNode;
};

const focusableSelector = [
  "a[href]",
  "button:not([disabled])",
  "textarea:not([disabled])",
  "input:not([disabled])",
  "select:not([disabled])",
  '[tabindex]:not([tabindex="-1"])',
].join(",");

const mobileMediaQuery = "(max-width: 767px)";

export function NavigationChrome({
  desktopSidebar,
  siteNavigation,
  layoutEditMode,
  sidebarBottom,
}: Props) {
  const pathname = usePathname();
  const chromeRef = useRef<HTMLDivElement>(null);
  const drawerRef = useRef<HTMLDivElement>(null);
  const previousFocusRef = useRef<HTMLElement | null>(null);
  const wasOpenRef = useRef(false);
  const [isMobileNavigationOpen, setMobileNavigationOpen] = useState(false);
  const [isMobileViewport, setMobileViewport] = useState(false);

  const isMobileNavigationModal = isMobileNavigationOpen && isMobileViewport;

  const openMobileNavigation = useCallback(() => {
    if (document.activeElement instanceof HTMLElement) {
      previousFocusRef.current = document.activeElement;
    }

    setMobileNavigationOpen(true);
  }, []);

  const closeMobileNavigation = useCallback(() => {
    setMobileNavigationOpen(false);
  }, []);

  useEffect(() => {
    const media = window.matchMedia(mobileMediaQuery);

    function handleMediaChange(event: MediaQueryList | MediaQueryListEvent) {
      setMobileViewport(event.matches);

      if (!event.matches) {
        setMobileNavigationOpen(false);
      }
    }

    handleMediaChange(media);
    media.addEventListener("change", handleMediaChange);

    return () => {
      media.removeEventListener("change", handleMediaChange);
    };
  }, []);

  useEffect(() => {
    setMobileNavigationOpen(false);
  }, [pathname]);

  useEffect(() => {
    if (!isMobileNavigationModal) {
      return;
    }

    const listener = (event: globalThis.KeyboardEvent) => {
      if (event.key === "Escape") {
        setMobileNavigationOpen(false);
      }
    };

    document.addEventListener("keydown", listener);

    return () => {
      document.removeEventListener("keydown", listener);
    };
  }, [isMobileNavigationModal]);

  useEffect(() => {
    const page = document.getElementById("navigation__scroll");

    if (!isMobileNavigationModal) {
      page?.removeAttribute("aria-hidden");
      page?.removeAttribute("inert");

      if (wasOpenRef.current) {
        const previousFocus = previousFocusRef.current;
        const returnTarget = previousFocus?.isConnected
          ? previousFocus
          : chromeRef.current?.querySelector<HTMLElement>(
              '[aria-controls="navigation__leftbar"]',
            );

        returnTarget?.focus({ preventScroll: true });
        previousFocusRef.current = null;
        wasOpenRef.current = false;
      }

      return;
    }

    page?.setAttribute("aria-hidden", "true");
    page?.setAttribute("inert", "");
    wasOpenRef.current = true;

    const frame = window.requestAnimationFrame(() => {
      const drawer = drawerRef.current;
      const focusTarget = drawer?.querySelector<HTMLElement>(focusableSelector);

      (focusTarget ?? drawer)?.focus({ preventScroll: true });
    });

    return () => {
      window.cancelAnimationFrame(frame);
      page?.removeAttribute("aria-hidden");
      page?.removeAttribute("inert");
    };
  }, [isMobileNavigationModal]);

  function handleTrapFocus(event: KeyboardEvent<HTMLDivElement>) {
    if (!isMobileNavigationModal || event.key !== "Tab") {
      return;
    }

    const chrome = chromeRef.current;
    if (!chrome) {
      return;
    }

    const focusableElements = Array.from(
      chrome.querySelectorAll<HTMLElement>(focusableSelector),
    ).filter((element) => {
      const style = window.getComputedStyle(element);

      return style.display !== "none" && style.visibility !== "hidden";
    });

    const first = focusableElements.at(0);
    const last = focusableElements.at(-1);

    if (!first || !last) {
      event.preventDefault();
      drawerRef.current?.focus({ preventScroll: true });
      return;
    }

    if (event.shiftKey && document.activeElement === first) {
      event.preventDefault();
      last.focus({ preventScroll: true });
    }

    if (!event.shiftKey && document.activeElement === last) {
      event.preventDefault();
      first.focus({ preventScroll: true });
    }
  }

  function handleDrawerClick(event: MouseEvent<HTMLDivElement>) {
    if (
      isMobileNavigationModal &&
      event.target instanceof Element &&
      event.target.closest("a[href]")
    ) {
      closeMobileNavigation();
    }
  }

  return (
    <div
      ref={chromeRef}
      className="navigation__grid navigation__fixed-grid"
      data-mobile-navigation-open={isMobileNavigationOpen ? "true" : "false"}
      data-mobile-navigation-modal={isMobileNavigationModal ? "true" : "false"}
      onKeyDown={handleTrapFocus}
    >
      <div
        aria-hidden="true"
        className="navigation__mobile-scrim"
        onClick={closeMobileNavigation}
      />

      <div
        id="navigation__leftbar"
        ref={drawerRef}
        aria-label={isMobileNavigationModal ? "Site navigation" : undefined}
        aria-modal={isMobileNavigationModal ? "true" : undefined}
        className="navigation__leftbar"
        onClick={handleDrawerClick}
        role={isMobileNavigationModal ? "dialog" : undefined}
        tabIndex={-1}
        data-sd-region="sidebar"
      >
        <div className="navigation__sidebar-top">
          <div
            aria-label="Primary navigation"
            className="navigation-pane__top-actions"
            role="toolbar"
          >
            <div className="navigation-pane__top-actions-start">
              <Tooltip.Root openDelay={300}>
                <Tooltip.Trigger asChild>
                  <span className="navigation-pane__tooltip-anchor">
                    <HomeAnchor hideLabel size="sm" />
                  </span>
                </Tooltip.Trigger>
                <Tooltip.Positioner>
                  <Tooltip.Content>Home</Tooltip.Content>
                </Tooltip.Positioner>
              </Tooltip.Root>

              <Tooltip.Root openDelay={300}>
                <Tooltip.Trigger asChild>
                  <span className="navigation-pane__tooltip-anchor">
                    <SearchAnchor hideLabel size="sm" />
                  </span>
                </Tooltip.Trigger>
                <Tooltip.Positioner>
                  <Tooltip.Content>Search</Tooltip.Content>
                </Tooltip.Positioner>
              </Tooltip.Root>
            </div>

            {layoutEditMode}
          </div>
        </div>
        <div className="navigation__desktop-sidebar">{desktopSidebar}</div>
        <div className="navigation__site-sidebar">{siteNavigation}</div>
        <div className="navigation__sidebar-bottom">{sidebarBottom}</div>
      </div>

      <div className="navigation__mobile-topbar" data-sd-region="topbar">
        <MobileCommandBar
          isOpen={isMobileNavigationOpen}
          onClose={closeMobileNavigation}
          onOpen={openMobileNavigation}
        />
      </div>
    </div>
  );
}
