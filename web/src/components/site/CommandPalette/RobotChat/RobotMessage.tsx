import { isDataUIPart, isToolUIPart } from "ai";
import Link from "next/link";

import { useNodeGet } from "@/api/openapi-client/nodes";
import { usePostLocationGet } from "@/api/openapi-client/posts";
import { useProfileGet } from "@/api/openapi-client/profiles";
import { useThreadGet } from "@/api/openapi-client/threads";
import { ProfileReference, Thread } from "@/api/openapi-schema";
import { RobotRenderCardData, StorydenUIMessage } from "@/api/robots-types";
import { ContentComposerMarkdown } from "@/components/content/ContentComposerMarkdown/ContentComposerMarkdown";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { CardBox } from "@/components/ui/card-box";
import { ReplyIcon } from "@/components/ui/icons/Reply";
import { Card } from "@/components/ui/surface";
import { Text } from "@/components/ui/text";
import { css } from "@/styled-system/css";
import {
  Box,
  HStack,
  LStack,
  VStack,
  WStack,
  styled,
} from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";
import { htmlToMarkdown } from "@/utils/markdown";

import { Timestamp } from "../../Timestamp";

import styles from "./RobotMessage.module.css";

import { RobotDelegation } from "./RobotDelegation";
import {
  ConfirmationPart,
  RobotToolCall,
  RobotToolConfirmationBatch,
  isConfirmationToolPart,
} from "./RobotToolCall";

type Props = {
  id: string;
  role: StorydenUIMessage["role"];
  parts: readonly StorydenUIMessage["parts"][number][];
  author?: ProfileReference;
  isCurrentMemberMessage?: boolean;
  isNewestUserMessage?: boolean;
};

export function RobotMessage({
  id,
  role,
  parts,
  author,
  isCurrentMemberMessage = false,
  isNewestUserMessage = false,
}: Props) {
  const isUser = role === "user";
  const isOwnMessage = isUser && (!author || isCurrentMemberMessage);
  let authorLabel = "Robot";
  if (isUser) {
    authorLabel =
      author && !isCurrentMemberMessage ? `@${author.handle}` : "You";
  }

  return (
    <VStack
      id={`robot-message-${id}`}
      role="article"
      aria-roledescription="chat message"
      aria-label={`${authorLabel} message`}
      w="full"
      minW="0"
      gap="2"
      alignItems={isUser ? "flex-end" : "flex-start"}
      className={isNewestUserMessage ? styles["newestUserMessage"] : undefined}
    >
      {isUser && author && !isOwnMessage && (
        <MemberBadge profile={author} avatar="hidden" size="xs" name="handle" />
      )}
      {renderMessageParts(parts, isUser)}
    </VStack>
  );
}

function renderMessageParts(
  parts: readonly StorydenUIMessage["parts"][number][],
  isUser: boolean,
) {
  const confirmationParts = parts.filter(
    isConfirmationToolPart,
  ) as ConfirmationPart[];
  const shouldBatchConfirmations = !isUser && confirmationParts.length > 1;
  let renderedConfirmationBatch = false;

  return parts.flatMap((part, idx) => {
    if (shouldBatchConfirmations && isConfirmationToolPart(part)) {
      if (renderedConfirmationBatch) {
        return [];
      }

      renderedConfirmationBatch = true;

      return [
        <RobotToolConfirmationBatch
          key={`confirmation-batch:${confirmationParts.map((p) => p.toolCallId).join(":")}`}
          parts={confirmationParts as ConfirmationPart[]}
        />,
      ];
    }

    return [
      <RobotMessagePart key={partKey(part, idx)} part={part} isUser={isUser} />,
    ];
  });
}

