<script lang="ts">
  import { onMount } from "svelte";
  import Icon from "./lib/Icon.svelte";
  import Sidebar from "./components/Sidebar.svelte";
  import Overlays from "./components/Overlays.svelte";
  import KeysPage from "./pages/KeysPage.svelte";
  import LightingPage from "./pages/LightingPage.svelte";
  import SettingsPage from "./pages/SettingsPage.svelte";
  import MacrosPage from "./pages/MacrosPage.svelte";
  import { app, init, isApplied, isSaved, saveProfile, applyToKeyboard, errorText } from "./lib/state.svelte";
  import { paint } from "./lib/light.svelte";
  import { t } from "./lib/i18n.svelte";

  let error = $state("");

  onMount(() => {
    init().catch((e) => (error = errorText(e)));
  });

  const applied = $derived(app.profile ? isApplied() : true);
  const saved = $derived(app.profile ? isSaved() : true);

  // En modo detección las teclas son para identificar, no para la interfaz
  function onKeydown(e: KeyboardEvent) {
    const tag = (document.activeElement as HTMLElement | null)?.tagName;
    if (app.detect && tag !== "INPUT") e.preventDefault();
  }
</script>

<svelte:window onkeydown={onKeydown} onmouseup={() => (paint.painting = false)} />

<div class="shell">
  <Sidebar />
  <div class="content">
    <header>
      <div>
        <h1>{t(`page.${app.page}.title`)}</h1>
        <p class="muted">{t(`page.${app.page}.subtitle`)}</p>
      </div>
      <div class="actions">
        <button class="btn" class:detecting={app.detect} onclick={() => (app.detect = !app.detect)}
          title={t("action.detectHint")}>
          <Icon name="target" size={16} />
          {app.detect ? t("action.detecting") : t("action.detect")}
        </button>
        <div class="split">
          <button class="btn" onclick={() => saveProfile(false)} title={app.profileName ? t("action.saveNamed", { name: app.profileName }) : t("action.saveNew")}>
            <Icon name="save" size={16} /> {t("action.save")}
            {#if app.profileName && !saved}<span class="dot warn"></span>{/if}
          </button>
          <button class="btn icon" onclick={() => saveProfile(true)} title={t("action.saveAs")}><Icon name="copy" size={15} /></button>
        </div>
        <button class="btn primary" disabled={!app.status.connected || app.busy || !app.profile} onclick={applyToKeyboard}>
          {#if app.busy}<Icon name="loader" size={16} class="spin" />{:else}<Icon name="zap" size={16} />{/if}
          {t("action.apply")}
          {#if !applied}<span class="dot"></span>{/if}
        </button>
      </div>
    </header>

    <main>
      {#if error}
        <div class="card fatal"><Icon name="alert" /> {error}</div>
      {:else if app.profile && app.meta}
        {#if app.page === "keys"}<KeysPage />
        {:else if app.page === "light"}<LightingPage />
        {:else if app.page === "macros"}<MacrosPage />
        {:else}<SettingsPage />{/if}
      {/if}
    </main>

    {#if app.profile && !applied}
      <div class="pending">{t("pending")}</div>
    {/if}
  </div>
</div>

<Overlays />

<style>
  .shell { display: flex; height: 100%; }
  .content { flex: 1; min-width: 0; display: flex; flex-direction: column; position: relative; }
  header {
    display: flex; align-items: center; justify-content: space-between; gap: 20px;
    padding: 16px 24px; border-bottom: 1px solid var(--border);
  }
  h1 { margin: 0; font-size: 18px; text-transform: uppercase; letter-spacing: .04em; }
  header p { margin: 2px 0 0; font-size: 12.5px; }
  .actions { display: flex; align-items: center; gap: 8px; }
  .split { display: flex; }
  .split .btn:last-child { border-left: 0; }
  .detecting { background: var(--success) !important; border-color: var(--success) !important; color: #0d0d0e; }
  .dot { width: 7px; height: 7px; background: #0d0d0e; }
  .dot.warn { background: var(--warning); }
  main { flex: 1; min-height: 0; padding: 20px 24px; overflow-y: auto; }
  .fatal { padding: 16px; display: flex; gap: 10px; align-items: center; color: var(--danger); }
  .pending {
    position: absolute; bottom: 0; left: 0; right: 0; padding: 6px 24px;
    background: var(--warning); color: #0d0d0e; font-size: 12px; font-weight: 600; pointer-events: none;
  }
</style>
