import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useRouter } from "@storybook/nextjs-vite/navigation.mock";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { fn } from "storybook/test";
import { SWRConfig, unstable_serialize } from "swr";

import { getAccountGetKey } from "@/api/openapi-client/accounts";
import { getGetInfoKey } from "@/api/openapi-client/misc";
import type { Account, ProfileReference, Reply } from "@/api/openapi-schema";
import { DefaultSettings } from "@/lib/settings/settings";
import { LStack } from "@/styled-system/jsx";

import { ReplyBoxView, type ReplyBoxViewProps } from "./ReplyBoxView";
import type { Form } from "./useReplyBox";

type Scenario =
  | "idle"
  | "editing"
  | "draft"
  | "submitting"
  | "reply-to"
  | "locked-member"
  | "locked-draft"
  | "locked-moderator"
  | "guest"
  | "guest-private"
  | "guest-locked"
  | "posted";

type PreviewProps = {
  scenario: Scenario;
  editor: "richtext" | "markdown";
  narrow: boolean;
  longHandle: boolean;
};

const author: ProfileReference = {
  id: "mara",
  handle: "mara",
  name: "Mara",
  joined: "2026-10-01T12:00:00Z",
  roles: [],
};

const reply: Reply = {
  id: "ellis-reply",
  createdAt: "2026-11-08T12:00:00Z",
  updatedAt: "2026-11-08T12:00:00Z",
  author: { ...author, id: "ellis", handle: "ellis", name: "Ellis" },
  title: "Planning our next community meetup",
  slug: "community-meetup",
  root_id: "community-meetup",
  root_slug: "community-meetup",
  visibility: "published",
  body: "<p>Could we meet on Saturday instead?</p>",
  body_links: [],
  assets: [],
  reacts: [],
  collections: { in_collections: 0, has_collected: false },
  likes: { likes: 0, liked: false },
};

const draft = "<p>Saturday works for me. I can bring coffee for everyone.</p>";

const meta = {
  title: "Compositions/Discussion/Reply Box",
  component: ReplyBoxPreview,
  tags: ["autodocs"],
  parameters: {
    layout: "padded",
    nextjs: { appDirectory: false },
    docs: {
      description: {
        component:
          "The thread reply task: recipient context, editor, posting action, and permission state. Stories render the production view with local form state; posting only resets this preview. The editor uses the shared content composer; reply targeting is above it, while lock context and the posting action sit below it. Use this for threaded conversation, rather than an ordinary multiline field or a standalone content card. Editor mode and narrow width are available in Controls. Submission errors currently use the application’s shared error handling and have no dedicated ReplyBox layout.",
      },
    },
  },
  beforeEach: () => {
    useRouter.mockReturnValue({
      back: fn(),
      forward: fn(),
      refresh: fn(),
      push: fn(),
      replace: fn(),
      prefetch: fn(),
      bfcacheId: "reply-box-story",
    });
  },
  args: {
    scenario: "idle",
    editor: "richtext",
    narrow: false,
    longHandle: false,
  },
  argTypes: {
    scenario: {
      control: "select",
      options: [
        "idle",
        "editing",
        "draft",
        "submitting",
        "reply-to",
        "locked-member",
        "locked-draft",
        "locked-moderator",
        "guest",
        "guest-private",
        "guest-locked",
        "posted",
      ],
    },
    editor: { control: "radio", options: ["richtext", "markdown"] },
  },
  render: (args) => (
    <ReplyBoxPreview key={`${args.scenario}-${args.editor}`} {...args} />
  ),
} satisfies Meta<PreviewProps>;

export default meta;
type Story = StoryObj<typeof meta>;

function ReplyBoxPreview(props: PreviewProps) {
  const account: Account = {
    ...author,
    roles: [],
    createdAt: author.joined,
    updatedAt: author.joined,
    bio: "",
    links: [],
    meta: { editor: { mode: props.editor } },
    verified_status: "none",
    kind: "human",
    email_addresses: [],
    auth_services: [],
    admin: false,
  };

  return (
    <SWRConfig
      value={{
        isPaused: () => true,
        fallback: {
          [unstable_serialize(getAccountGetKey())]: account,
          [unstable_serialize(getGetInfoKey())]: DefaultSettings,
        },
      }}
    >
      <LStack w="full" maxW={props.narrow ? "sm" : "3xl"} alignItems="stretch">
        <ReplyBoxScenario {...props} />
      </LStack>
    </SWRConfig>
  );
}

