import { UnreadyBanner } from "src/components/site/Unready";
import { Default } from "src/layouts/Default";

export default function Page() {
  return (
    <Default>
      <UnreadyBanner error="The link to this page did not lead anywhere." />
    </Default>
  );
}
