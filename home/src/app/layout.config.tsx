import type { BaseLayoutProps, MenuItemType } from "fumadocs-ui/layouts/shared";
import {
  BracesIcon,
  HeartPlusIcon,
  LibraryBigIcon,
  MessageCircleHeartIcon,
} from "lucide-react";
import Image from "next/image";

import { Logo } from "@/components/Logo";
import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";

import { GitHub } from "./(home)/Footer";

export const githubLink = {
  type: "icon",
  text: "GitHub",
  label: "GitHub",
  url: "https://github.com/Southclaws/storyden",
  icon: <GitHub aria-hidden="true" />,
} satisfies NonNullable<BaseLayoutProps["links"]>[number];

const documentationQuickLinks = [
  {
    icon: <HeartPlusIcon />,
    text: "What is Storyden?",
    description: "Get to know your new favourite community + content platform.",
    url: "/docs/introduction/what-is-storyden",
    menu: {
      className: css({
        gridColumn: { md: "2" },
        gridRow: { md: "1 / span 1" },
      }),
    },
  },
  {
    icon: <MessageCircleHeartIcon />,
    text: "For discussion",
    description: "Conversations worth coming back to.",
    url: "/docs/introduction/discussion",
    menu: {
      className: css({
        gridColumn: { md: "2" },
        gridRow: { md: "2 / span 1" },
      }),
    },
  },
  {
    icon: <LibraryBigIcon />,
    text: "Library",
    description: "A shared home for the things your community knows.",
    url: "/docs/introduction/library",
    menu: {
      className: css({
        gridColumn: { md: "1", lg: "3" },
        gridRow: { md: "3 / span 1", lg: "1 / span 1" },
      }),
    },
  },
  {
    icon: <BracesIcon />,
    text: "API",
    description: "Browse the API reference and build something awesome.",
    url: "/docs/api",
    menu: {
      className: css({
        gridColumn: { md: "2", lg: "3" },
        gridRow: { md: "3 / span 1", lg: "2 / span 1" },
      }),
    },
  },
] satisfies MenuItemType["items"];

export const baseOptions: BaseLayoutProps = {
  themeSwitch: {
    enabled: false,
  },
  nav: {
    title: (
      <>
        <Box w="8">
          <Logo />
        </Box>
        Storyden
      </>
    ),
  },
  links: [
    {
      text: "Documentation",
      url: "/docs/introduction",
      on: "menu",
    },
    ...documentationQuickLinks.map((item) => ({
      ...item,
      on: "menu" as const,
    })),
    {
      type: "menu",
      text: "Documentation",
      on: "nav",
      items: [
        {
          active: "nested-url",
          menu: {
            className: css({
              gridColumn: { md: "1" },
              gridRow: { md: "1 / span 2" },
            }),
            banner: (
              <Box
                className={css({
                  margin: "-3",
                })}
              >
                <Image
                  className={css({
                    borderRadius: "md",
                    boxShadow: "lg",
                    maxWidth: "100%",
                    objectFit: "cover",
                  })}
                  style={{
                    maskImage:
                      "linear-gradient(to bottom,white 50%,transparent)",
                  }}
                  width="1200"
                  height="630"
                  src="/docs_get_started_banner.png"
                  alt=""
                />
              </Box>
            ),
          },
          url: "/docs/introduction",
          text: "Get started with Storyden",
          description: "Deploy a community in minutes",
        },
        ...documentationQuickLinks,
      ],
    },
    {
      text: "Blog",
      url: "/blog",
      active: "nested-url",
    },
    githubLink,
  ],
};
