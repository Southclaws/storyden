import { flatten, zip } from "lodash";
import { NextResponse } from "next/server";

import { getInfo } from "src/utils/info";

/**
 *
 * @returns A fully static server-side-rendered CSS document that uses the
 *          Storyden installation's accent colour set by the administrator.
 */

export async function GET() {
  const info = await getInfo();

  const document = css`
    :root {
      --accent-colour: ${info.accent_colour};
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
