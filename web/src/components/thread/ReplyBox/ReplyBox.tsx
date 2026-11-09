import { useState } from "react";

import { usePublicRegistration } from "@/lib/settings/registration";

import { useReplyContext } from "../ReplyContext";

import { ReplyBoxView } from "./ReplyBoxView";
import { Props, useReplyBox } from "./useReplyBox";

export function ReplyBox(props: Props) {
  const { replyTo, clearReplyTo } = useReplyContext();
  const state = useReplyBox(props);
  const registrationAllowed = usePublicRegistration(props.initialSettings);
  const [isActivated, setIsActivated] = useState(false);
  const [focusOnReady, setFocusOnReady] = useState(false);
  const isActive = isActivated || !!replyTo;
  const isReplyingDisabled = Boolean(props.thread.locked_at) && !state.canReply;

  function handleActivate() {
    if (isReplyingDisabled) {
      return;
    }

    setIsActivated(true);
    setFocusOnReady(true);
  }

  function handleWarmEditor() {
    setIsActivated(true);
  }

  function handleClearReplyTo() {
    setIsActivated(true);
    clearReplyTo();
  }

  return (
    <ReplyBoxView
      thread={props.thread}
      state={state}
      replyTo={replyTo?.reply ?? null}
      registrationAllowed={registrationAllowed}
      isActive={isActive}
      focusOnReady={focusOnReady}
      onActivate={handleActivate}
      onWarmEditor={handleWarmEditor}
      onClearReplyTo={handleClearReplyTo}
    />
  );
}
