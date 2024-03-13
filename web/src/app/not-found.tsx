import ErrorBanner from "src/components/site/ErrorBanner";
import { Default } from "src/layouts/Default";

export default function Page() {
  return (
    <Default>
      <ErrorBanner message="The link to this page did not lead anywhere." />
    </Default>
  );
}
