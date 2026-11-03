<script lang="ts">
  import Icon from "../lib/Icon.svelte";
  import Keyboard from "../components/Keyboard.svelte";
  import ColorPanel from "../components/ColorPanel.svelte";
  import { app, onPhysical } from "../lib/state.svelte";
  import { paint, useTool, keyColor, keyPosition, setKeyColor } from "../lib/light.svelte";
  import { fromHex, hsvToRgb, dim, type RGB } from "../lib/color";
  import type { Key } from "../../bindings/d75studio/internal/d75/models";
  import { t, tEffect } from "../lib/i18n.svelte";

  const CONSTANT = 16;
  const BREATHING = [3, 8];

  const led = $derived(app.profile!.led);
  const custom = $derived(led.effect === app.meta!.customEffect);
  const palette = $derived(app.meta!.palette ?? []);

  // Reloj de la vista previa animada
  let clock = $state(0);
  $effect(() => {
    const id = setInterval(() => (clock += 0.04), 40);
    return () => clearInterval(id);
  });

  $effect(() => {
    onPhysical((name) => {
      if (!custom) return;
      useTool(name, true);
      paint.painting = false;
    });
    return () => onPhysical(null);
  });

  function colorFor(k: Key): RGB | null {
    if (custom) {
      const c = keyColor(k.name);
      return c ? dim(c, led.brightness) : null;
    }
    const pos = keyPosition(k);
    const speed = 0.4 + led.speed * 0.3;
    const animated = led.effect !== CONSTANT;
    const base = led.color === 7
      ? hsvToRgb((pos * 300 + (animated ? clock * 60 * speed : 0)) % 360, 1, 1)
      : fromHex(palette[led.color]);
    let f = 1;
    if (BREATHING.includes(led.effect)) f = 0.12 + 0.88 * (0.5 + 0.5 * Math.sin(clock * speed * 2));
    else if (animated && led.color !== 7) f = 0.3 + 0.7 * (0.5 + 0.5 * Math.sin(clock * speed * 2 - pos * 9));
    return dim(base.map((v) => v * f), led.brightness);
  }

  function setEffect(i: number) {
    led.effect = i;
  }

  // "Color libre": pasa al modo personalizado con todo el teclado de un color
  function freeColor() {
    if (led.color < 7) paint.color = palette[led.color];
    const wasEmpty = !led.custom || Object.keys(led.custom).length === 0;
    led.effect = app.meta!.customEffect;
    if (wasEmpty) for (const k of app.meta!.layout!) setKeyColor(k.name, fromHex(paint.color));
  }

  // Miniatura de cada efecto: bandas de color planas
  function effectStyle(i: number): string {
    if (i === app.meta!.customEffect)
      return "background: conic-gradient(#ff4d6d 0 25%, #ffb347 0 50%, #29c5ff 0 75%, #8b7dff 0)";
    const hue = (i * 47) % 360;
    if (i === CONSTANT) return `background: hsl(${hue} 75% 55%)`;
    const c = [0, 70, 160].map((d) => `hsl(${(hue + d) % 360} 80% 58%)`);
    return `background: linear-gradient(90deg, ${c[0]} 0 33%, ${c[1]} 33% 66%, ${c[2]} 66%)`;
  }
</script>

