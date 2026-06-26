import { RobotConfigurationScreen } from "@/screens/robots/RobotConfigurationScreen";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

type Props = {
  params: Promise<{
    robot: string;
  }>;
};

export default async function Page(props: Props) {
  const { robot } = await props.params;

  return <RobotConfigurationScreen robotId={robot} />;
}
