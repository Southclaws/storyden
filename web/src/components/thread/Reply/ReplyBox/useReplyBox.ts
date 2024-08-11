"use client";

import { FormEvent, useState } from "react";

import { replyCreate } from "src/api/openapi-client/replies";
import { useThreadGet } from "src/api/openapi-client/threads";
import { Thread } from "src/api/openapi-schema";
import { handleError } from "src/components/site/ErrorBanner";

type Value = {
  body: string;
  isEmpty: boolean;
};

export function useReplyBox(thread: Thread) {
  const { mutate } = useThreadGet(thread.slug);
  const [isLoading, setLoading] = useState(false);
  const [value, setValue] = useState<Value>({
    body: "",
    isEmpty: true,
  });
  const [resetKey, setResetKey] = useState("");

  function onChange(v: string, isEmpty: boolean) {
    setValue({
      body: v,
      isEmpty,
    });
  }

  async function doReply() {
    if (value.isEmpty) {
      return;
    }

    setLoading(true);
    await replyCreate(thread.id, { body: value.body })
      .catch(handleError)
      .then(async () => {
        await mutate();
        setValue({
          body: "",
          isEmpty: true,
        });

        // This is a little hack tbh, essentially if this prop for the
        // ContentComposer component changes, its value is reset. Could have
        // done it with a hook but... meh this is simpler (albeit not idiomatic)
        setResetKey(new Date().toISOString());

        setLoading(false);
        setTimeout(
          () =>
            window.scrollTo({
              behavior: "smooth",
              top: document.body.scrollHeight,
            }),
          100,
        );
      });
  }

  async function onReply(e: FormEvent) {
    e.preventDefault();
    doReply();
  }

  function onKeyDown(event: React.KeyboardEvent<HTMLFormElement>) {
    if (event.key == "Enter" && (event.metaKey || event.ctrlKey)) {
      doReply();
    }
  }

  return {
    onReply,
    onChange,
    onKeyDown,
    isLoading,
    resetKey,
    isEmpty: value.isEmpty,
  };
}
