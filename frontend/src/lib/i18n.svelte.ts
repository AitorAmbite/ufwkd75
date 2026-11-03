import { en, enCodes, enKeys, enEffects, type Dict } from "./locales/en";
import { es, esCodes, esKeys, esEffects } from "./locales/es";
import { zh, zhCodes, zhKeys, zhEffects } from "./locales/zh";

export type Locale = "es" | "en" | "zh";

export const LOCALES: { id: Locale; name: string }[] = [
  { id: "es", name: "Español" },
  { id: "en", name: "English" },
  { id: "zh", name: "中文" },
];

const tables = {
  es: { dict: es, codes: esCodes, keys: esKeys, effects: esEffects },
  en: { dict: en, codes: enCodes, keys: enKeys, effects: enEffects },
  zh: { dict: zh, codes: zhCodes, keys: zhKeys, effects: zhEffects },
};

function detect(): Locale {
  const saved = localStorage.getItem("d75.locale");
  if (saved === "es" || saved === "en" || saved === "zh") return saved;
  const lang = navigator.language.toLowerCase();
  if (lang.startsWith("zh")) return "zh";
  if (lang.startsWith("es")) return "es";
  return "en";
}

export const i18n = $state({ locale: detect() });
document.documentElement.lang = i18n.locale;

export function setLocale(locale: Locale) {
  i18n.locale = locale;
  localStorage.setItem("d75.locale", locale);
  document.documentElement.lang = locale;
}

/** Traduce una clave; {nombre} se sustituye por params.nombre. */
export function t(key: keyof Dict, params: Record<string, string | number> = {}): string {
  const text = tables[i18n.locale].dict[key] ?? en[key] ?? key;
  return text.replace(/\{(\w+)\}/g, (_, k) => String(params[k] ?? `{${k}}`));
}

/** Etiqueta de una función asignable; si no hay traducción se usa la del backend. */
export const tCode = (code: number, fallback: string) => tables[i18n.locale].codes[code] ?? fallback;

/** Leyenda de una tecla física; los símbolos de la distribución no se traducen. */
export const tKey = (name: string, fallback: string) => tables[i18n.locale].keys[name] ?? fallback;

export const tEffect = (index: number, fallback: string) => tables[i18n.locale].effects[index] ?? fallback;

export const tGroup = (group: string) => t(`group.${group}` as keyof Dict);

/** Traduce los errores con código del backend ("código: detalle"). */
export function tError(message: string): string {
  const m = /^([a-z_]+)(?::\s*([\s\S]*))?$/.exec(message.trim());
  if (!m) return message;
  const key = `err.${m[1]}` as keyof Dict;
  if (!(key in en)) return message;
  return m[2] ? `${t(key)}: ${m[2]}` : t(key);
}
