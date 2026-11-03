<script lang="ts">
  import Icon from "../lib/Icon.svelte";
  import { app, MODS, layerName, keyByName, layerMap, assign, resetKey, assignLabel, codeLabel } from "../lib/state.svelte";
  import { t, tKey, tCode, tGroup } from "../lib/i18n.svelte";

  let query = $state("");
  let group = $state("all");

  const key = $derived(app.selected ? keyByName(app.selected) : null);
  const current = $derived(key ? layerMap()[key.name] : undefined);
  const mods = $derived(current?.mods ?? 0);
  const groups = $derived(["all", ...new Set((app.meta?.codes ?? []).map((c) => c.group))]);
  const visible = $derived(
    (app.meta?.codes ?? []).filter((c) => {
      const q = query.trim().toLowerCase();
      if (q) return tCode(c.code, c.label).toLowerCase().includes(q) || tGroup(c.group).toLowerCase().includes(q);
      return group === "all" || c.group === group;
    }),
  );
  const remaps = $derived(Object.entries(app.profile ? layerMap() : {}));

  function toggleMod(bit: number) {
    if (!key) return;
    const code = current?.code ?? key.default;
    assign(key.name, code, mods ^ bit);
  }
</script>

<div class="inspector card">
  {#if key}
    <div class="head">
      <div class="preview">
        <div class="cap">{tKey(key.name, key.label)}</div>
        <Icon name="arrow" size={18} class="faint" />
        <div class="cap target" class:changed={!!current}>{current ? assignLabel(current) : codeLabel(key.default)}</div>
      </div>
      <button class="btn ghost sm icon" title={t("inspector.close")} onclick={() => (app.selected = null)}><Icon name="x" size={16} /></button>
    </div>
    <div class="faint small">{t("inspector.layerState", { layer: layerName(app.layer), state: t(current ? "inspector.remapped" : "inspector.original") })}</div>

    <div class="block">
      <div class="label">{t("inspector.modifiers")}</div>
      <div class="mods">
        {#each MODS as [bit, name]}
          <button class="chip" class:active={mods & bit} onclick={() => toggleMod(bit)}>{name}</button>
        {/each}
      </div>
    </div>

    <div class="block grow">
      <div class="search">
        <Icon name="search" size={15} class="faint" />
        <input class="input" placeholder={t("inspector.search")} bind:value={query} />
      </div>
      {#if !query}
        <div class="groups">
          {#each groups as g}
            <button class="chip" class:active={group === g} onclick={() => (group = g)}>{tGroup(g)}</button>
          {/each}
        </div>
      {/if}
      <div class="codes">
        {#each visible as c (c.group + c.code)}
          <button
            class="code"
            class:active={current ? current.code === c.code : key.default === c.code}
            onclick={() => assign(key.name, c.code, mods)}
            title={tGroup(c.group)}
          >{tCode(c.code, c.label)}</button>
        {/each}
        {#if visible.length === 0}<p class="faint small">{t("inspector.noMatch", { q: query })}</p>{/if}
      </div>
    </div>

    <button class="btn full" disabled={!current} onclick={() => resetKey(key.name)}>
      <Icon name="rotate" size={15} /> {t("inspector.restore")}
    </button>
  {:else}
    <div class="empty">
      <div class="empty-icon"><Icon name="keyboard" size={26} /></div>
      <h3>{t("inspector.emptyTitle")}</h3>
      <p class="muted">{@html t("inspector.emptyText")}</p>
    </div>
    <div class="block grow">
      <div class="label">{t("inspector.changes", { layer: layerName(app.layer), n: remaps.length })}</div>
      <div class="remaps">
        {#each remaps as [name, a] (name)}
          <div class="remap">
            <button class="remap-btn" onclick={() => (app.selected = name)}>
              <span class="mini-cap">{tKey(name, keyByName(name).label)}</span>
              <Icon name="arrow" size={14} class="faint" />
              <span class="ellipsis">{assignLabel(a!)}</span>
            </button>
            <button class="btn ghost sm icon" title={t("inspector.restoreOne")} onclick={() => resetKey(name)}><Icon name="rotate" size={14} /></button>
          </div>
        {:else}
          <p class="faint small">{t("inspector.noChanges")}</p>
        {/each}
      </div>
    </div>
  {/if}
</div>

<style>
  .inspector { width: 340px; flex: none; display: flex; flex-direction: column; gap: 16px; padding: 16px; min-height: 0; }
  .head { display: flex; align-items: flex-start; justify-content: space-between; }
  .preview { display: flex; align-items: center; gap: 10px; min-width: 0; }
  .cap {
    min-width: 52px; height: 52px; padding: 0 12px; display: grid; place-items: center;
    font-weight: 700; font-size: 15px; white-space: nowrap; max-width: 150px; overflow: hidden; text-overflow: ellipsis;
    background: #1c1c1f; border: 1px solid #3a3a41; border-radius: 2px;
  }
  .cap.target.changed { color: var(--accent-2); border-color: var(--accent); background: #8b7dff14; font-size: 13px; }
  .small { font-size: 11.5px; margin-top: -10px; font-family: var(--mono); }
  .block { display: flex; flex-direction: column; gap: 8px; }
  .grow { flex: 1; min-height: 0; }
  .mods { display: flex; flex-wrap: wrap; gap: 4px; }
  .groups { display: flex; gap: 4px; overflow-x: auto; padding-bottom: 6px; margin-bottom: -4px; }
  .groups .chip { flex: none; }
  .search { position: relative; display: flex; align-items: center; }
  .search :global(svg) { position: absolute; left: 10px; }
  .search .input { width: 100%; padding-left: 32px; }
  .codes {
    flex: 1; overflow-y: auto; display: grid; grid-template-columns: repeat(3, 1fr); gap: 4px;
    align-content: start; min-height: 120px;
  }
  .code {
    height: 34px; padding: 0 6px; border-radius: 2px; border: 1px solid var(--border-strong); background: var(--surface-2);
    cursor: pointer; font-size: 12.5px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
  }
  .code:hover { border-color: #5a5a62; background: var(--surface-3); }
  .code.active { background: var(--accent); border-color: var(--accent); color: #0d0d0e; }
  .full { width: 100%; }

  .empty { padding: 8px 0 4px; }
  .empty-icon { width: 44px; height: 44px; margin-bottom: 12px; display: grid; place-items: center; border: 1px solid var(--border-strong); color: var(--accent-2); }
  .empty h3 { margin: 0 0 6px; font-size: 14px; }
  .empty p { margin: 0; font-size: 12.5px; }
  .remaps { overflow-y: auto; display: flex; flex-direction: column; border-top: 1px solid var(--border); }
  .remap { display: flex; align-items: center; gap: 4px; border-bottom: 1px solid var(--border); padding-right: 2px; }
  .remap:hover { background: var(--surface-2); }
  .remap-btn {
    flex: 1; min-width: 0; display: flex; align-items: center; gap: 10px; padding: 8px 4px; border: 0;
    background: transparent; cursor: pointer; text-align: left;
  }
  .mini-cap { padding: 1px 7px; border: 1px solid var(--border-strong); font-size: 12px; font-weight: 600; flex: none; }
  .ellipsis { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