function RobotMessagePart({
  part,
  isUser,
}: {
  part: StorydenUIMessage["parts"][number];
  isUser: boolean;
}) {
  if (isToolUIPart(part)) {
    return <RobotToolCall part={part} />;
  }

  if (isDataUIPart(part)) {
    switch (part.type) {
      case "data-delegation":
        return (
          <RobotDelegation data={part.data} renderParts={renderMessageParts} />
        );
      case "data-render_card":
        return (
          // padding 1 since cardbox has shadow (should remove shadow in future)
          <Box px="1" w="full">
            <RobotRenderCard data={part.data} />
          </Box>
        );
    }
  }

  if (part.type === "reasoning") {
    if (!("text" in part) || !part.text) {
      return null;
    }

    return (
      <styled.details w="full" color="text.subtle" fontSize="sm">
        <styled.summary
          cursor="pointer"
          w="fit"
          fontSize="xs"
          _hover={{ color: "text.default" }}
        >
          Thought
        </styled.summary>
        <Box
          mt="2"
          pl="3"
          borderLeftWidth="thin"
          borderLeftColor="border.muted"
          color="text.muted"
          className={styles["messageText"]}
        >
          <ContentComposerMarkdown
            disabled
            initialValue={part.text}
            initialValueFormat="markdown"
          />
        </Box>
      </styled.details>
    );
  }

  if (part.type === "text") {
    if (!("text" in part) || !part.text) {
      return null;
    }

    return (
      <Box
        className={styles["messageText"]}
        bg={isUser ? "background.inset" : "transparent"}
        borderRadius={isUser ? "2xl" : "none"}
        px={isUser ? "4" : "0"}
        pt={isUser ? "2" : "1"}
        w={isUser ? "fit" : "full"}
        minW="0"
        maxW={isUser ? "3/4" : "full"}
        overflow="hidden"
        color="text.default"
        lineHeight="relaxed"
      >
        <ContentComposerMarkdown
          disabled
          initialValue={part.text}
          initialValueFormat="markdown"
        />
      </Box>
    );
  }

  return null;
}

function RobotRenderCard({ data }: { data: RobotRenderCardData }) {
  switch (data.kind) {
    case "node":
      return <RobotNodeCard data={data} />;
    case "thread":
      return <RobotThreadCard data={data} />;
    case "profile":
      return <RobotProfileCard data={data} />;
    case "reply":
      return <RobotReplyCard data={data} />;
    default:
      return <RobotUnavailableCard label="Referenced resource" />;
  }
}

function RobotNodeCard({ data }: { data: RobotRenderCardData }) {
  const {
    data: page,
    error,
    isLoading,
  } = useNodeGet(data.id, undefined, {
    swr: { enabled: Boolean(data.id) },
  });

  if (!data.id || error) {
    return <RobotUnavailableCard label="Library page" />;
  }

  if (isLoading || !page) {
    return <RobotLoadingCard label="Library page" />;
  }

  const url = `/l/${page.slug}`;

  return (
    <Box w="full" maxW="4/5">
      <Card
        id={page.id}
        title={page.name}
        text={page.description}
        url={url}
        image={getAssetURL(page.primary_image?.path)}
        shape="row"
      />
    </Box>
  );
}

function RobotThreadCard({ data }: { data: RobotRenderCardData }) {
  const {
    data: thread,
    error,
    isLoading,
  } = useThreadGet(data.id, undefined, {
    swr: { enabled: Boolean(data.id) },
  });

  if (!data.id || error) {
    return <RobotUnavailableCard label="Thread" />;
  }

  if (isLoading || !thread) {
    return <RobotLoadingCard label="Thread" />;
  }

  return (
    <Box w="full" maxW="4/5">
      <RobotThreadReferenceCard thread={thread} />
    </Box>
  );
}

function RobotProfileCard({ data }: { data: RobotRenderCardData }) {
  const {
    data: profile,
    error,
    isLoading,
  } = useProfileGet(data.id, {
    swr: { enabled: Boolean(data.id) },
  });

  if (!data.id || error) {
    return <RobotFallbackLinkCard data={data} label="Profile" />;
  }

  if (isLoading || !profile) {
    return <RobotLoadingCard label="Profile" />;
  }

  return (
    <Box w="full" maxW="4/5">
      {/* <Card
        id={profile.id}
        title={profile.name}
        image={getAvatarURL(profile.id)}
        url={`/m/${profile.handle}`}
        shape="row"
        content={profile.bio}
      >
        <MemberBadge profile={profile} size="md" name="full-vertical" />
      </Card> */}

      <CardBox>
        <WStack>
          <MemberBadge profile={profile} size="md" name="full-vertical" />

          <Text as="span" variant="supporting">
            {"joined "}
            <Timestamp created={profile.createdAt} />
            {" ago"}
          </Text>
        </WStack>
      </CardBox>
    </Box>
  );
}

