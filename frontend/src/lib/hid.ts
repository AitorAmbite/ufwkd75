import type { Key } from "../../bindings/d75studio/internal/d75/models";

// Teclas que no están en el D75 pero se pueden grabar desde otro teclado (KeyboardEvent.code → usage HID)
const extra: Record<string, number> = {
  IntlBackslash: 0x64, PrintScreen: 0x46, ScrollLock: 0x47, Pause: 0x48, Home: 0x4a, End: 0x4d,
  ContextMenu: 0x65, ControlRight: 0xe4, MetaRight: 0xe7, NumLock: 0x53, NumpadDivide: 0x54,
  NumpadMultiply: 0x55, NumpadSubtract: 0x56, NumpadAdd: 0x57, NumpadEnter: 0x58, NumpadDecimal: 0x63,
  Numpad0: 0x62, Numpad1: 0x59, Numpad2: 0x5a, Numpad3: 0x5b, Numpad4: 0x5c, Numpad5: 0x5d,
  Numpad6: 0x5e, Numpad7: 0x5f, Numpad8: 0x60, Numpad9: 0x61,
};
for (let i = 13; i <= 24; i++) extra[`F${i}`] = 0x68 + i - 13;

/** Usage HID de una tecla según KeyboardEvent.code (posición física, no depende de la distribución). */
export function hidFromEvent(code: string, layout: Key[]): number | undefined {
  const k = layout.find((k) => k.name === code);
  return k && k.default > 0 && k.default < 0xe8 ? k.default : extra[code];
}

/** Grupos de funciones que el firmware acepta dentro de una macro (usages HID de teclado). */
export const MACRO_GROUPS = ["letters", "numbers", "symbols", "function", "editing", "navigation", "modifiers", "numpad"];
