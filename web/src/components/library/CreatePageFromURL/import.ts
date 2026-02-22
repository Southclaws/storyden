import { linkCreate } from "@/api/openapi-client/links";
import {
  nodeCreate,
  nodeGenerateContent,
  nodeGenerateTags,
  nodeGenerateTitle,
} from "@/api/openapi-client/nodes";
import { Asset, Node } from "@/api/openapi-schema";
import { deriveError } from "@/utils/error";
import { generateXid } from "@/utils/xid";

export type ImportStep =
  | "fetching_link"
  | "generating_tags"
  | "generating_title"
  | "generating_content"
  | "creating_node"
  | "complete"
  | "failed";

export const importStateLabel: Record<ImportStep, string> = {
  fetching_link: "Fetching page...",
  generating_tags: "Generating tags",
  generating_title: "Generating title",
  generating_content: "Generating content",
  creating_node: "Creating page",
  complete: "Complete",
  failed: "Failed to import",
};

export type ImportState = {
  step: ImportStep;
  data?: {
    title?: string;
    description?: string;
    primary_image?: Asset;
    tag_suggestions?: string[];
    title_suggestion?: string;
    content_suggestion?: string;
    created_node?: Node;
  };
  error?: string;
};

export type ImportOptions = {
  url: string;
  parentSlug?: string;
  genaiAvailable: boolean;
};

export async function* importFromURLGenerator({
  url,
  parentSlug,
  genaiAvailable,
}: ImportOptions): AsyncGenerator<ImportState, ImportState, unknown> {
  try {
    // Step 1: Fetch link data
    yield {
      step: "fetching_link",
    };

    const { title, description, primary_image } = await linkCreate({ url });

    const baseData = {
      title,
      description,
      primary_image,
    };

    yield {
      step: "fetching_link",
      data: baseData,
    };

    const finalData = {
      ...baseData,
      title_suggestion: undefined as string | undefined,
      tag_suggestions: [] as string[],
      content_suggestion: undefined as string | undefined,
    };

    // If GenAI is available and we have description, generate suggestions
    if (genaiAvailable && description) {
      // Step 2: Generate tags
      yield {
        step: "generating_tags",
        data: { ...finalData },
      };

      // Step 3: Generate title
      yield {
        step: "generating_title",
        data: { ...finalData },
      };

      // Step 4: Generate content
      yield {
        step: "generating_content",
        data: { ...finalData },
      };

      // We need a temporary slug to call the AI functions
      // This is a limitation of the current API design
      // The slug isn't actually used by the backend lol...
      const tempSlug = `temp-${Date.now()}`;

      // Enhance the content with metadata for better AI suggestions
      const enhancedContent = `URL: ${url}
Original Title: ${title || "N/A"}

Content:
${description}`;

      try {
        const [tag_suggestions, title_suggestion, content_suggestion] =
          await Promise.all([
            nodeGenerateTags(tempSlug, { content: enhancedContent })
              .then((r) => r.tags)
              .catch(() => undefined),
            nodeGenerateTitle(tempSlug, { content: enhancedContent })
              .then((r) => r.title)
              .catch(() => undefined),
            nodeGenerateContent(tempSlug, { content: enhancedContent })
              .then((r) => r.content)
              .catch(() => undefined),
          ]);

        if (tag_suggestions?.length) {
          finalData.tag_suggestions = tag_suggestions;
        }

        if (title_suggestion) {
          finalData.title_suggestion = title_suggestion;
        }

        if (content_suggestion) {
          finalData.content_suggestion = content_suggestion;
        }

        yield {
          step: "generating_content",
          data: finalData,
        };
      } catch (aiError) {
        // If AI fails, continue with non-AI data
        yield {
          step: "generating_content",
          data: finalData,
          error:
            aiError instanceof Error ? aiError.message : "AI generation failed",
        };
      }
    }

    // Step 5: Create the node
    yield {
      step: "creating_node",
      data: finalData,
    };

    const name = finalData.title_suggestion || finalData.title || "";
    const content = finalData.content_suggestion || finalData.description;

    // Generate a slug if there's no name to fill. We do this because the API
    // cannot generate a slug from an empty name. It probably should though.
    const slug = name === "" ? generateXid() : undefined;

    const created_node = await nodeCreate({
      name,
      slug,
      parent: parentSlug,
      description: finalData.description,
      primary_image_asset_id: finalData.primary_image?.id,
      tags: finalData.tag_suggestions,
      content,
      url,
    });

    // Step 6: Complete
    const completeState: ImportState = {
      step: "complete",
      data: {
        ...finalData,
        created_node,
      },
    };

    yield completeState;
    return completeState;
  } catch (error) {
    console.log("generator error:", error);

    const errorState: ImportState = {
      step: "failed",
      error: deriveError(error),
    };

    yield errorState;
    return errorState;
  }
}
