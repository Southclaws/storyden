import { Logo } from "@/components/Logo";
import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";
import type { BaseLayoutProps } from "fumadocs-ui/layouts/shared";
import {
  BracesIcon,
  HeartPlusIcon,
  LibraryBigIcon,
  MessageCircleHeartIcon,
} from "lucide-react";
import Image from "next/image";

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
      type: "menu",
      text: "Documentation",
      url: "/docs/introduction",
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
        {
          icon: <HeartPlusIcon />,
          text: "What is Storyden?",
          description:
            "Get to know your new favourite community + content platform.",
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
          description: "How does Storyden replace forums?",
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
          description:
            "An intelligent community knowledgebase, spend less time organising and more time sharing.",
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
      ],
    },
    {
      text: "Blog",
      url: "/blog",
      active: "nested-url",
    },
  ],
  githubUrl: "https://github.com/Southclaws/storyden",
};

export const baseOptionsZh: BaseLayoutProps = {
  ...baseOptions,
  links: [
    {
      type: "menu",
      text: "文档",
      url: "/zh/docs/introduction",
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
          url: "/zh/docs/introduction",
          text: "开始使用 Storyden",
          description: "几分钟内部署一个社区",
        },
        {
          icon: <HeartPlusIcon />,
          text: "Storyden 是什么？",
          description: "认识你的新社区与内容平台。",
          url: "/zh/docs/introduction/what-is-storyden",
          menu: {
            className: css({
              gridColumn: { md: "2" },
              gridRow: { md: "1 / span 1" },
            }),
          },
        },
        {
          icon: <MessageCircleHeartIcon />,
          text: "用于讨论",
          description: "Storyden 如何替代传统论坛？",
          url: "/zh/docs/introduction/discussion",
          menu: {
            className: css({
              gridColumn: { md: "2" },
              gridRow: { md: "2 / span 1" },
            }),
          },
        },
        {
          icon: <LibraryBigIcon />,
          text: "资料库",
          description: "智能社区知识库，把时间花在分享上，而不是整理上。",
          url: "/zh/docs/introduction/library",
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
          description: "浏览 API 参考并构建你的集成。",
          url: "/zh/docs/api",
          menu: {
            className: css({
              gridColumn: { md: "2", lg: "3" },
              gridRow: { md: "3 / span 1", lg: "2 / span 1" },
            }),
          },
        },
      ],
    },
    {
      text: "博客",
      url: "/zh/blog",
      active: "nested-url",
    },
  ],
};
