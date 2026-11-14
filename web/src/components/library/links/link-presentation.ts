import type { LinkReference } from "@/api/openapi-schema";
import type { StatusBadgeTone } from "@/components/ui/status-badge";

type LinkStatus = {
  label: string;
  tone: StatusBadgeTone;
};

/** The compact hint prioritises reachability; saved metadata remains usable. */
export function getLinkStatus(link: LinkReference): LinkStatus | undefined {
  const hasSavedPreview = Boolean(
    link.scraped_at ||
    link.primary_image ||
    link.favicon_image ||
    link.title ||
    link.description,
  );

  if (link.liveness_status === "offline") {
    return { label: "Couldn't reach site", tone: "warning" };
  }

  if (link.scrape_status === "pending" && !hasSavedPreview) {
    return { label: "Preview pending", tone: "neutral" };
  }

  if (link.scrape_status === "failed") {
    if (hasSavedPreview) {
      return { label: "Preview refresh failed", tone: "neutral" };
    }

    return { label: "Preview unavailable", tone: "neutral" };
  }

  return undefined;
}
