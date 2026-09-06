import { flatten, zip } from "lodash";

import { getSettings } from "@/lib/settings/settings-server";
import { createThemeResourceResponse } from "@/lib/theme/theme-resource";
import { getServerThemeBundle } from "@/lib/theme/theme-server";
import { getColourVariants } from "@/utils/colour";

export async function GET(request: Request) {
  const [settings, theme] = await Promise.all([
    getSettings(),
    getServerThemeBundle(),
  ]);

  const cv = getColourVariants(settings.accent_colour);

  const rules = Object.entries(cv).map(([k, v]) => `${k}: ${v};`);

  const document = css`
    :root {
      ${rules.join("\n      ")}
    }

    ${theme.stylesheet}
  `;

  return createThemeResourceResponse(request, document, "text/css");
}

// NOTE: literally just so we get syntax highlighting above...
function css(s: TemplateStringsArray, ...v: string[]): string {
  return flatten(zip(s, v)).join("");
}
