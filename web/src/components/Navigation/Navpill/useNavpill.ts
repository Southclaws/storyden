import { useNavigation } from "../useNavigation";

export function useNavpill() {
  const navigation = useNavigation();

  return navigation;
}
