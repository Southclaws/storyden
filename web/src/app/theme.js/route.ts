import { createThemeResourceResponse } from "@/lib/theme/theme-resource";
import { getServerThemeBundle } from "@/lib/theme/theme-server";

export async function GET(request: Request) {
  const theme = await getServerThemeBundle();

  return createThemeResourceResponse(
    request,
    theme.script,
    "application/javascript",
  );
}
