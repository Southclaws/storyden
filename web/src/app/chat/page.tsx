"use client";

import { useChat } from "@ai-sdk/react";
import { DefaultChatTransport } from "ai";
import { FormEvent, useMemo, useState } from "react";

import { API_ADDRESS } from "@/config";

function summariseParts(parts: { type: string }[] = []) {
  if (!Array.isArray(parts)) {
    return "";
  }

  return parts
    .map((part) => {
      if (!part || typeof part !== "object") {
        return "";
      }

      switch (part.type) {
        case "text":
          return "text" in part ? (part.text ?? "") : "";
        case "reasoning":
          return "text" in part ? (part.text ?? "") : "";
        default:
          return `[${part.type}]`;
      }
    })
    .filter(Boolean)
    .join("\n");
}

export default function ChatPage() {
  const [input, setInput] = useState("");
  const transport = useMemo(
    () =>
      new DefaultChatTransport({
        api: `${API_ADDRESS}/sse/chat`,
        credentials: "include",
      }),
    [],
  );

  const { messages, sendMessage, status, error, clearError } = useChat({
    transport,
  });

  const isBusy = status === "submitted" || status === "streaming";

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!input.trim()) {
      return;
    }

    try {
      await sendMessage({ text: input.trim() });
      setInput("");
      if (error) {
        clearError();
      }
    } catch (err) {
      console.error("sendMessage failed", err);
    }
  };

  return (
    <div
      style={{
        padding: "2rem",
        maxWidth: "720px",
        margin: "0 auto",
        display: "flex",
        flexDirection: "column",
        gap: "1.5rem",
      }}
    >
      <div>
        <h1
          style={{
            fontSize: "1.5rem",
            fontWeight: 600,
            marginBottom: "0.5rem",
          }}
        >
          Storyden Agent Chat
        </h1>
        <p style={{ color: "#666" }}>
          Talk to the Storyden agent and watch the SSE bridge stream responses
          in real time.
        </p>
      </div>

      <div
        style={{
          border: "1px solid #e2e2e2",
          borderRadius: "12px",
          padding: "1rem",
          minHeight: "320px",
          display: "flex",
          flexDirection: "column",
          gap: "0.75rem",
          background: "#fff",
        }}
      >
        {messages.length === 0 && (
          <p style={{ color: "#888" }}>
            No conversation yet. Ask about your library pages to begin.
          </p>
        )}
        {messages.map((message) => {
          const content = summariseParts(message.parts);
          return (
            <div
              key={message.id}
              style={{
                display: "flex",
                flexDirection: "column",
                gap: "0.35rem",
              }}
            >
              <span style={{ fontSize: "0.85rem", fontWeight: 600 }}>
                {message.role === "user" ? "You" : "Agent"}
              </span>
              <div
                style={{
                  background: message.role === "user" ? "#f4f4f5" : "#eef6ff",
                  borderRadius: "8px",
                  padding: "0.75rem",
                  whiteSpace: "pre-wrap",
                }}
              >
                {content}
              </div>
            </div>
          );
        })}
      </div>

      {error && <p style={{ color: "#d14343" }}>{error.message}</p>}

      <form
        onSubmit={onSubmit}
        style={{ display: "flex", gap: "0.75rem", alignItems: "flex-start" }}
      >
        <textarea
          value={input}
          onChange={(event) => setInput(event.target.value)}
          placeholder="Ask the Storyden agent..."
          rows={3}
          style={{
            flex: 1,
            borderRadius: "10px",
            border: "1px solid #d7d7db",
            padding: "0.75rem",
            resize: "vertical",
          }}
        />
        <button
          type="submit"
          disabled={isBusy || !input.trim()}
          style={{
            background: "#111827",
            color: "#fff",
            border: "none",
            borderRadius: "10px",
            padding: "0.75rem 1.25rem",
            cursor: "pointer",
            opacity: isBusy || !input.trim() ? 0.6 : 1,
          }}
        >
          {isBusy ? "Sending…" : "Send"}
        </button>
        <button
          type="button"
          onClick={() => setInput("")}
          style={{
            border: "none",
            background: "transparent",
            color: "#6b7280",
            cursor: "pointer",
          }}
        >
          Clear
        </button>
      </form>
    </div>
  );
}
