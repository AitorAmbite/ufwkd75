<script lang="ts">
  import Keyboard from "../components/Keyboard.svelte";
  import KeyInspector from "../components/KeyInspector.svelte";
  import { app, layerName, onPhysical } from "../lib/state.svelte";
  import { t } from "../lib/i18n.svelte";

  $effect(() => {
    onPhysical((name) => (app.selected = name));
    return () => onPhysical(null);
  });

</script>

<div class="page">
  <div class="main">
    <div class="toolbar">
      <div class="segmented">
        {#each [0, 1, 2] as i}
          <button class:active={app.layer === i} onclick={() => (app.layer = i)}>{layerName(i)}</button>
        {/each}
      </div>
      <span class="muted hint">{t(`keys.hint.${app.layer}` as "keys.hint.0")}</span>
    </div>
    <div class="board">
      <Keyboard mode="keys" onpress={(name) => (app.selected = name)} />
    </div>
    <div class="legend faint">
      <span><i class="sw remapped"></i> {t("legend.remapped")}</span>
      <span><i class="sw selected"></i> {t("legend.selected")}</span>
      <span><i class="sw pressed"></i> {t("legend.pressed")}</span>
    </div>
  </div>
  <KeyInspector />
</div>

<style>
  .page { display: flex; gap: 16px; height: 100%; min-height: 0; }
  .main { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 16px; }
  .toolbar { display: flex; align-items: center; gap: 16px; }
  .hint { font-size: 12.5px; }
  .board { display: flex; }
  .legend { display: flex; gap: 20px; font-size: 11.5px; }
  .legend span { display: flex; align-items: center; gap: 7px; }
  .sw { width: 10px; height: 10px; display: inline-block; }
  .sw.remapped { background: #8b7dff14; border: 1px solid var(--accent); }
  .sw.selected { background: var(--accent); }
  .sw.pressed { border: 1px solid var(--success); }
</style>
