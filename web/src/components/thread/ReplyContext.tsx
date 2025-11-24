"use client";

import { createContext, useContext, useState } from "react";

import { Reply, Thread } from "@/api/openapi-schema";

type ReplyToState = {
  thread: Thread;
  reply: Reply;
} | null;

type ReplyContextValue = {
  replyTo: ReplyToState;
  setReplyTo: (thread: Thread, reply: Reply) => void;
  clearReplyTo: () => void;
};

const ReplyContext = createContext<ReplyContextValue | undefined>(undefined);

export function ReplyProvider({ children }: { children: React.ReactNode }) {
  const [replyTo, setReplyToState] = useState<ReplyToState>(null);

  const setReplyTo = (thread: Thread, reply: Reply) => {
    setReplyToState({ thread, reply });
  };

  const clearReplyTo = () => {
    setReplyToState(null);
  };

  return (
    <ReplyContext.Provider value={{ replyTo, setReplyTo, clearReplyTo }}>
      {children}
    </ReplyContext.Provider>
  );
}

export function useReplyContext() {
  const context = useContext(ReplyContext);
  if (context === undefined) {
    throw new Error("useReplyContext must be used within a ReplyProvider");
  }
  return context;
}
