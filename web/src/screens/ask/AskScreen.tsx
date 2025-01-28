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
  kind: string;
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

type References = {
  refs: DatagraphRef[];
  urls: string[];
};

export function Ask() {
  const [question, setQuestion] = useState("");
  const [content, setContent] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [sources, setSources] = useState<References>({
    refs: [],
    urls: [],
  });
  const isEnabled = useCapability("semdex");

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
        setSources({
          refs: [],
          urls: [],
        });

        const source = new EventSource(
          `${API_ADDRESS}/api/datagraph/ask?q=${encodeURIComponent(question)}`,
        );

        source.onerror = (err) => {
          source.close();
        };

        source.addEventListener("text", (e) => {
          const data = JSON.parse(e.data);
          const { chunk } = data;
          if (chunk) {
            setContent((prev) => prev + chunk);
          }
        });

        source.addEventListener("meta", (e) => {
          const data = JSON.parse(e.data);
          const { refs, urls } = data;

          setSources({ refs, urls });
        });
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

        {sources.refs.length > 0 && (
          <LStack>
            <Heading>Sources from the community</Heading>
            {sources.refs.map((source) => (
              <SourceCard key={source.id} kind={source.kind} id={source.id} />
            ))}
          </LStack>
        )}

        {sources.urls.length > 0 && (
          <LStack>
            <Heading>Sources from the web</Heading>
            <div className="typography">
              <ul>
                {sources.urls.map((url) => (
                  <li key={url}>
                    <Link className="link" href={url}>
                      {url}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          </LStack>
        )}
      </LStack>
    </LStack>
  );
}

function SourceCard({ kind, id }: { kind: string; id: string }) {
  const kpath = getRouteForKind(kind as any);
  const href = `${WEB_ADDRESS}/${kpath}/${id}`;
  switch (kind) {
    case "thread":
      return <SourceCardThread href={href} id={id} />;
    case "node":
      return <SourceCardNode href={href} id={id} />;
    default:
      return <Link href={href}>{href}</Link>;
  }
}

function SourceCardThread({ href, id }: { href: string; id: string }) {
  const { error, data } = useThreadGet(id);
  if (!data) {
    if (error) {
      return null;
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
      return null;
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
