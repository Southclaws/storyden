import type { Preview, ReactRenderer } from "@storybook/nextjs-vite";
import { useEffect } from "react";
import { Toaster } from "sonner";
import type { DecoratorFunction } from "storybook";
import { SWRConfig } from "swr";

import "../src/app/global.css";
import { FALLBACK_COLOUR, getColourVariants } from "../src/utils/colour";

import "./preview.css";

const storydenThemeCss = `:root {
  ${Object.entries(getColourVariants(FALLBACK_COLOUR))
    .map(([property, value]) => `${property}: ${value};`)
    .join("\n  ")}
}`;

const withStorydenProviders: DecoratorFunction<ReactRenderer> = (
  Story,
  context,
) => {
  useEffect(() => {
    const storydenWindow = window as typeof window & {
      __storyden__?: {
        API_ADDRESS: string;
        WEB_ADDRESS: string;
        source: string;
      };
    };

    storydenWindow.__storyden__ = {
      API_ADDRESS: "http://localhost:8000",
      WEB_ADDRESS: "http://localhost:6006",
      source: "storybook",
    };
  }, []);

  const layout = context.parameters.layout ?? "padded";

  return (
    <>
      <style>{storydenThemeCss}</style>
      <div className="storybook-preview">
        <SWRConfig value={{ provider: () => new Map() }}>
          <div className="storybook-preview__story" data-layout={layout}>
            <Story />
          </div>
          <Toaster />
        </SWRConfig>
      </div>
    </>
  );
};

const preview: Preview = {
  decorators: [withStorydenProviders],
  parameters: {
    layout: "padded",
    nextjs: {
      appDirectory: true,
    },
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
    a11y: {
      test: "todo",
    },
  },
};

export default preview;
