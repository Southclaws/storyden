import { OAuthConsentScreen } from "@/screens/auth/OAuthConsentScreen/OAuthConsentScreen";

export default function Page() {
  return <OAuthConsentScreen />;
}

export function generateMetadata() {
  return {
    title: "OAuth consent",
  };
}
