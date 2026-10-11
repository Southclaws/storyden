import { nodeList } from "@/api/openapi-server/nodes";
import { UnreadyBanner } from "@/components/site/Unready";
import { nodeListCached } from "@/lib/library/server-node-list";
import { canonicalPage, jsonLD } from "@/lib/metadata/public";
import { LibraryIndexScreen } from "@/screens/library/LibraryIndexScreen";

export default async function Page() {
  let nodes;

  try {
    nodes = await nodeListCached();
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }

  let structuredData;

  try {
    const { data } = await nodeList(
      { visibility: ["published"], depth: "0" },
      { anonymous: true, cache: "no-store" },
    );

    structuredData = {
      "@context": "https://schema.org",
      "@type": "CollectionPage",
      name: "Library",
      url: canonicalPage("/l"),
      mainEntity: {
        "@type": "ItemList",
        itemListElement: data.nodes
          .filter(
            (node) =>
              node.visibility === "published" &&
              nodes.data.nodes.some(
                (visibleNode) => visibleNode.id === node.id,
              ),
          )
          .map((node, index) => ({
            "@type": "ListItem",
            position: index + 1,
            name: node.name,
            url: canonicalPage(`/l/${encodeURIComponent(node.slug)}`),
          })),
      },
    };
  } catch {
    // The signed-in library can still render without public structured data.
  }

  return (
    <>
      {structuredData && (
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: jsonLD(structuredData) }}
        />
      )}
      <LibraryIndexScreen nodes={nodes.data} />
    </>
  );
}
