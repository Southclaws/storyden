import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import type { ReactNode } from "react";

import {
  SidebarNavigationLink,
  SidebarNavigationPane,
  SidebarNavigationSection,
} from "@/components/site/Navigation/SidebarNavigation";
import { AddIcon } from "@/components/ui/icons/Add";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { LinkIcon } from "@/components/ui/icons/Link";
import { RobotIcon } from "@/components/ui/icons/Robot";
import { SettingsIcon } from "@/components/ui/icons/Settings";
import { ToolIcon } from "@/components/ui/icons/Tool";
import { TrailIcon } from "@/components/ui/icons/Trail";
import { LinkButton } from "@/components/ui/link-button";
import { styled } from "@/styled-system/jsx";

const meta = {
  title: "Compositions/Navigation/Sidebar",
  parameters: {
    layout: "fullscreen",
    docs: {
      description: {
        story:
          "**Sidebar item row button pattern.** Use a compact, icon-only `action` at the trailing edge only for an operation scoped to that row's destination. It remains a separate, accessible control; use the `plain` button variant so it does not introduce a second hover surface. Keep the shared `4px` end inset used by navigation-tree rows. Product areas may own secondary destinations and recent-resource sections: the Robots example places Trails in its local navigation, then lists recent Trails below the Robots section instead of promoting Trails to the site navigation.",
      },
    },
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
            href="/l"
            icon={<LibraryIcon />}
            label="Settings documentation"
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
      <SidebarNavigationPane label="Robots navigation" title="Robots">
        <SidebarNavigationSection>
          <SidebarNavigationLink
            active
            href="/robots"
            icon={<RobotIcon />}
            label="All robots"
            action={
              <LinkButton
                aria-label="New robot"
                href="/robots/new"
                px="0"
                size="sm"
                variant="plain"
              >
                <AddIcon />
              </LinkButton>
            }
          />
          <SidebarNavigationLink
            href="/robots/chats"
            icon={<DiscussionIcon />}
            label="All chats"
            action={
              <LinkButton
                aria-label="New chat"
                href="/robots/chats/new"
                px="0"
                size="sm"
                variant="plain"
              >
                <AddIcon />
              </LinkButton>
            }
          />
          <SidebarNavigationLink
            href="/robots/trails"
            icon={<TrailIcon />}
            label="Trails"
            action={
              <LinkButton
                aria-label="New Trail"
                href="/robots/trails/new"
                px="0"
                size="sm"
                variant="plain"
              >
                <AddIcon />
              </LinkButton>
            }
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

        <SidebarNavigationSection label="Recent Trails">
          {[
            "Scheduled moderation check",
            "Monthly photo thread",
            "Weekly library review",
            "Daily report digest",
            "Quarterly welcome refresh",
          ].map((trail, index) => (
            <SidebarNavigationLink
              key={trail}
              href={`/robots/trails/${index}`}
              icon={<TrailIcon />}
              label={trail}
            />
          ))}
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
