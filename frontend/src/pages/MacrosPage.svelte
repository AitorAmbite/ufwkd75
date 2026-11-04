<script lang="ts">
  import Icon from "../lib/Icon.svelte";
  import {
    app, macros, addMacro, deleteMacro, macroUsers, keyByName, layerName, codeLabel, ask, toast,
  } from "../lib/state.svelte";
  import { t, tKey, tCode, tGroup } from "../lib/i18n.svelte";
  import { hidFromEvent, MACRO_GROUPS } from "../lib/hid";

  const MODES = ["count", "hold", "toggle"] as const;

  let current = $state(0);
  let recording = $state(false);
  let fixedDelay = $state(30);
  let lastTime = 0;

  const list = $derived(macros());
  const macro = $derived(list[current]);
  const max = $derived(app.meta!.maxMacroEvents);
  const minDelay = $derived(app.meta!.minMacroDelay);
  const maxDelay = $derived(app.meta!.maxMacroDelay);
  const events = $derived(macro?.events ?? []);
  const users = $derived(macro ? macroUsers(current) : []);
  const duration = $derived(events.reduce((s, e) => s + Math.max(e.delay, minDelay), 0));
  const groups = $derived(
    MACRO_GROUPS.map((g) => ({ g, codes: (app.meta?.codes ?? []).filter((c) => c.group === g) })),
  );

  const clampDelay = (ms: number) => Math.min(maxDelay, Math.max(minDelay, Math.round(ms) || minDelay));

  // Si se borra la macro seleccionada, o se cambia de perfil, la selección sigue siendo válida
  $effect(() => {
    if (current >= list.length) current = Math.max(0, list.length - 1);
  });
  // Al cambiar de macro o salir de la página se deja de grabar
  $effect(() => {
    current;
    return () => (recording = false);
  });

  function create() {
    current = addMacro();
  }

  async function remove() {
    const ok = await ask({
      title: t("macros.deleteTitle", { name: macro.name }),
      message: t("macros.deleteText"),
      confirm: t("dialog.delete"),
      danger: true,
    });
    if (ok) deleteMacro(current);
  }

  function ensureEvents() {
    if (!macro.events) macro.events = [];
    return macro.events;
  }

  function toggleRecording() {
    recording = !recording;
    lastTime = 0;
    if (recording) (document.activeElement as HTMLElement | null)?.blur();
  }

  // Graba desde el teclado del sistema: KeyboardEvent.code es la posición física, como el usage HID
  function onKey(e: KeyboardEvent, down: boolean) {
    if (!recording) return;
    e.preventDefault();
    e.stopPropagation();
    if (e.repeat) return;
    const code = hidFromEvent(e.code, app.meta!.layout!);
    if (code === undefined) {
      toast(t("macros.unknownKey", { key: e.code }), "error");
      return;
    }
    const evs = ensureEvents();
    if (evs.length >= max) {
      recording = false;
      toast(t("macros.full", { max }), "error");
      return;
    }
    if (evs.length && lastTime) evs[evs.length - 1].delay = clampDelay(e.timeStamp - lastTime);
    evs.push({ code, down, delay: minDelay });
    lastTime = e.timeStamp;
  }

  function addKey(e: Event) {
    const sel = e.currentTarget as HTMLSelectElement;
    const code = Number(sel.value);
    sel.value = "";
    if (!code) return;
    const evs = ensureEvents();
    if (evs.length + 2 > max) return toast(t("macros.full", { max }), "error");
    evs.push({ code, down: true, delay: fixedDelay }, { code, down: false, delay: fixedDelay });
  }

  function applyFixedDelay() {
    fixedDelay = clampDelay(fixedDelay);
    for (const ev of events) ev.delay = fixedDelay;
  }

  function move(i: number, dir: -1 | 1) {
    const evs = ensureEvents();
    const j = i + dir;
    if (j < 0 || j >= evs.length) return;
    [evs[i], evs[j]] = [evs[j], evs[i]];
  }
</script>

<svelte:window onkeydowncapture={(e) => onKey(e, true)} onkeyupcapture={(e) => onKey(e, false)} />

