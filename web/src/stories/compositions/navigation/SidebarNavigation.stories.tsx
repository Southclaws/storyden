import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import type { ReactNode } from "react";

import {
  SidebarNavigationLink,
  SidebarNavigationPane,
  SidebarNavigationSection,
} from "@/components/site/Navigation/SidebarNavigation";
import { AddIcon } from "@/components/ui/icons/Add";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { HomeIcon } from "@/components/ui/icons/Home";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { LinkIcon } from "@/components/ui/icons/Link";
import { RobotIcon } from "@/components/ui/icons/Robot";
import { SettingsIcon } from "@/components/ui/icons/Settings";
import { ToolIcon } from "@/components/ui/icons/Tool";
import { styled } from "@/styled-system/jsx";

const meta = {
  title: "Compositions/Navigation/Sidebar",
  parameters: {
    layout: "fullscreen",
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

function SidebarPreview({ children }: { children: ReactNode }) {
  return (
    <styled.main
      backgroundColor="background.canvas"
      boxSizing="border-box"
      display="grid"
      minHeight="dvh"
      padding={{ base: "2", md: "6" }}
      placeItems="start"
    >
      <styled.div
        style={{
          height: "42rem",
          maxHeight: "calc(100dvh - 3rem)",
          width: "18rem",
        }}
      >
        {children}
      </styled.div>
    </styled.main>
  );
}

export const SiteNavigation: Story = {
  render: () => (
    <SidebarPreview>
      <SidebarNavigationPane label="Site navigation" title="Storyden">
        <SidebarNavigationSection>
          <SidebarNavigationLink
            active
            href="/"
            icon={<DiscussionIcon />}
            label="Discussion"
          />
          <SidebarNavigationLink
            href="/l"
            icon={<LibraryIcon />}
            label="Library"
          />
        </SidebarNavigationSection>

        <SidebarNavigationSection label="Community">
          <SidebarNavigationLink
            href="/links"
            icon={<LinkIcon />}
            label="Links"
          />
          <SidebarNavigationLink
            href="/settings"
            icon={<SettingsIcon />}
            label="Settings"
          />
        </SidebarNavigationSection>
      </SidebarNavigationPane>
    </SidebarPreview>
  ),
};

export const SettingsNavigation: Story = {
  render: () => (
    <SidebarPreview>
      <SidebarNavigationPane
        label="Settings navigation"
        title="Settings"
        footer={
          <SidebarNavigationLink
            href="/"
            icon={<HomeIcon />}
            label="Back to community"
          />
        }
      >
        <SidebarNavigationSection>
          <SidebarNavigationLink
            active
            description="Theme and display preferences"
            href="/settings/interface"
            icon={<SettingsIcon />}
            label="Interface"
          />
          <SidebarNavigationLink
            description="Password and sign-in methods"
            href="/settings/authentication"
            icon={<SettingsIcon />}
            label="Authentication"
          />
          <SidebarNavigationLink
            description="Personal API credentials"
            href="/settings/access-keys"
            icon={<ToolIcon />}
            label="Access keys"
          />
          <SidebarNavigationLink
            description="Connected application access"
            href="/settings/oauth"
            icon={<LinkIcon />}
            label="OAuth"
          />
        </SidebarNavigationSection>
      </SidebarNavigationPane>
    </SidebarPreview>
  ),
};

export const RobotsNavigation: Story = {
  render: () => (
    <SidebarPreview>
      <SidebarNavigationPane
        label="Robots navigation"
        title="Robots"
        footer={
          <SidebarNavigationLink
            href="/"
            icon={<HomeIcon />}
            label="Back to community"
          />
        }
      >
        <SidebarNavigationSection>
          <SidebarNavigationLink
            active
            href="/robots"
            icon={<RobotIcon />}
            label="All robots"
          />
          <SidebarNavigationLink
            href="/robots/new"
            icon={<AddIcon />}
            label="New robot"
          />
          <SidebarNavigationLink
            href="/robots/chats"
            icon={<DiscussionIcon />}
            label="All chats"
          />
          <SidebarNavigationLink
            href="/robots/chats/new"
            icon={<AddIcon />}
            label="New chat"
          />
        </SidebarNavigationSection>

        <SidebarNavigationSection label="Robots">
          <SidebarNavigationLink
            href="/robots/moderator"
            icon={<RobotIcon />}
            label="Moderator"
          />
          <SidebarNavigationLink
            href="/robots/librarian"
            icon={<RobotIcon />}
            label="Librarian"
          />
        </SidebarNavigationSection>

        <SidebarNavigationSection label="Recent chats">
          {Array.from({ length: 10 }, (_, index) => (
            <SidebarNavigationLink
              key={index}
              href={`/robots/chats/${index}`}
              icon={<DiscussionIcon />}
              label={`Community planning ${index + 1}`}
            />
          ))}
        </SidebarNavigationSection>
      </SidebarNavigationPane>
    </SidebarPreview>
  ),
};
