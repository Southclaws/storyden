import { flatten, zip } from "lodash";
import { NextResponse } from "next/server";

import { getColourVariants } from "src/utils/colour";

import { getSettings } from "@/lib/settings/settings-server";

/**
 *
 * @returns A fully static server-side-rendered CSS document that uses the
 *          Storyden installation's accent colour set by the administrator.
 */

export async function GET() {
  const settings = await getSettings();

  const cv = getColourVariants(settings.accent_colour);

  const rules = Object.entries(cv).map(([k, v]) => `${k}: ${v};`);

  const document = css`
    :root {
      ${rules.join("\n      ")}
    }
  `;

  return new NextResponse(document, {
    headers: {
      "Content-Type": "text/css",
    },
  });
}

// NOTE: literally just so we get syntax highlighting above...
function css(s: TemplateStringsArray, ...v: string[]): string {
  return flatten(zip(s, v)).join("");
}