function ReplyBoxScenario(props: PreviewProps) {
  const { scenario } = props;
  const isLocked = scenario.startsWith("locked") || scenario === "guest-locked";
  const hasDraft = ["draft", "submitting", "reply-to", "locked-draft"].includes(
    scenario,
  );
  const form = useForm<Form>({
    defaultValues: { body: hasDraft ? draft : "" },
  });
  const [isEmpty, setEmpty] = useState(!hasDraft);
  const [isActive, setActive] = useState(
    scenario !== "idle" &&
      scenario !== "locked-member" &&
      scenario !== "guest-locked",
  );
  const [focusOnReady, setFocusOnReady] = useState(false);
  const [draftKey, setDraftKey] = useState(0);
  const [replyTo, setReplyTo] = useState<Reply | null>(
    scenario === "reply-to" ? reply : null,
  );
  const [postedReply, setPostedReply] = useState<
    ReplyBoxViewProps["state"]["postedReply"]
  >(
    scenario === "posted"
      ? { id: "posted-reply", pageNumber: 3, permalink: "#posted-reply" }
      : null,
  );
  const thread = {
    author: {
      ...author,
      handle: props.longHandle
        ? "mara-with-a-very-long-community-handle"
        : author.handle,
    },
    locked_at: isLocked ? "2026-11-09T10:00:00Z" : undefined,
  };
  const state: ReplyBoxViewProps["state"] = {
    isLoggedIn: !scenario.startsWith("guest"),
    canReply: !isLocked || scenario === "locked-moderator",
    isEmpty,
    isLoading: scenario === "submitting",
    draftKey: String(draftKey),
    postedReply,
    form,
    handlers: {
      handleSubmit: form.handleSubmit(() => {
        form.reset({ body: "" });
        setDraftKey((value) => value + 1);
        setEmpty(true);
        setReplyTo(null);
        setPostedReply({
          id: "posted-reply",
          pageNumber: 3,
          permalink: "#posted-reply",
        });
      }),
      handleEmptyStateChange: setEmpty,
      handleReplyPostedAdmonitionClose: () => setPostedReply(null),
      handleReplyNavigation: () => setPostedReply(null),
    },
  };
  const viewProps: ReplyBoxViewProps = {
    thread,
    state,
    replyTo: replyTo && {
      ...replyTo,
      author: {
        ...replyTo.author,
        handle: props.longHandle
          ? "ellis-with-a-very-long-community-handle"
          : replyTo.author.handle,
      },
    },
    registrationAllowed: scenario !== "guest-private",
    isActive,
    focusOnReady,
    onActivate: () => {
      setActive(true);
      setFocusOnReady(true);
    },
    onWarmEditor: () => setActive(true),
    onClearReplyTo: () => {
      setActive(true);
      setReplyTo(null);
    },
  };

  return <ReplyBoxView {...viewProps} />;
}

export const Idle: Story = {};
export const Editing: Story = { args: { scenario: "editing" } };
export const WithDraft: Story = { args: { scenario: "draft" } };
export const Submitting: Story = {
  args: { scenario: "submitting" },
  parameters: {
    docs: {
      description: {
        story:
          "The posting action shows its pending state while the editor retains the draft.",
      },
    },
  },
};
export const ReplyToMessage: Story = { args: { scenario: "reply-to" } };
export const LockedMember: Story = { args: { scenario: "locked-member" } };
export const LockedWithDraft: Story = {
  args: { scenario: "locked-draft" },
  parameters: {
    docs: {
      description: {
        story:
          "The thread was locked after the member started writing. The draft remains visible in a disabled composer.",
      },
    },
  },
};
export const LockedModerator: Story = {
  args: { scenario: "locked-moderator" },
};
export const Guest: Story = { args: { scenario: "guest" } };
export const GuestRegistrationClosed: Story = {
  args: { scenario: "guest-private" },
};
export const GuestLockedThread: Story = { args: { scenario: "guest-locked" } };
export const PostedOnAnotherPage: Story = { args: { scenario: "posted" } };
export const NarrowLongHandle: Story = {
  args: { scenario: "reply-to", narrow: true, longHandle: true },
};
export const Markdown: Story = {
  args: { scenario: "editing", editor: "markdown" },
};
export const NarrowLockedModerator: Story = {
  args: { scenario: "locked-moderator", narrow: true },
};
