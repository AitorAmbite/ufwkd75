import { app } from "./state.svelte";
import { fromHex, toHex, hsvToRgb, lerp, isLit, type RGB } from "./color";
import type { Key } from "../../bindings/d75studio/internal/d75/models";

export type Tool = "brush" | "picker" | "eraser";

const MAX_RECENT = 14;

export const paint = $state({
  tool: "brush" as Tool,
  color: "#7c6cff",
  gradA: "#ff2d95",
  gradB: "#29c5ff",
  painting: false,
  recent: JSON.parse(localStorage.getItem("d75.recent") || "[]") as string[],
});

export const ZONES: { id: "all" | "letters" | "numbers" | "frow" | "mods" | "nav"; test: (k: Key) => boolean }[] = [
  { id: "all", test: () => true },
  { id: "letters", test: (k) => /^Key[A-Z]$/.test(k.name) || k.name === "Semicolon" },
  { id: "numbers", test: (k) => /^Digit\d$/.test(k.name) || ["Backquote", "Minus", "Equal"].includes(k.name) },
  { id: "frow", test: (k) => /^F\d+$/.test(k.name) || k.name === "Escape" },
  { id: "mods", test: (k) => /Shift|Control|Meta|Alt|Fn|CapsLock|Tab|Enter|Backspace|Space/.test(k.name) },
  { id: "nav", test: (k) => /^Arrow|PageUp|PageDown|Insert|Delete/.test(k.name) },
];

export function rememberColor(hex: string) {
  paint.recent = [hex, ...paint.recent.filter((c) => c !== hex)].slice(0, MAX_RECENT);
  localStorage.setItem("d75.recent", JSON.stringify(paint.recent));
}

function custom() {
  const led = app.profile!.led;
  if (!led.custom) led.custom = {};
  return led.custom;
}

export function setKeyColor(name: string, rgb: RGB | null) {
  const c = custom();
  if (rgb && isLit(rgb)) c[name] = rgb;
  else delete c[name];
}

export const keyColor = (name: string): RGB | null => (app.profile?.led.custom?.[name] as RGB) ?? null;

/** Aplica la herramienta activa sobre una tecla. */
export function useTool(name: string, start: boolean) {
  if (paint.tool === "picker") {
    if (!start) return;
    const rgb = keyColor(name);
    if (rgb) paint.color = toHex(rgb);
    paint.tool = "brush";
    return;
  }
  if (start) {
    paint.painting = true;
    if (paint.tool === "brush") rememberColor(paint.color);
  }
  if (!paint.painting) return;
  setKeyColor(name, paint.tool === "eraser" ? null : fromHex(paint.color));
}

function keyPosition(k: Key) {
  const layout = app.meta!.layout!;
  const maxX = Math.max(...layout.map((k) => k.x + k.w));
  return (k.x + k.w / 2) / maxX;
}

export function fillZone(test: (k: Key) => boolean) {
  rememberColor(paint.color);
  const rgb = fromHex(paint.color);
  for (const k of app.meta!.layout!) if (test(k)) setKeyColor(k.name, rgb);
}

export function fillGradient() {
  const a = fromHex(paint.gradA), b = fromHex(paint.gradB);
  for (const k of app.meta!.layout!) setKeyColor(k.name, lerp(a, b, keyPosition(k)));
}

export function fillRainbow() {
  for (const k of app.meta!.layout!) setKeyColor(k.name, hsvToRgb(keyPosition(k) * 300, 1, 1));
}

export function clearAll() {
  app.profile!.led.custom = {};
}

export { keyPosition };
