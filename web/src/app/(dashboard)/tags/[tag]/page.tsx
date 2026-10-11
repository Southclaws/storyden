import { tagGet } from "@/api/openapi-server/tags";
import { UnreadyBanner } from "@/components/site/Unready";
import { canonicalPage, jsonLD } from "@/lib/metadata/public";
import { TagScreen } from "@/screens/tags/TagScreen";

type Props = {
  params: Promise<{
    tag: string;
  }>;
};

export default async function Page(props: Props) {
  const { tag } = await props.params;
  let data;

  try {
    ({ data } = await tagGet(tag));
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }

  let structuredData;

  try {
    const { data: publicTag } = await tagGet(tag, {
      anonymous: true,
      cache: "no-store",
    });
    const items = publicTag.items.filter(
      (item) =>
        (item.kind === "profile" || item.ref.visibility === "published") &&
        data.items.some((visibleItem) => visibleItem.ref.id === item.ref.id),
    );

    structuredData = {
      "@context": "https://schema.org",
      "@type": "CollectionPage",
      name: publicTag.name,
      url: canonicalPage(`/tags/${encodeURIComponent(publicTag.name)}`),
      mainEntity: {
        "@type": "ItemList",
        itemListElement: items.map((item, index) => {
          let path;
          let name;

          switch (item.kind) {
            case "profile":
              path = `/m/${encodeURIComponent(item.ref.handle)}`;
              name = item.ref.name;
              break;

            case "node":
              path = `/l/${encodeURIComponent(item.ref.slug)}`;
              name = item.ref.name;
              break;

            case "reply":
              path = `/t/locate/${encodeURIComponent(item.ref.id)}`;
              name = item.ref.title;
              break;

            case "post":
            case "thread":
              path = `/t/${encodeURIComponent(item.ref.slug)}`;
              name = item.ref.title;
              break;
          }

          return {
            "@type": "ListItem",
            position: index + 1,
            name,
            url: canonicalPage(path),
          };
        }),
      },
    };
  } catch {
    // The signed-in tag can still render without public structured data.
  }

  return (
    <>
      {structuredData && (
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: jsonLD(structuredData) }}
        />
      )}
      <TagScreen initialTag={data} slug={tag} />
    </>
  );
}
