export function useAvatar(handle: string) {
  return { src: `/api/accounts/${handle}/avatar` };
}
