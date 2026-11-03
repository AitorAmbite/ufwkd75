<script lang="ts">
  import { app, pressed, layerMap, assignLabel } from "../lib/state.svelte";
  import { toHex, luminance, isLit, type RGB } from "../lib/color";
  import type { Key } from "../../bindings/d75studio/internal/d75/models";
  import { tKey } from "../lib/i18n.svelte";

  let {
    mode = "keys",
    colorFor = () => null,
    onpress = () => {},
    onenter = () => {},
  }: {
    mode?: "keys" | "light";
    colorFor?: (k: Key) => RGB | null;
    onpress?: (name: string) => void;
    onenter?: (name: string) => void;
  } = $props();

  const UNITS_W = 16.5;
  const UNITS_H = 6.5;
  const PAD = 18;

  let width = $state(900);
  const u = $derived(Math.max(30, Math.min(64, (width - PAD * 2) / UNITS_W)));
  const layout = $derived(app.meta?.layout ?? []);
  const assigned = $derived(mode === "keys" && app.profile ? layerMap() : {});

  // El relleno es el 72% del color (ver .key.lit): el texto contrasta con ese tono
  const textOn = (rgb: RGB) => (luminance(rgb) * 0.72 > 95 ? "#0d0d0e" : "#ffffff");

  // Las etiquetas largas en teclas pequeñas encogen para no cortarse
  function fontSize(k: Key, label: string) {
    const base = Math.max(9.5, u * 0.21);
    const room = (k.w * u - 12) / (label.length * base * 0.64);
    return Math.max(8.5, base * Math.min(1, room));
  }

</script>

<div class="wrap" bind:clientWidth={width}>
  <div
    class="case"
    class:light={mode === "light"}
    style:width="{u * UNITS_W + PAD * 2}px"
    style:height="{u * UNITS_H + PAD * 2}px"
  >
    {#each layout as k (k.name)}
      {@const a = assigned[k.name]}
      {@const rgb = mode === "light" ? colorFor(k) : null}
      {@const lit = isLit(rgb)}
      {@const label = tKey(k.name, k.label)}
      <button
        class="key"
        class:selected={app.selected === k.name && mode === "keys"}
        class:remapped={!!a}
        class:pressed={pressed.has(k.index)}
        class:lit
        style:left="{PAD + k.x * u}px"
        style:top="{PAD + k.y * u}px"
        style:width="{k.w * u - 5}px"
        style:height="{u - 5}px"
        style:--c={lit ? toHex(rgb!) : null}
        style:color={lit ? textOn(rgb!) : null}
        style:font-size="{fontSize(k, a ? assignLabel(a) : label)}px"
        onmousedown={(e) => { e.preventDefault(); onpress(k.name); }}
        onmouseenter={() => onenter(k.name)}
        title={a ? `${label} → ${assignLabel(a)}` : label}
      >
        {#if a}
          <span class="main">{assignLabel(a)}</span>
          <span class="sub">{label}</span>
        {:else}
          <span class="main">{label}</span>
        {/if}
      </button>
    {/each}
  </div>
</div>

<style>
  .wrap { width: 100%; display: flex; justify-content: center; }

  .case { position: relative; flex: none; background: #111113; border: 1px solid var(--border-strong); }

  .key {
    position: absolute; display: flex; flex-direction: column; align-items: center; justify-content: center;
    gap: 2px; padding: 3px; border-radius: 2px; cursor: pointer; overflow: hidden;
    color: #d4d4d8; line-height: 1.1; text-align: center;
    background: #1c1c1f; border: 1px solid #303036;
    transition: background .08s, border-color .08s;
  }
  .key:hover { background: #26262a; border-color: #4a4a52; }
  .key .main { font-weight: 600; max-width: 100%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .key .sub { font-size: .78em; opacity: .55; font-family: var(--mono); }

  .key.remapped { color: var(--accent-2); border-color: var(--accent); background: #8b7dff14; }
  .key.remapped .main { font-size: .9em; white-space: normal; }

  .key.selected { background: var(--accent); border-color: var(--accent); color: #0d0d0e; }
  .key.selected .sub { opacity: .7; }

  .key.pressed { border-color: var(--success); box-shadow: inset 0 0 0 1px var(--success); }

  /* Modo iluminación: relleno plano con el color de la luz */
  .case.light .key { color: #5d5d66; background: #161618; border-color: #26262b; }
  .key.lit {
    background: color-mix(in srgb, var(--c) 72%, #000) !important;
    border-color: var(--c) !important;
  }
  .case.light .key:hover { border-color: #fff; }
</style>