<div class="page">
  <div class="main">
    <div class="preview">
      <Keyboard
        mode="light"
        {colorFor}
        onpress={(name) => custom && useTool(name, true)}
        onenter={(name) => custom && useTool(name, false)}
      />
      <div class="faint note">
        {t(custom ? "light.noteCustom" : "light.notePreview")}
      </div>
    </div>

    <div class="effects">
      <div class="label">{t("light.effects")}</div>
      <div class="grid">
        {#each app.meta!.effects ?? [] as name, i}
          <button class="effect" class:active={led.effect === i} onclick={() => setEffect(i)}>
            <span class="thumb" class:breathe={BREATHING.includes(i)} style={effectStyle(i)}></span>
            <span class="ename">{tEffect(i, name)}</span>
          </button>
        {/each}
      </div>
    </div>
  </div>

  <aside class="side card">
    <section>
      <div class="label">{t("light.brightness")}</div>
      <div class="slider-row">
        <Icon name="sun" size={15} class="faint" />
        <input type="range" min="0" max="4" bind:value={led.brightness} style:--fill="{(led.brightness / 4) * 100}%" />
        <span class="value">{led.brightness}</span>
      </div>
    </section>

    {#if custom}
      <ColorPanel />
    {:else}
      <section>
        <div class="label">{t("light.speed")}</div>
        <div class="slider-row">
          <Icon name="zap" size={15} class="faint" />
          <input type="range" min="1" max="5" bind:value={led.speed} style:--fill="{((led.speed - 1) / 4) * 100}%" />
          <span class="value">{led.speed}</span>
        </div>
      </section>

      <section>
        <div class="label">{t("light.effectColor")}</div>
        <div class="swatches">
          {#each palette as c, i}
            <button class="sw" class:active={led.color === i} style:background={c} aria-label={c} onclick={() => (led.color = i)}></button>
          {/each}
          <button class="sw multi" class:active={led.color === 7} title={t("light.multicolor")} aria-label={t("light.multicolor")} onclick={() => (led.color = 7)}></button>
        </div>
        <p class="faint info">{t("light.colorLimit")}</p>
      </section>

      <div class="free card">
        <Icon name="palette" size={22} />
        <div>
          <div class="free-title">{t("light.freeTitle")}</div>
          <div class="faint info">{t("light.freeText")}</div>
        </div>
        <button class="btn primary sm" onclick={freeColor}>{t("light.freeButton")}</button>
      </div>
    {/if}
  </aside>
</div>

<style>
  .page { display: flex; gap: 16px; height: 100%; min-height: 0; }
  .main { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 20px; overflow-y: auto; padding-right: 4px; }
  .preview { display: flex; flex-direction: column; gap: 8px; }
  .note { font-size: 11.5px; }

  .effects { display: flex; flex-direction: column; gap: 10px; }
  .grid {
    display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    border-top: 1px solid var(--border); border-left: 1px solid var(--border);
  }
  .effect {
    display: flex; align-items: center; gap: 10px; padding: 9px 10px; text-align: left;
    border: 0; border-right: 1px solid var(--border); border-bottom: 1px solid var(--border);
    background: var(--surface); cursor: pointer;
  }
  .effect:hover { background: var(--surface-2); }
  .effect.active { background: var(--text); color: var(--bg); }
  .thumb { width: 24px; height: 24px; flex: none; border: 1px solid #ffffff26; }
  .thumb.breathe { animation: breathe 2.4s steps(6) infinite; }
  .ename { font-size: 12.5px; line-height: 1.2; }
  @keyframes breathe { 0%, 100% { opacity: .35; } 50% { opacity: 1; } }

  .side { width: 320px; flex: none; padding: 16px; display: flex; flex-direction: column; gap: 22px; overflow-y: auto; }
  section { display: flex; flex-direction: column; gap: 10px; }
  .slider-row { display: grid; grid-template-columns: 16px 1fr 18px; gap: 10px; align-items: center; }
  .value { text-align: right; font-family: var(--mono); }
  .swatches { display: grid; grid-template-columns: repeat(8, 1fr); border: 1px solid var(--border-strong); }
  .sw { height: 30px; border: 0; border-right: 1px solid #0d0d0e; cursor: pointer; position: relative; }
  .sw:last-child { border-right: 0; }
  .sw.active { outline: 2px solid var(--text); outline-offset: -2px; z-index: 1; }
  .sw.multi { background: linear-gradient(180deg, #f33 0 25%, #fd3 25% 50%, #3c6 50% 75%, #36f 75%); }
  .info { font-size: 12px; margin: 0; line-height: 1.4; }
  .free { display: flex; flex-direction: column; gap: 10px; padding: 14px; border-color: var(--border-strong); background: var(--bg); }
  .free :global(svg) { color: var(--accent-2); }
  .free-title { color: var(--text); font-weight: 600; margin-bottom: 2px; }
  .free .btn { align-self: flex-start; }
</style>
