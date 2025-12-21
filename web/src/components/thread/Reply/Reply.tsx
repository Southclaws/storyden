import Link from "next/link";
import { Controller, ControllerProps } from "react-hook-form";

import { Reply as ReplyType, Thread } from "@/api/openapi-schema";
import { ContentComposer } from "@/components/content/ContentComposer/ContentComposer";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { CancelAction } from "@/components/site/Action/Cancel";
import { SaveAction } from "@/components/site/Action/Save";
import { Timestamp } from "@/components/site/Timestamp";
import { ReplyIcon } from "@/components/ui/icons/Reply";
import { CardBox, HStack, WStack, styled } from "@/styled-system/jsx";
import { hstack } from "@/styled-system/patterns";

import { Byline } from "../../content/Byline";
import { ReactList } from "../ReactList/ReactList";
import { ReplyMenu } from "../ReplyMenu/ReplyMenu";

import { ReplyToButton } from "./ReplyToButton";
import { useFragmentScroll } from "./useFragmentScroll";
import { Form, Props, useReply } from "./useReply";

export function Reply(props: Props) {
  const { isEmpty, isEditing, resetKey, form, handlers } = useReply(props);
  const isTargeted = useFragmentScroll(props.reply.id);

  const { thread, reply, currentPage } = props;

  return (
    <CardBox
      id={reply.id}
      data-targeted={isTargeted || undefined}
      _target={{
        scrollMarginTop: {
          base: "0",
          md: "20",
        },
        animation: "target-pulse",
      }}
      css={{
        "&[data-targeted]": {
          animation: "target-pulse",
        },
      }}
    >
      <styled.form
        display="flex"
        flexDirection="column"
        gap="2"
        onSubmit={handlers.handleSave}
      >
        <WStack>
          <Byline
            href={`#${reply.id}`}
            author={reply.author}
            time={new Date(reply.createdAt)}
            updated={new Date(reply.updatedAt)}
          />

          {isEditing ? (
            <HStack>
              <>
                <CancelAction
                  type="button"
                  onClick={handlers.handleDiscardChanges}
                >
                  Discard
                </CancelAction>
                <SaveAction type="submit" disabled={isEmpty}>
                  Save
                </SaveAction>
              </>
            </HStack>
          ) : (
            <HStack>
              <ReplyToButton thread={thread} reply={reply} />
              <ReplyMenu
                thread={thread}
                reply={reply}
                currentPage={currentPage}
                onEdit={handlers.handleSetEditing}
              />
            </HStack>
          )}
        </WStack>

        {reply.reply_to && <InReplyTo to={reply.reply_to} thread={thread} />}

        <ReplyBodyInput
          control={form.control}
          name="body"
          initialValue={reply.body}
          resetKey={resetKey}
          disabled={!isEditing}
          handleEmptyStateChange={handlers.handleEmptyStateChange}
        />
      </styled.form>

      <ReactList thread={thread} reply={reply} currentPage={currentPage} />
    </CardBox>
  );
}

type ReplyBodyInputProps = Omit<ControllerProps<Form>, "render"> & {
  initialValue: string;
  resetKey: string;
  handleEmptyStateChange: (isEmpty: boolean) => void;
};

function ReplyBodyInput({
  control,
  name,
  initialValue,
  resetKey,
  disabled,
  handleEmptyStateChange,
}: ReplyBodyInputProps) {
  return (
    <Controller<Form>
      render={({ field: { onChange } }) => {
        function handleChange(value: string, isEmpty: boolean) {
          handleEmptyStateChange(isEmpty);
          onChange(value);
        }

        return (
          <ContentComposer
            initialValue={initialValue}
            onChange={handleChange}
            resetKey={resetKey}
            disabled={disabled}
          />
        );
      }}
      control={control}
      name={name}
    />
  );
}

function InReplyTo({ to, thread }: { to: ReplyType; thread: Thread }) {
  // figure out if the reply-to is on the current page, then  do a fragment link
  // if on same page, otherwise use /t/locate to navigate to the right page.
  const isOnCurrentPage = thread.replies.replies.some((r) => r.id === to.id);
  const permalink = isOnCurrentPage ? `#${to.id}` : `/t/locate/${to.id}`;

  // NOTE: because nextjs does some weird shit, we gotta use a normal anchor
  // for fragment navigation, otherwise it breaks :target etc for some reason.
  const AnchorComponent = isOnCurrentPage ? styled.a : Link;

  return (
    <WStack
      gap="1"
      fontSize="xs"
      color="fg.muted"
      px="2"
      py="1"
      borderRadius="md"
      bgColor="bg.subtle"
      w="full"
      minW="0"
    >
      <AnchorComponent
        href={permalink}
        className={hstack({
          minW: "0",
          flexShrink: "1",
        })}
      >
        <ReplyIcon w="4" minW="4" />
        <styled.span
          minW="0"
          overflow="hidden"
          textOverflow="ellipsis"
          whiteSpace="nowrap"
          lineClamp="1"
        >
          “{to.description}”
        </styled.span>
      </AnchorComponent>

      <HStack flexShrink="0" minW="0">
        <MemberBadge
          profile={to.author}
          size="xs"
          name="handle"
          avatar="visible"
        />
        <AnchorComponent href={permalink}>
          <Timestamp created={to.createdAt} />
        </AnchorComponent>
      </HStack>
    </WStack>
  );
}