<div class="page">
  <div class="card list">
    <div class="list-head">
      <span class="label">{t("macros.list", { n: list.length, max: app.meta!.maxMacros })}</span>
      <button class="btn ghost sm icon" title={t("macros.new")} disabled={list.length >= app.meta!.maxMacros} onclick={create}>
        <Icon name="plus" size={16} />
      </button>
    </div>
    {#each list as m, i}
      <button class="item" class:active={i === current} onclick={() => (current = i)}>
        <span class="num">M{i + 1}</span>
        <span class="iname">{m.name}</span>
        {#if !m.events?.length}<span class="warn" title={t("macros.emptyWarn")}><Icon name="alert" size={14} /></span>{/if}
      </button>
    {:else}
      <p class="faint small pad">{t("macros.emptyList")}</p>
    {/each}
  </div>

  {#if macro}
    <div class="editor">
      <div class="card section row">
        <div class="field grow">
          <span class="label">{t("macros.name")}</span>
          <input class="input" maxlength="32" bind:value={macro.name} />
        </div>
        <div class="field">
          <span class="label">{t("macros.repeat")}</span>
          <div class="row">
            <div class="segmented">
              {#each MODES as mode}
                <button class:active={(macro.mode || "count") === mode} onclick={() => (macro.mode = mode)}>{t(`macros.mode.${mode}`)}</button>
              {/each}
            </div>
            {#if (macro.mode || "count") === "count"}
              <input class="input num-input" type="number" min="1" max="65534" value={macro.count ?? 1}
                onchange={(e) => (macro.count = Math.min(65534, Math.max(1, Number(e.currentTarget.value) || 1)))} />
            {/if}
          </div>
        </div>
        <button class="btn ghost icon del" title={t("macros.delete")} onclick={remove}><Icon name="trash" size={16} /></button>
      </div>
      <p class="muted hint">{t(`macros.modeHint.${macro.mode || "count"}` as "macros.modeHint.count")}</p>

      <div class="card section events-card">
        <div class="row toolbar">
          <button class="btn" class:rec={recording} onclick={toggleRecording}>
            <Icon name={recording ? "stop" : "record"} size={15} />
            {recording ? t("macros.stop") : t("macros.record")}
          </button>
          <select class="input" onchange={addKey} disabled={recording}>
            <option value="">{t("macros.addKey")}</option>
            {#each groups as { g, codes }}
              <optgroup label={tGroup(g)}>
                {#each codes as c}<option value={c.code}>{tCode(c.code, c.label)}</option>{/each}
              </optgroup>
            {/each}
          </select>
          <div class="spacer"></div>
          <label class="delay-set faint">
            {t("macros.fixedDelay")}
            <input class="input num-input" type="number" min={minDelay} max={maxDelay} bind:value={fixedDelay} />
            ms
          </label>
          <button class="btn sm" disabled={!events.length} onclick={applyFixedDelay}>{t("macros.applyDelay")}</button>
          <button class="btn sm ghost" disabled={!events.length} onclick={() => (macro.events = [])}>{t("macros.clear")}</button>
        </div>

        {#if recording}
          <div class="recording"><span class="rec-dot"></span>{t("macros.recording")}</div>
        {/if}

        <div class="label">
          {t("macros.events", { n: events.length, max })}
          {#if events.length}<span class="faint"> · ≈ {(duration / 1000).toFixed(2)} s</span>{/if}
        </div>
        <div class="events">
          {#each events as ev, i}
            <div class="event">
              <span class="idx faint">{i + 1}</span>
              <button class="dir" class:down={ev.down} title={t("macros.toggleDir")} onclick={() => (ev.down = !ev.down)}>
                <Icon name={ev.down ? "arrowDown" : "arrowUp"} size={14} />
                {ev.down ? t("macros.press") : t("macros.release")}
              </button>
              <select class="input key-select" bind:value={ev.code}>
                {#if !groups.some(({ codes }) => codes.some((c) => c.code === ev.code))}
                  <option value={ev.code}>{codeLabel(ev.code)}</option>
                {/if}
                {#each groups as { g, codes }}
                  <optgroup label={tGroup(g)}>
                    {#each codes as c}<option value={c.code}>{tCode(c.code, c.label)}</option>{/each}
                  </optgroup>
                {/each}
              </select>
              <label class="delay faint" title={t("macros.delayHint")}>
                <Icon name="clock" size={13} />
                <input class="input num-input" type="number" min={minDelay} max={maxDelay} value={ev.delay}
                  onchange={(e) => (ev.delay = clampDelay(Number(e.currentTarget.value)))} />
                ms
              </label>
              <div class="row-actions">
                <button class="btn ghost sm icon" disabled={i === 0} onclick={() => move(i, -1)}><Icon name="chevronUp" size={14} /></button>
                <button class="btn ghost sm icon" disabled={i === events.length - 1} onclick={() => move(i, 1)}><Icon name="chevronDown" size={14} /></button>
                <button class="btn ghost sm icon" title={t("macros.remove")} onclick={() => macro.events!.splice(i, 1)}><Icon name="x" size={14} /></button>
              </div>
            </div>
          {:else}
            <p class="faint small pad">{t("macros.noEvents")}</p>
          {/each}
        </div>
      </div>

      <div class="card section">
        <div class="label">{t("macros.usedBy")}</div>
        {#if users.length}
          <div class="users">
            {#each users as u}
              <button class="chip" onclick={() => { app.layer = u.layer; app.selected = u.name; app.page = "keys"; }}>
                {layerName(u.layer)} · {tKey(u.name, keyByName(u.name).label)}
              </button>
            {/each}
          </div>
        {:else}
          <p class="faint small">{t("macros.unused")}</p>
        {/if}
      </div>
    </div>
  {:else}
    <div class="card empty">
      <div class="empty-icon"><Icon name="macro" size={26} /></div>
      <h3>{t("macros.emptyTitle")}</h3>
      <p class="muted">{t("macros.emptyText")}</p>
      <button class="btn primary" onclick={create}><Icon name="plus" size={16} /> {t("macros.new")}</button>
    </div>
  {/if}
</div>

<style>
  .page { display: flex; gap: 16px; height: 100%; min-height: 0; }
  .list { width: 220px; flex: none; display: flex; flex-direction: column; overflow-y: auto; }
  .list-head { display: flex; align-items: center; justify-content: space-between; padding: 10px 8px 6px 14px; }
  .item {
    display: flex; align-items: center; gap: 10px; height: 34px; padding: 0 12px; border: 0;
    border-left: 2px solid transparent; background: transparent; color: var(--text-2); cursor: pointer; text-align: left;
  }
  .item:hover { background: var(--surface-2); color: var(--text); }
  .item.active { border-left-color: var(--accent); background: var(--surface-2); color: var(--text); }
  .num { font-family: var(--mono); font-size: 11px; color: var(--text-3); width: 24px; }
  .iname { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .warn { color: var(--warning); display: grid; }
  .small { font-size: 12px; }
  .pad { padding: 4px 14px; margin: 0; }

  .editor { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 12px; min-height: 0; }
  .section { padding: 14px 16px; display: flex; flex-direction: column; gap: 10px; }
  .row { display: flex; align-items: center; gap: 10px; flex-direction: row; }
  .section.row { align-items: flex-end; flex-wrap: wrap; }
  .field { display: flex; flex-direction: column; gap: 6px; }
  .grow { flex: 1; min-width: 160px; }
  .del:hover { color: var(--danger); }
  .hint { margin: -4px 2px 0; font-size: 12.5px; }
  .num-input { width: 76px; }

  .events-card { flex: 1; min-height: 220px; }
  .toolbar { flex-wrap: wrap; }
  .toolbar select { width: 180px; }
  .spacer { flex: 1; }
  .delay-set { display: flex; align-items: center; gap: 6px; font-size: 12.5px; }
  .rec { background: var(--danger) !important; border-color: var(--danger) !important; color: #0d0d0e; }
  .recording { display: flex; align-items: center; gap: 8px; padding: 8px 10px; border: 1px solid var(--danger); color: var(--danger); font-size: 12.5px; }
  .rec-dot { width: 8px; height: 8px; background: var(--danger); animation: blink 1s steps(2) infinite; }
  @keyframes blink { 50% { opacity: 0; } }

  .events { flex: 1; overflow-y: auto; border-top: 1px solid var(--border); min-height: 0; }
  .event { display: flex; align-items: center; gap: 10px; padding: 5px 4px; border-bottom: 1px solid var(--border); }
  .event:hover { background: var(--surface-2); }
  .idx { width: 22px; text-align: right; font-family: var(--mono); font-size: 11px; }
  .dir {
    width: 92px; display: flex; align-items: center; gap: 6px; height: 28px; padding: 0 8px; border: 1px solid var(--border-strong);
    background: transparent; color: var(--text-2); cursor: pointer; font-size: 12px;
  }
  .dir.down { color: var(--accent-2); border-color: var(--accent); }
  .key-select { width: 160px; height: 28px; }
  .delay { display: flex; align-items: center; gap: 6px; font-size: 12px; }
  .delay .num-input { height: 28px; }
  .row-actions { margin-left: auto; display: flex; gap: 2px; }

  .users { display: flex; flex-wrap: wrap; gap: 6px; }

  .empty { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px; text-align: center; padding: 32px; }
  .empty-icon { width: 52px; height: 52px; display: grid; place-items: center; border: 1px solid var(--border-strong); color: var(--accent-2); }
  .empty h3 { margin: 6px 0 0; font-size: 15px; }
  .empty p { margin: 0 0 8px; max-width: 380px; font-size: 13px; }
</style>
