"use client";

import { cva } from "@/styled-system/css";
import { HStack, styled } from "@/styled-system/jsx";
import { Copy, Rss } from "lucide-react";
import { useState } from "react";

const buttonGroupStyles = cva({
  base: {
    display: "flex",
    alignItems: "center",
    gap: "1.5",
    px: "1.5",
    py: "0.5",
    bg: "Shades.newspaper/60",
    color: "Primary.forest",
    border: "1px solid",
    borderColor: "Shades.stone/80",
    fontSize: "xs",
    fontWeight: "medium",
    cursor: "pointer",
    textDecoration: "none",
    _hover: {
      bg: "Shades.stone/60",
      color: "Mono.ink",
    },
    transition: "all 0.2s",
  },
  variants: {
    position: {
      left: {
        borderRadius: "lg",
        borderRightRadius: "0",
        borderRightWidth: "0",
      },
      right: {
        borderRadius: "lg",
        borderLeftRadius: "0",
      },
    },
    state: {
      default: {},
      success: {
        bg: "Primary.campfire",
        color: "Mono.ink",
        borderColor: "Primary.saddle",
        _hover: {
          bg: "Primary.saddle",
          color: "Mono.slush",
        },
      },
    },
  },
});

export function RSSButtons() {
  const [copied, setCopied] = useState(false);

  const rssUrl = "https://www.storyden.org/rss.xml";

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(rssUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy RSS URL:", err);
    }
  };

  return (
    <HStack gap="0">
      <styled.a
        href={rssUrl}
        target="_blank"
        rel="noopener noreferrer"
        className={buttonGroupStyles({
          position: "left",
          state: "default",
        })}
      >
        <Rss size={16} />
        RSS Feed
      </styled.a>

      <styled.button
        onClick={handleCopy}
        className={buttonGroupStyles({
          position: "right",
          state: copied ? "success" : "default",
        })}
      >
        <Copy size={16} />
        {copied ? "Copy URL" : "Copy URL"}
      </styled.button>
    </HStack>
  );
}
