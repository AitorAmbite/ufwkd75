<script lang="ts">
  import Icon from "../lib/Icon.svelte";
  import { app, factoryReset, installPermissions } from "../lib/state.svelte";
  import { t, i18n, LOCALES, setLocale } from "../lib/i18n.svelte";

  const p = $derived(app.profile!);
</script>

<div class="page">
  <div class="card setting">
    <div class="icon"><Icon name="info" size={20} /></div>
    <div class="text">
      <h3>{t("settings.language")}</h3>
      <p class="muted">{t("settings.languageText")}</p>
    </div>
    <div class="control wide">
      <div class="segmented">
        {#each LOCALES as l}
          <button class:active={i18n.locale === l.id} onclick={() => setLocale(l.id)}>{l.name}</button>
        {/each}
      </div>
    </div>
  </div>

  <div class="card setting">
    <div class="icon"><Icon name="zap" size={20} /></div>
    <div class="text">
      <h3>{t("settings.debounce")}</h3>
      <p class="muted">{t("settings.debounceText")}</p>
    </div>
    <div class="control">
      <input type="range" min="1" max="20" bind:value={p.debounce} style:--fill="{((p.debounce - 1) / 19) * 100}%" />
      <span class="value">{p.debounce} ms</span>
    </div>
  </div>

  {#if app.meta?.platform === "linux"}
    <div class="card setting">
      <div class="icon"><Icon name="shield" size={20} /></div>
      <div class="text">
        <h3>{t("settings.perms")}</h3>
        <p class="muted">{t("settings.permsText")}</p>
      </div>
      <div class="control">
        {#if app.status.connected}
          <span class="ok"><Icon name="check" size={16} /> {t("settings.permsGranted")}</span>
        {:else}
          <button class="btn" onclick={installPermissions}>{t("settings.permsInstall")}</button>
        {/if}
      </div>
    </div>
  {/if}

  <div class="card setting danger-zone">
    <div class="icon"><Icon name="rotate" size={20} /></div>
    <div class="text">
      <h3>{t("settings.reset")}</h3>
      <p class="muted">{t("settings.resetText")}</p>
    </div>
    <div class="control">
      <button class="btn danger" disabled={!app.status.connected || app.busy} onclick={factoryReset}>{t("settings.resetButton")}</button>
    </div>
  </div>

  <div class="card setting">
    <div class="icon"><Icon name="keyboard" size={20} /></div>
    <div class="text">
      <h3>{t("settings.about")}</h3>
      <p class="muted">{t("settings.aboutText")}</p>
    </div>
  </div>
</div>

<style>
  .page { display: flex; flex-direction: column; max-width: 860px; border-top: 1px solid var(--border); }
  .setting { display: flex; align-items: center; gap: 18px; padding: 18px 4px; border: 0; border-bottom: 1px solid var(--border); background: transparent; }
  .icon { width: 36px; height: 36px; display: grid; place-items: center; flex: none; border: 1px solid var(--border-strong); color: var(--accent-2); }
  .text { flex: 1; }
  h3 { margin: 0 0 4px; font-size: 14px; }
  p { margin: 0; font-size: 12.5px; line-height: 1.5; }
  .control { width: 220px; flex: none; display: flex; align-items: center; justify-content: flex-end; gap: 12px; }
  .control.wide { width: auto; }
  .value { width: 48px; text-align: right; font-family: var(--mono); }
  .ok { display: flex; align-items: center; gap: 6px; color: var(--success); }
  .danger-zone .icon { border-color: var(--danger); color: var(--danger); }
</style>
