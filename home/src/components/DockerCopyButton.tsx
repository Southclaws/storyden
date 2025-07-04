"use client";

import { styled } from "@/styled-system/jsx";
import { useCopyToClipboard, useIsClient } from "@uidotdev/usehooks";
import { CopyIcon } from "lucide-react";
import { useState } from "react";

import "./hearts.css";

type heart = {
  id: number;
  x: string;
  y: string;
  char: string;
};

const floatyEmojis = [
  // Very common
  "💖",
  "💖",
  "❤️",
  "❤️",
  "✨",
  "✨",
  "💕",
  "💕",
  // Common
  "💗",
  "💗",
  "💓",
  "💓",
  "⭐",
  "🌟",
  // Storyden!
  "ᛟ",
  "ᛟ",
  "ᛟ",
  "ᛟ",
  "ᛟ",
  // Mid
  "💞",
  "💘",
  "💝",
  "🩷",
  "🌈",
  "☁️",
  "🌸",
  "🎉",
  "🎀",
  // Rare/fun
  "🦄",
  "🪽",
  "👻",
  "💫",
  "🍄",
  "🎈",
  "🌙",
  "🔮",
  "🕊️",
  "🦋",
  // Ultra-rare easter eggs
  "💀",
  "🦆",
  "🥹",
  "🧃",
  "🥲",
  "🌚",
  "🧸",
];
const numFloatyEmojis = floatyEmojis.length;

export function DockerCopyButton() {
  const [copiedText, copyToClipboard] = useCopyToClipboard();
  const hasCopiedText = Boolean(copiedText);

  const isClient = useIsClient();

  const [hearts, setHearts] = useState<heart[]>([]);

  function handleCopy() {
    copyToClipboard("docker run -p 8000:8000 ghcr.io/southclaws/storyden");

    const slots = 5;

    if (!(window as any).__storydenConsoleEasterEggShown) {
      (window as any).__storydenConsoleEasterEggShown = true;

      console.log(
        "%c ᛟ  Welcome to the homestead. Enjoy your corner of the web!",
        "font-size: 1.5rem; font-weight: bold; color: #d8dbcd; background: linear-gradient(90deg, #307343, #104059); padding: 8px 12px; border-radius: 6px;"
      );
      console.log("You're curious, what are you looking for?");
      console.log(
        "docs for your brand new copied docker command? https://www.storyden.org/docs/introduction/vps/docker"
      );
      console.log(
        "feeling contribute-y? check out hub de la git: https://github.com/Southclaws/storyden"
      );
      console.log(
        "%c ᛟ stay fresh",
        "font-style: italic; color: #854627; background: #d68e4d; padding: 4px 8px; border-radius: 4px;"
      );
    }

    for (const heart of [1, 2, 3, 4, 5]) {
      const base = (heart - 0.5) / slots;
      const jitter = (Math.random() - 0.5) * (25 / window.innerWidth);

      setTimeout(() => {
        setHearts((h) => [
          ...h,
          {
            id: Math.random(),
            x: `${Math.min(1, Math.max(0, base + jitter)) * 100}%`,
            y: `-${49 + Math.random()}px`,
            char: floatyEmojis[(Math.random() * numFloatyEmojis) | 0],
          },
        ]);

        setTimeout(
          () =>
            setHearts((h) =>
              h.filter((i) => {
                return !h.some((n) => n.id === i.id);
              })
            ),
          1000
        );
      }, heart * 100);
    }
  }

  return (
    <styled.pre
      position="relative"
      display="flex"
      alignItems="center"
      bgColor="Shades.iron"
      py="1"
      pl="2"
      pr="1"
      borderRadius="md"
      gap="1"
    >
      {hearts.map(({ id, x, y, char }) => (
        <styled.span
          _motionReduce={{
            display: "none",
          }}
          key={id}
          style={{
            position: "absolute",
            top: y,
            left: x,
            animation: "floatUp 1s ease-out",
            pointerEvents: "none",
            fontSize: "2rem",
          }}
        >
          {char}
        </styled.span>
      ))}
      <styled.span py="1">
        docker run -p 8000:8000 ghcr.io/southclaws/storyden
      </styled.span>

      {isClient && (
        <styled.button
          type="button"
          p="1"
          borderRadius="sm"
          cursor="pointer"
          title="Copy, run, get Storydenning in seconds!"
          _hover={{
            bgColor: "Shades.stone/50",
          }}
          aria-label="copy docker command to clipboard"
          onClick={handleCopy}
        >
          <CopyIcon width="16" height="16" />
        </styled.button>
      )}
    </styled.pre>
  );
}
