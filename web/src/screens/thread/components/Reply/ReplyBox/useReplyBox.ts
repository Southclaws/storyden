import { useToast } from "@chakra-ui/react";
import { useState } from "react";

import { postCreate } from "src/api/openapi/posts";
import { Thread } from "src/api/openapi/schemas";
import { useThreadGet } from "src/api/openapi/threads";
import { errorToast } from "src/components/ErrorBanner";

export function useReplyBox(thread: Thread) {
  const toast = useToast();
  const { mutate } = useThreadGet(thread.slug);
  const [isLoading, setLoading] = useState(false);
  const [value, setValue] = useState("");
  const [resetKey, setResetKey] = useState("");

  function onChange(v: string) {
    setValue(v);
  }

  async function onReply() {
    setLoading(true);
    await postCreate(thread.id, {
      body: {
        type: "application/json",
        value,
      },
    })
      .catch(errorToast(toast))
      .then(async () => {
        await mutate();
        setValue("");

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
          100
        );
      });
  }

  function onKeyDown(event: React.KeyboardEvent<HTMLDivElement>) {
    if (event.key == "Enter" && (event.metaKey || event.ctrlKey)) {
      onReply();
    }
  }

  return {
    onReply,
    onChange,
    onKeyDown,
    isLoading,
    resetKey,
  };
}
