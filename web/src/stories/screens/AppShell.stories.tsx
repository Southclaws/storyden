import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import {
  SidebarNavigationLink,
  SidebarNavigationPane,
  SidebarNavigationSection,
} from "@/components/site/Navigation/SidebarNavigation";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { FormControl } from "@/components/ui/form-control";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { LinkIcon } from "@/components/ui/icons/Link";
import { SettingsIcon } from "@/components/ui/icons/Settings";
import { ToolIcon } from "@/components/ui/icons/Tool";
import { Input } from "@/components/ui/input";
import { PageHeader } from "@/components/ui/page-header";
import { SectionHeading } from "@/components/ui/section-heading";
import { CardRows } from "@/components/ui/surface";
import { Switch } from "@/components/ui/switch";
import { Text } from "@/components/ui/text";

import "./screens.css";

const meta = {
  title: "Screens/App Shell",
  parameters: {
    layout: "fullscreen",
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

const discussionItems = [
  {
    id: "command-line-interface",
    title: "Storyden command line interface",
    url: "/d/command-line-interface",
    text: "Ship, inspect, and operate a Storyden community from the terminal.",
  },
  {
    id: "community-handbook",
    title: "Planning our community handbook",
    url: "/d/community-handbook",
    text: "Collecting the policies and practical guidance that members need most often.",
  },
  {
    id: "moderation-workflow",
    title: "A calmer moderation workflow",
    url: "/d/moderation-workflow",
    text: "Ideas for keeping review work focused without turning every screen into an admin console.",
  },
];

export const CommunityFeed: Story = {
  render: () => (
    <main className="screen-preview">
      <aside className="screen-preview__sidebar">
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
      </aside>

      <section className="screen-preview__main">
        <div className="screen-preview__content">
          <CardBox
            as="section"
            aria-label="Share something"
            background="background.control"
            borderRadius="sm"
          >
            <Text color="text.muted">
              Share a thought, a link, something cool...
            </Text>
            <div className="screen-preview__composer-actions">
              <Button>Share</Button>
            </div>
          </CardBox>

          <CardRows items={discussionItems} />
        </div>
      </section>
    </main>
  ),
};

export const Settings: Story = {
  render: () => (
    <main className="screen-preview">
      <aside className="screen-preview__sidebar">
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
          </SidebarNavigationSection>
        </SidebarNavigationPane>
      </aside>

      <section className="screen-preview__main">
        <div className="screen-preview__content screen-preview__content--settings">
          <PageHeader
            title="Interface settings"
            description="Choose how Storyden looks and behaves for you."
            actions={<Button>Save</Button>}
          />

          <section className="screen-preview__settings-section">
            <div>
              <SectionHeading>Community identity</SectionHeading>
              <Text variant="supporting">
                This name appears in navigation and browser titles.
              </Text>
            </div>
            <FormControl>
              <label htmlFor="community-name">Community name</label>
              <Input id="community-name" defaultValue="Storyden" />
            </FormControl>
          </section>

          <section className="screen-preview__settings-section">
            <div>
              <SectionHeading>Display</SectionHeading>
              <Text variant="supporting">
                Keep compact controls visible while working through dense pages.
              </Text>
            </div>
            <Switch defaultChecked>Use compact interface</Switch>
            <div>
              <Badge colorPalette="green">Saved locally</Badge>
            </div>
          </section>
        </div>
      </section>
    </main>
  ),
};
