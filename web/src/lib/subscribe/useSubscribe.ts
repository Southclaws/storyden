"use client";

import { useEffect } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { ZodSchema, z } from "zod";

import { WS_ADDRESS } from "@/config";

export type Channel = "thread";

export function useSubscribe<T>(
  chan: Channel,
  messageSchema: ZodSchema<T>,
  cb: (event: T) => Promise<void>,
) {
  const { readyState, sendJsonMessage, lastJsonMessage } = useWebSocket(
    WS_ADDRESS,
    { share: true },
  );

  useEffect(() => {
    if (readyState !== ReadyState.OPEN) {
      return;
    }
    const event = messageSchema.parse(lastJsonMessage);

    cb(event);
  }, [readyState, cb, messageSchema, lastJsonMessage]);

  useEffect(() => {
    if (readyState !== ReadyState.OPEN) {
      return;
    }

    sendJsonMessage({ type: "subscribe", channel: chan });
  }, [readyState, chan, sendJsonMessage]);
}