function RobotReplyCard({ data }: { data: RobotRenderCardData }) {
  const {
    data: location,
    error,
    isLoading,
  } = usePostLocationGet(
    { id: data.id },
    {
      swr: { enabled: Boolean(data.id) },
    },
  );
  const isReplyLocation = location?.kind === "reply";
  const page = isReplyLocation ? (location.page ?? 1) : undefined;
  const {
    data: thread,
    error: threadError,
    isLoading: isThreadLoading,
  } = useThreadGet(
    isReplyLocation ? location.slug : "",
    page ? { page: String(page) } : undefined,
    {
      swr: { enabled: Boolean(data.id && isReplyLocation && location.slug) },
    },
  );

  if (!data.id || error) {
    return <RobotFallbackLinkCard data={data} label="Reply" />;
  }

  if (isLoading || !location) {
    return <RobotLoadingCard label="Reply" />;
  }

  if (!isReplyLocation) {
    return <RobotFallbackLinkCard data={data} label="Reply" />;
  }

  const url =
    (page ?? 1) > 1
      ? `/t/${location.slug}?page=${page}#post-${data.id}`
      : `/t/${location.slug}#post-${data.id}`;

  if (threadError) {
    return <RobotFallbackLinkCard data={data} label="Reply" />;
  }

  if (isThreadLoading || !thread) {
    return <RobotLoadingCard label="Reply" />;
  }

  const reply = thread.replies.replies.find((reply) => reply.id === data.id);

  return (
    <CardBox w="full" maxW="4/5">
      <LStack>
        <WStack color="text.subtle" fontSize="xs">
          <Link
            href={url}
            className={css({
              color: "accent.text",
              fontWeight: "medium",
              _hover: { textDecoration: "underline" },
            })}
          >
            <HStack gap="1" alignItems="center">
              <ReplyIcon width="4" height="4" aria-hidden />
              <span>Reply in this thread</span>
            </HStack>
          </Link>
          <Box>
            {reply && (
              <MemberBadge
                profile={reply.author}
                avatar="visible"
                size="xs"
                name="handle"
              />
            )}
          </Box>
        </WStack>

        <styled.blockquote fontSize="sm" display="flex">
          <span>“</span>
          <styled.span lineClamp={1}>{reply?.description}</styled.span>
          <span>”</span>
        </styled.blockquote>

        <RobotThreadReferenceCard thread={thread} url={url} />
      </LStack>
    </CardBox>
  );
}

function RobotThreadReferenceCard({
  thread,

  url = `/t/${thread.slug}`,
}: {
  thread: Thread;

  url?: string;
}) {
  const title = thread.title || thread.link?.title || "Untitled thread";
  const image = getAssetURL(
    thread.assets?.[0]?.path ?? thread.link?.primary_image?.path,
  );
  const replyCount = thread.reply_status.replies;
  const replyLabel = replyCount === 1 ? "1 reply" : `${replyCount} replies`;
  const text = resourceSnippet(thread.description, thread.body);

  return (
    <Card
      id={thread.id}
      title={title}
      text={text}
      url={url}
      image={image}
      shape="row"
    >
      <WStack>
        <HStack gap="2" minW="0" color="text.subtle" fontSize="sm">
          <MemberBadge
            profile={thread.author}
            avatar="visible"
            size="xs"
            name="handle"
          />
          <styled.span color="text.muted">·</styled.span>
          <Timestamp created={thread.createdAt} />
        </HStack>

        <Text as="span" variant="supporting">
          {replyLabel}
        </Text>
      </WStack>
    </Card>
  );
}

function resourceSnippet(description?: string, body?: string) {
  const value =
    description?.trim() || (body ? htmlToMarkdown(body).trim() : "");

  if (!value) {
    return undefined;
  }

  return value.length > 220 ? `${value.slice(0, 220).trim()}…` : value;
}

function RobotFallbackLinkCard({
  data,
  label,
}: {
  data: RobotRenderCardData;
  label: string;
}) {
  return (
    <Box w="full" maxW="4/5">
      <Card
        id={data.ref}
        title={label}
        text="Open this referenced resource."
        url={`/_/resolve/${data.kind}/${data.id}`}
        shape="row"
      />
    </Box>
  );
}

function RobotLoadingCard({ label }: { label: string }) {
  return <Text variant="supporting">Loading {label.toLowerCase()}...</Text>;
}

function RobotUnavailableCard({ label }: { label: string }) {
  return <Text variant="supporting">{label} unavailable.</Text>;
}

function partKey(part: StorydenUIMessage["parts"][number], idx: number) {
  if (isToolUIPart(part)) {
    return `${part.toolCallId}:${part.state}`;
  }

  if (isDataUIPart(part)) {
    return `${part.type}:${part.id ?? idx}`;
  }

  return `${part.type}:${idx}`;
}
