import { notFound } from "next/navigation";

// Complete the children slot when the sidebar catch-all matches a sibling route
// Unknown dashboard URLs must still render the not-found page.
export default function Default() {
  notFound();
}
