import { Events } from "@wailsio/runtime";
import { SvelteSet } from "svelte/reactivity";
import { KeyboardService, ProfileService } from "../../bindings/d75studio";
import type { Assign, KeyEvent, Macro, Meta, Profile, Status } from "../../bindings/d75studio/internal/d75/models";
import { t, tCode, tError } from "./i18n.svelte";

export type Page = "keys" | "light" | "macros" | "settings";

export const MODS: [number, string][] = [[1, "Ctrl"], [2, "Shift"], [4, "Alt"], [8, "Win"]];
export const layerName = (i: number) => t(`layer.${i}` as "layer.0");

export const app = $state({
  meta: null as Meta | null,
  profile: null as Profile | null,
  page: "keys" as Page,
  layer: 0,
  selected: null as string | null,
  status: { connected: false, needsPermission: false, detail: "" } as Status,
  profiles: [] as string[],
  profileName: "",
  appliedJSON: "",
  savedJSON: "",
  detect: false,
  busy: false,
});

/** Posiciones de matriz pulsadas ahora mismo en el teclado físico. */
export const pressed = new SvelteSet<number>();

const json = (p: Profile | null) => JSON.stringify($state.snapshot(p));

export const isApplied = () => json(app.profile) === app.appliedJSON;
export const isSaved = () => app.profileName !== "" && json(app.profile) === app.savedJSON;

// ---------- avisos y diálogos ----------

export type Toast = { id: number; text: string; kind: "ok" | "error" | "info" };
export const toasts = $state<Toast[]>([]);
let toastSeq = 0;

export function toast(text: string, kind: Toast["kind"] = "ok") {
  const id = ++toastSeq;
  toasts.push({ id, text, kind });
  setTimeout(() => {
    const i = toasts.findIndex((t) => t.id === id);
    if (i >= 0) toasts.splice(i, 1);
  }, kind === "error" ? 6000 : 2800);
}

export type Dialog = {
  title: string;
  message?: string;
  input?: string;
  confirm: string;
  danger?: boolean;
  resolve: (v: string | boolean | null) => void;
};
export const dialog = $state<{ current: Dialog | null }>({ current: null });

export function ask(opts: Omit<Dialog, "resolve">): Promise<string | boolean | null> {
  return new Promise((resolve) => (dialog.current = { ...opts, resolve }));
}

export function errorText(e: unknown): string {
  const msg = e && typeof e === "object" && "message" in e ? String((e as Error).message) : String(e);
  return tError(msg);
}

// ---------- teclas ----------

export const keyByName = (name: string) => app.meta!.layout!.find((k) => k.name === name)!;

export function codeLabel(code: number): string {
  const m = macroAt(code);
  if (m) return t("macros.keyLabel", { name: m.name });
  const c = app.meta!.codes!.find((c) => c.code === code);
  return c ? tCode(c.code, c.label) : "0x" + code.toString(16).padStart(2, "0");
}

export function assignLabel(a: Assign): string {
  const mods = MOD_NAMES(a.mods ?? 0);
  return [...mods, codeLabel(a.code)].join(" + ");
}

export const MOD_NAMES = (mods: number) => MODS.filter(([bit]) => mods & bit).map(([, n]) => n);

export function layerMap(layer = app.layer) {
  const p = app.profile!;
  if (!p.layers[layer]) p.layers[layer] = {};
  return p.layers[layer]!;
}

export function assign(name: string, code: number, mods: number) {
  const key = keyByName(name);
  const map = layerMap();
  // Asignar a una tecla su propio valor de fábrica equivale a no tocarla
  if (code === key.default && mods === 0) delete map[name];
  else map[name] = mods ? { code, mods } : { code };
}

export function resetKey(name: string, layer = app.layer) {
  delete layerMap(layer)[name];
}

// ---------- macros ----------

export function macros(): Macro[] {
  const p = app.profile!;
  if (!p.macros) p.macros = [];
  return p.macros;
}

/** Índice de la macro que asigna un código del keymap (0xF0 + n), o -1. */
export function macroIndex(code: number): number {
  const n = code - (app.meta?.macroCodeBase ?? 0xf0);
  return n >= 0 && n < (app.meta?.maxMacros ?? 10) ? n : -1;
}

const macroAt = (code: number) => {
  const n = macroIndex(code);
  return n >= 0 ? app.profile?.macros?.[n] : undefined;
};

export const macroCode = (n: number) => app.meta!.macroCodeBase + n;

export function addMacro(): number {
  const list = macros();
  list.push({ name: t("macros.defaultName", { n: list.length + 1 }), mode: "count", count: 1, events: [] });
  return list.length - 1;
}

/** Teclas (capa y nombre) que tienen asignada la macro n. */
export function macroUsers(n: number): { layer: number; name: string }[] {
  const code = macroCode(n);
  return (app.profile?.layers ?? []).flatMap((map, layer) =>
    Object.entries(map ?? {}).filter(([, a]) => a?.code === code).map(([name]) => ({ layer, name })),
  );
}

/** Borra la macro n: las teclas que la usaban vuelven a su función original y las siguientes se renumeran. */
export function deleteMacro(n: number) {
  const base = app.meta!.macroCodeBase;
  for (const map of app.profile!.layers) {
    for (const [name, a] of Object.entries(map ?? {})) {
      const i = macroIndex(a!.code);
      if (i === n) delete map![name];
      else if (i > n) a!.code = base + i - 1;
    }
  }
  macros().splice(n, 1);
}

// ---------- perfiles y teclado ----------

