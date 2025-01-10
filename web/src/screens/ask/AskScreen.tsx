"use client";

import { values } from "lodash";
import Link from "next/link";
import React, {
  AnchorHTMLAttributes,
  ClassAttributes,
  useEffect,
  useState,
} from "react";
import ReactMarkdown, { Components } from "react-markdown";

import { handle } from "@/api/client";
import { useNodeGet } from "@/api/openapi-client/nodes";
import { useThreadGet } from "@/api/openapi-client/threads";
import { Account, DatagraphItemKind } from "@/api/openapi-schema";
import {
  DatagraphItemNodeCard,
  DatagraphItemPostGenericCard,
} from "@/components/datagraph/DatagraphItemCard";
import {
  LoginAnchor,
  RegisterAnchor,
} from "@/components/site/Navigation/Anchors/Login";
import { UnreadyBanner } from "@/components/site/Unready";
import { Spinner } from "@/components/ui/Spinner";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { Input } from "@/components/ui/input";
import { API_ADDRESS, WEB_ADDRESS } from "@/config";
import { useCapability } from "@/lib/settings/capabilities";
import { css } from "@/styled-system/css";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { hstack, lstack } from "@/styled-system/patterns";
import { deriveError } from "@/utils/error";

type DatagraphRef = {
  id: string;
  kp: string;
  href: string;
};

type Props = {
  session?: Account;
};

export function AskScreen({ session }: Props) {
  if (!session) {
    return (
      <UnreadyBanner error="You must be logged in to use the knowledgebase Ask tool.">
        <WStack>
          <RegisterAnchor />
          <LoginAnchor />
        </WStack>
      </UnreadyBanner>
    );
  }

  return <Ask />;
}

export function Ask() {
  const [question, setQuestion] = useState("");
  const [content, setContent] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [sources, setSources] = useState<Record<string, DatagraphRef>>({});
  const isEnabled = useCapability("semdex");

  useEffect(() => {
    // Helper to extract SDR references from content
    const extractSources = (text: string): Record<string, DatagraphRef> => {
      const sdrRegex = /sdr:(\w+)\/([\w-]+)/g;
      const refs = {};

      text.replace(sdrRegex, (_, kind, id) => {
        if (id.length === 20) {
          const kp = getRouteForKind(kind);
          refs[id] = { id, kp: kp, href: `${WEB_ADDRESS}/${kp}/${id}` };
        }
        return "";
      });

      return refs;
    };

    const newSources = extractSources(content);

    setSources((current) => ({ ...current, ...newSources }));
  }, [content]);

  // Helper to replace SDR URLs with frontend links
  const replaceSdrUrls = (text: string): string => {
    const sdrRegex = /sdr:(\w+)\/([\w-]+)/g;
    return text.replace(sdrRegex, (_, kind, id) => {
      const kindRoute = getRouteForKind(kind);

      const url = `${WEB_ADDRESS}/${kindRoute}/${id}`;

      return `[${url}](${url})`;
    });
  };

  const fetchAnswer = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setContent("");
    setIsLoading(true);

    await handle(
      async () => {
        const response = await fetch(
          `${API_ADDRESS}/api/datagraph/qna?q=${encodeURIComponent(question)}`,
          {
            method: "GET",
            mode: "cors",
            credentials: "include",
          },
        );

        if (response.status === 404) {
          throw new Error("No answer could be found for this question.");
        }

        if (!response.ok) {
          throw new Error(`Error: ${response.statusText}`);
        }

        if (!response.body) {
          throw new Error(`Error: response is empty`);
        }

        setSources({});

        const reader = response.body.getReader();
        const decoder = new TextDecoder("utf-8");
        let done = false;

        while (!done) {
          const { value, done: readerDone } = await reader.read();
          done = readerDone;
          const chunk = decoder.decode(value, { stream: true });
          setContent((prev) => prev + chunk);
        }
      },
      {
        errorToast: false,
        async onError(e) {
          setContent(deriveError(e));
        },
        async cleanup() {
          setIsLoading(false);
        },
      },
    );
  };

  const components: Components = {
    a: ({ href, children }: AnchorHTMLAttributes<HTMLAnchorElement>) => {
      if (href === undefined) {
        return <a href={href}>{children}</a>;
      }

      try {
        const [kind, id] = sourceDataFromURL(href);
        if (!kind || !id) {
          throw new Error("Invalid SDR");
        }

        return (
          <Link
            className={css({
              background: "bg.muted",
              px: "2",
              py: "1",
              borderRadius: "sm",
            })}
            href={href}
          >
            {href}
          </Link>
        );
      } catch (e) {
        return <a href={href}>{children}</a>;
      }
    },
    ul: ({ children }) => <ul className={lstack({ gap: "2" })}>{children}</ul>,
  };

  const sourceList = values(sources);

  if (!isEnabled) {
    return (
      <UnreadyBanner error="Ask mode is not enabled for this installation." />
    );
  }

  return (
    <LStack>
      <styled.form
        className={hstack({ w: "full", gap: "0" })}
        onSubmit={fetchAnswer}
      >
        <Input
          type="text"
          w="full"
          value={question}
          onChange={(e) => setQuestion(e.target.value)}
          placeholder="Wisdom begins in wonder..."
          disabled={isLoading}
          borderRightRadius="none"
        />
        <Button
          type="submit"
          disabled={!question || isLoading}
          loading={isLoading}
          borderLeftRadius="none"
        >
          Ask
        </Button>
      </styled.form>
      <LStack>
        <ReactMarkdown className="typography" components={components}>
          {replaceSdrUrls(content)}
        </ReactMarkdown>

        {sourceList.length > 0 && (
          <LStack>
            <Heading>Sources from the community</Heading>
            {sourceList.map((source) => (
              <SourceCard
                key={source.id}
                href={source.href}
                kp={source.kp}
                id={source.id}
              />
            ))}
          </LStack>
        )}
      </LStack>
    </LStack>
  );
}

function SourceCard({
  href,
  kp,
  id,
}: {
  href: string;
  kp: string;
  id: string;
}) {
  switch (kp) {
    case "t":
      return <SourceCardThread href={href} id={id} />;
    case "l":
      return <SourceCardNode href={href} id={id} />;
    default:
      return <Link href={href}>{href}</Link>;
  }
}

function SourceCardThread({ href, id }: { href: string; id: string }) {
  const { error, data } = useThreadGet(id);
  if (!data) {
    if (error) {
      return <a href={href}>href</a>;
    }
    return (
      <Box display="inline">
        <Spinner />
      </Box>
    );
  }

  return (
    <DatagraphItemPostGenericCard
      item={{
        kind: "thread",
        ref: data,
      }}
    />
  );
}
function SourceCardNode({ href, id }: { href: string; id: string }) {
  const { error, data } = useNodeGet(id);
  if (!data) {
    if (error) {
      return <a href={href}>href</a>;
    }
    return (
      <Box display="inline">
        <Spinner />
      </Box>
    );
  }

  return (
    <DatagraphItemNodeCard
      item={{
        kind: "node",
        ref: data,
      }}
    />
  );
}

function sourceDataFromURL(href: string) {
  const url = new URL(href);

  const [_, kind, id] = url.pathname.split("/");

  if (!kind || !id) {
    throw new Error("Invalid SDR");
  }

  if (id.length != 20) {
    throw new Error("Invalid XID");
  }

  return [kind, id];
}

function getRouteForKind(kind: DatagraphItemKind): string {
  switch (kind) {
    case DatagraphItemKind.thread:
      return "t";
    case DatagraphItemKind.node:
      return "l";
    default:
      return "";
  }
}
