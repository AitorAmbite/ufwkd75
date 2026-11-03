<script lang="ts">
  import Icon from "../lib/Icon.svelte";
  import { paint, ZONES, fillZone, fillGradient, fillRainbow, clearAll, type Tool } from "../lib/light.svelte";
  import { fromHex, toHex, isHex } from "../lib/color";
  import { t } from "../lib/i18n.svelte";

  let hexDraft = $state(paint.color);
  $effect(() => { hexDraft = paint.color; });

  const rgb = $derived(fromHex(paint.color));

  function setChannel(i: number, v: number) {
    const c = [...rgb];
    c[i] = v;
    paint.color = toHex(c);
  }

  function onHex(v: string) {
    hexDraft = v;
    const h = (v.startsWith("#") ? v : "#" + v).toLowerCase();
    if (isHex(h)) paint.color = h;
  }

  const tools: { id: Tool; icon: string }[] = [
    { id: "brush", icon: "brush" },
    { id: "picker", icon: "pipette" },
    { id: "eraser", icon: "eraser" },
  ];
  const channels = ["R", "G", "B"];
</script>

<div class="panel">
  <section>
    <div class="label">{t("color.title")}</div>
    <div class="picker">
      <label class="swatch-big" style:background={paint.color}>
        <input type="color" bind:value={paint.color} />
      </label>
      <div class="picker-side">
        <input class="input mono" value={hexDraft} oninput={(e) => onHex(e.currentTarget.value)} maxlength="7" spellcheck="false" />
        <div class="faint small">rgb({rgb.join(", ")})</div>
      </div>
    </div>
    {#each channels as ch, i}
      <div class="channel">
        <span class="faint">{ch}</span>
        <input type="range" min="0" max="255" value={rgb[i]}
          style:--fill="{(rgb[i] / 255) * 100}%"
          oninput={(e) => setChannel(i, +e.currentTarget.value)} />
        <span class="mono faint">{rgb[i]}</span>
      </div>
    {/each}
    {#if paint.recent.length}
      <div class="recent">
        {#each paint.recent as c}
          <button class="dot" style:background={c} title={c} aria-label={c} onclick={() => (paint.color = c)}></button>
        {/each}
      </div>
    {/if}
  </section>

  <section>
    <div class="label">{t("color.tool")}</div>
    <div class="tools">
      {#each tools as tool}
        <button class="tool" class:active={paint.tool === tool.id} title={t(`tool.${tool.id}Hint`)} onclick={() => (paint.tool = tool.id)}>
          <Icon name={tool.icon} size={17} />
          <span>{t(`tool.${tool.id}`)}</span>
        </button>
      {/each}
    </div>
  </section>

  <section>
    <div class="label">{t("color.fill")}</div>
    <div class="zones">
      {#each ZONES as z}
        <button class="chip" onclick={() => fillZone(z.test)}>{t(`zone.${z.id}`)}</button>
      {/each}
      <button class="chip danger" onclick={clearAll}>{t("color.clearAll")}</button>
    </div>
  </section>

  <section>
    <div class="label">{t("color.gradient")}</div>
    <div class="gradient">
      <label class="swatch-sm" style:background={paint.gradA}><input type="color" bind:value={paint.gradA} /></label>
      <div class="grad-preview" style:background="linear-gradient(90deg, {paint.gradA}, {paint.gradB})"></div>
      <label class="swatch-sm" style:background={paint.gradB}><input type="color" bind:value={paint.gradB} /></label>
    </div>
    <div class="row">
      <button class="btn sm" onclick={fillGradient}>{t("color.applyGradient")}</button>
      <button class="btn sm" onclick={fillRainbow}><span class="rainbow-bar"></span> {t("color.rainbow")}</button>
    </div>
  </section>
</div>

<style>
  .panel { display: flex; flex-direction: column; gap: 20px; }
  section { display: flex; flex-direction: column; gap: 10px; }
  .picker { display: flex; gap: 10px; align-items: stretch; }
  .swatch-big, .swatch-sm { position: relative; cursor: pointer; flex: none; border: 1px solid #ffffff33; }
  .swatch-big { width: 60px; height: 60px; }
  .swatch-sm { width: 30px; height: 30px; }
  .swatch-big input, .swatch-sm input { position: absolute; inset: 0; opacity: 0; cursor: pointer; width: 100%; height: 100%; }
  .picker-side { display: flex; flex-direction: column; justify-content: space-between; gap: 6px; flex: 1; min-width: 0; }
  .picker-side .input { width: 100%; }
  .small { font-size: 11px; font-family: var(--mono); }
  .channel { display: grid; grid-template-columns: 14px 1fr 30px; gap: 10px; align-items: center; font-size: 12px; font-family: var(--mono); }
  .channel .mono { text-align: right; }
  .recent { display: flex; flex-wrap: wrap; gap: 4px; }
  .dot { width: 20px; height: 20px; border: 1px solid #ffffff33; cursor: pointer; }
  .dot:hover { border-color: #fff; }
  .tools { display: grid; grid-template-columns: repeat(3, 1fr); border: 1px solid var(--border-strong); }
  .tool {
    display: flex; flex-direction: column; align-items: center; gap: 5px; padding: 10px 4px;
    border: 0; border-right: 1px solid var(--border-strong); background: transparent; cursor: pointer; font-size: 12px; color: var(--text-2);
  }
  .tool:last-child { border-right: 0; }
  .tool:hover { color: var(--text); background: var(--surface-2); }
  .tool.active { background: var(--text); color: var(--bg); }
  .zones { display: flex; flex-wrap: wrap; gap: 4px; }
  .chip.danger { color: var(--danger); }
  .gradient { display: flex; align-items: center; gap: 8px; }
  .grad-preview { flex: 1; height: 30px; border: 1px solid #ffffff33; }
  .row { display: flex; gap: 6px; }
  .row .btn { flex: 1; }
  .rainbow-bar {
    width: 14px; height: 14px; flex: none;
    background: linear-gradient(90deg, #ff4d6d 0 20%, #ffb347 20% 40%, #3ddc97 40% 60%, #29c5ff 60% 80%, #8b7dff 80%);
  }
</style>
