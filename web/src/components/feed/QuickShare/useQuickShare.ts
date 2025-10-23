import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { linkCreate } from "@/api/openapi-client/links";
import { Account, Category, LinkReference } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { NO_CATEGORY_VALUE } from "@/components/category/CategorySelect/useCategorySelect";
import { useFeedMutations } from "@/lib/feed/mutation";
import { useClickAway } from "@/utils/useClickAway";

export type Props = {
  initialSession?: Account;
  initialCategory?: Category | null;
  showCategorySelect: boolean;
};

export const FormSchema = z.object({
  body: z.string(),
  category: z.string().optional(),
});
export type Form = z.infer<typeof FormSchema>;

export function useQuickShare({ initialCategory }: Props) {
  const session = useSession();
  const [editing, setEditing] = useState(false);
  const [postURL, setPostURL] = useState<string | null>(null);
  const [hydratedLink, setHydratedLink] = useState<
    LinkReference | "loading" | null
  >(null);
  const formRef = useClickAway<HTMLFormElement>(() => setEditing(false));
  const [resetKey, setResetKey] = useState("");

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      body: undefined,
      category: initialCategory?.id,
    },
  });

  // Watch body for changes - because onChange is already used by RHF.
  const bodyContent = form.watch("body");

  // When the body changes, find the first URL in the content and remember it.
  useEffect(() => {
    const parsed = new DOMParser().parseFromString(bodyContent, "text/html");
    const url = getFirstURL(parsed);

    if (url) {
      setPostURL(url);
    } else {
      setPostURL(null);
      setHydratedLink(null);
    }
  }, [bodyContent]);

  // When the URL found in the body changes, index/fetch its rich preview.
  useEffect(() => {
    if (postURL) {
      setHydratedLink("loading");
      linkCreate({ url: postURL }).then((link) => {
        // Only hydrate link if there's a preview available. Otherwise bail out.
        if (link.title && link.description) {
          setHydratedLink(link);
        } else {
          setHydratedLink(null);
        }
      });
    }
  }, [postURL]);

  const { createThread, revalidate } = useFeedMutations(session, {
    categories:
      initialCategory === undefined
        ? undefined
        : initialCategory === null
          ? ["null"]
          : [initialCategory.slug],
  });

  const handlePost = form.handleSubmit((data: Form) => {
    handle(
      async () => {
        const parsed = new DOMParser().parseFromString(
          bodyContent,
          "text/html",
        );
        const { title, body, isFallback } = splitTitleBody(parsed);

        const linkAvailable =
          hydratedLink && hydratedLink !== "loading" ? hydratedLink : undefined;

        const threadTitle = isFallback
          ? linkAvailable?.title
            ? linkAvailable.title
            : title
          : title;

        if (!title) {
          throw new Error("Cannot post an empty thread.");
        }

        const newThread = {
          title: threadTitle,
          body,
          url: postURL ?? undefined,
          category:
            data.category === NO_CATEGORY_VALUE ? undefined : data.category,
          visibility: "published" as const,
        };

        await createThread(newThread, linkAvailable);

        // Awful hack to reset the rich text editor...
        setResetKey(new Date().toISOString());

        // Only reset the body content, the member might want to post again.
        form.resetField("body");

        setHydratedLink(null);
      },
      {
        async cleanup() {
          revalidate();
        },
      },
    );
  });

  function handleFocus() {
    setEditing(true);
  }

  return {
    form,
    state: {
      formRef,
      editing,
      hydratedLink,
      resetKey,
    },
    handlers: {
      handleFocus,
      handlePost,
    },
  };
}

function getFirstURL(html: Document) {
  const result = html.querySelector("a");
  if (!result) {
    return undefined;
  }

  const href = result?.attributes.getNamedItem("href")?.value;
  if (!href?.startsWith("https")) {
    return undefined;
  }

  return href;
}

function splitTitleBody(html: Document) {
  const bodyEl = html.querySelector("body");

  // The title is the first text node of the first paragraph element. If the
  // paragraph contains any more tags after the initial text node, do not include
  // them in the title, this ensures that <a> tags are not included.
  const firstChild = bodyEl?.querySelector("p")?.childNodes[0] as Node;

  if (!firstChild) {
    throw new Error("Not enough text content to post a new thread.");
  }

  let isFallback = false;

  const title = ((): string => {
    const textContent = firstChild.textContent;

    switch (firstChild.nodeType) {
      case Node.ELEMENT_NODE:
        if (firstChild.nodeName === "A") {
          // Mark this as a fallback strategy, if the Link acquired in the
          // actual link fetch yields a title from the opengraph data, and this
          // is true, then we'll use the opengraph title instead.
          isFallback = true;

          if (textContent?.startsWith("http")) {
            // If the anchor tag is just a bare tag (where the text content is
            // the actual URL itself, no title), then we'll just use that.
            const parsed = new URL(textContent);
            return parsed.hostname;
          } else if (textContent === "") {
            // Otherwise, if there's no text content in the node, try to get the
            // href attribute and use the hostname from that as the title.
            const href = (firstChild as any) /* Types are broken here */
              ?.getAttribute("href");

            if (href) {
              const parsed = new URL(href);
              return parsed.hostname;
            }

            return textContent ?? "";
          } else {
            // Finally, if none of the above conditions are met, just use the
            // text content of the node. Which might be empty, but caught below.
            return textContent ?? "";
          }
        }
    }

    return textContent ?? "";
  })();

  // We want something of substance to post a new thread. This may bug out tho.
  if (!title) {
    throw new Error("Not enough text content to post a new thread.");
  }

  // Now remove the first text node from the first paragraph of the bodyEl.
  bodyEl?.querySelector("p")?.childNodes[0]?.remove();

  const body = bodyEl?.getHTML() ?? "";

  return {
    title,
    body,
    isFallback,
  };
}
