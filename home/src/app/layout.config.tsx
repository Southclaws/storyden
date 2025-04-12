import { Logo } from "@/components/Logo";
import { Box } from "@/styled-system/jsx";
import type { BaseLayoutProps } from "fumadocs-ui/layouts/shared";

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
      active: "nested-url",
    },
    {
      text: "Blog",
      url: "/blog",
      active: "nested-url",
    },
  ],
  githubUrl: "https://github.com/Southclaws/storyden",
};
