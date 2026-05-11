export type TranslationParams = Record<string, string | number>;
export type Translate = (key: string, params?: TranslationParams) => string;

export function interpolate(message: string, params?: TranslationParams) {
  if (!params) {
    return message;
  }

  return message.replace(/\{\{(\w+)\}\}/g, (match, name) => {
    const value = params[name];
    return value === undefined ? match : String(value);
  });
}
