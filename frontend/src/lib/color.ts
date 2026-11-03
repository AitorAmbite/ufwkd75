export type RGB = number[];

export const toHex = (rgb: RGB) => "#" + rgb.map((v) => v.toString(16).padStart(2, "0")).join("");

export const fromHex = (hex: string): RGB => [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16));

export const isHex = (s: string) => /^#[0-9a-f]{6}$/i.test(s);

export const lerp = (a: RGB, b: RGB, t: number): RGB => a.map((v, i) => Math.round(v + (b[i] - v) * t));

export function hsvToRgb(h: number, s: number, v: number): RGB {
  const f = (n: number) => {
    const k = (n + h / 60) % 6;
    return Math.round(255 * (v - v * s * Math.max(0, Math.min(k, 4 - k, 1))));
  };
  return [f(5), f(3), f(1)];
}

export const isLit = (rgb?: RGB | null) => !!rgb && rgb.some((v) => v > 0);

export const luminance = (rgb: RGB) => 0.299 * rgb[0] + 0.587 * rgb[1] + 0.114 * rgb[2];

/** Escala un color por el brillo del teclado (0-4) para la vista previa. */
export const dim = (rgb: RGB, brightness: number): RGB => {
  const f = brightness === 0 ? 0.12 : 0.35 + (brightness / 4) * 0.65;
  return rgb.map((v) => Math.round(v * f));
};
