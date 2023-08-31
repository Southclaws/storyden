import { Content } from "./Content/Content";

export async function HomeScreen() {
  return (
    <>
      {/* @ts-expect-error Server Component */}
      <Content showEmptyState />
    </>
  );
}