export async function init() {
  app.meta = await KeyboardService.Meta();
  const applied = await KeyboardService.Applied();
  app.appliedJSON = JSON.stringify(applied);
  app.profile = applied;
  app.status = await KeyboardService.Status();
  await refreshProfiles();

  Events.On("d75:status", (e) => (app.status = e.data));
  Events.On("d75:key", (e) => onPhysicalKey(e.data));
}

export async function refreshProfiles() {
  app.profiles = (await ProfileService.List()) ?? [];
}

export function loadProfile(p: Profile, name: string) {
  app.profile = p;
  app.profileName = name;
  app.savedJSON = name ? JSON.stringify(p) : "";
  app.selected = null;
}

export async function openProfile(name: string) {
  try {
    loadProfile(await ProfileService.Load(name), name);
  } catch (e) {
    toast(errorText(e), "error");
  }
}

export async function newProfile() {
  if (!isApplied() && !isSaved()) {
    const ok = await ask({
      title: t("dialog.discardTitle"),
      message: t("dialog.discardText"),
      confirm: t("dialog.discard"),
      danger: true,
    });
    if (!ok) return;
  }
  loadProfile(structuredClone($state.snapshot(app.meta!.default)), "");
}

export async function saveProfile(asNew = false) {
  let name = app.profileName;
  if (asNew || !name) {
    const answer = await ask({
      title: t(asNew ? "dialog.saveAsTitle" : "dialog.saveTitle"),
      message: t("dialog.saveText"),
      input: asNew ? (name ? name + t("dialog.copySuffix") : "") : "",
      confirm: t("action.save"),
    });
    if (typeof answer !== "string" || !answer.trim()) return;
    name = answer.trim();
  }
  try {
    const snap = $state.snapshot(app.profile!) as Profile;
    await ProfileService.Save(name, snap);
    app.profileName = name;
    app.savedJSON = JSON.stringify(snap);
    await refreshProfiles();
    toast(t("toast.saved", { name }));
  } catch (e) {
    toast(errorText(e), "error");
  }
}

/** Exporta lo que hay en el editor (guardado o no) a un archivo elegido por el usuario. */
export async function exportProfile() {
  try {
    const snap = $state.snapshot(app.profile!) as Profile;
    const ok = await ProfileService.Export(t("dialog.exportTitle"), app.profileName || "d75-perfil", snap);
    if (ok) toast(t("toast.exported"));
  } catch (e) {
    toast(errorText(e), "error");
  }
}

/** Importa un archivo como perfil nuevo y lo abre. */
export async function importProfile() {
  if (!isApplied() && !isSaved()) {
    const ok = await ask({
      title: t("dialog.discardTitle"),
      message: t("dialog.importDiscardText"),
      confirm: t("dialog.discard"),
      danger: true,
    });
    if (!ok) return;
  }
  try {
    const res = await ProfileService.Import(t("dialog.importTitle"));
    if (!res.name) return;
    await refreshProfiles();
    loadProfile(res.profile, res.name);
    toast(t("toast.imported", { name: res.name }));
  } catch (e) {
    toast(errorText(e), "error");
  }
}

export async function deleteProfile(name: string) {
  const ok = await ask({
    title: t("dialog.deleteTitle", { name }),
    message: t("dialog.deleteText"),
    confirm: t("dialog.delete"),
    danger: true,
  });
  if (!ok) return;
  try {
    await ProfileService.Delete(name);
    if (app.profileName === name) {
      app.profileName = "";
      app.savedJSON = "";
    }
    await refreshProfiles();
    toast(t("toast.deleted"));
  } catch (e) {
    toast(errorText(e), "error");
  }
}

export async function applyToKeyboard() {
  app.busy = true;
  try {
    const snap = $state.snapshot(app.profile!) as Profile;
    await KeyboardService.Apply(snap);
    app.appliedJSON = JSON.stringify(snap);
    toast(t("toast.applied"));
  } catch (e) {
    toast(errorText(e), "error");
  } finally {
    app.busy = false;
  }
}

export async function factoryReset() {
  const ok = await ask({
    title: t("dialog.resetTitle"),
    message: t("dialog.resetText"),
    confirm: t("settings.resetButton"),
    danger: true,
  });
  if (!ok) return;
  app.busy = true;
  try {
    await KeyboardService.FactoryReset();
    loadProfile(structuredClone($state.snapshot(app.meta!.default)), "");
    app.appliedJSON = JSON.stringify($state.snapshot(app.profile));
    toast(t("toast.reset"));
  } catch (e) {
    toast(errorText(e), "error");
  } finally {
    app.busy = false;
  }
}

export async function installPermissions() {
  try {
    await KeyboardService.InstallPermissions();
    app.status = await KeyboardService.Status();
    toast(t(app.status.connected ? "toast.permsOk" : "toast.permsPartial"), "info");
  } catch (e) {
    toast(errorText(e), "error");
  }
}

// ---------- detección de teclas físicas ----------

type KeyHandler = (name: string) => void;
let physicalHandler: KeyHandler | null = null;

/** Cada página registra qué hacer con una pulsación física en modo detección. */
export function onPhysical(fn: KeyHandler | null) {
  physicalHandler = fn;
}

function onPhysicalKey(ev: KeyEvent) {
  if (ev.pressed) pressed.add(ev.index);
  else pressed.delete(ev.index);
  if (!app.detect || !ev.pressed) return;
  const key = app.meta?.layout?.find((k) => k.index === ev.index);
  if (!key) {
    toast(t("toast.unknownPosition", { n: ev.index }), "error");
    return;
  }
  physicalHandler?.(key.name);
}
