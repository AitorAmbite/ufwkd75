<script lang="ts">
  import Icon from "../lib/Icon.svelte";
  import {
    app, isSaved, newProfile, openProfile, deleteProfile, importProfile, exportProfile, installPermissions, type Page,
  } from "../lib/state.svelte";
  import { t, tError } from "../lib/i18n.svelte";

  const nav: { id: Page; icon: string }[] = [
    { id: "keys", icon: "keyboard" },
    { id: "light", icon: "sun" },
    { id: "macros", icon: "macro" },
    { id: "settings", icon: "settings" },
  ];
</script>

<aside class="sidebar">
  <div class="brand">
    <div class="logo"><Icon name="keyboard" size={20} /></div>
    <div>
      <div class="name">D75 Studio</div>
      <div class="faint small">Weikav D75</div>
    </div>
  </div>

  <nav>
    {#each nav as item}
      <button class="nav-item" class:active={app.page === item.id} onclick={() => (app.page = item.id)}>
        <Icon name={item.icon} />
        {t(`page.${item.id}.title`)}
      </button>
    {/each}
  </nav>

  <div class="section-head">
    <span class="label">{t("profiles.title")}</span>
    <div class="head-actions">
      <button class="btn ghost sm icon" title={t("profiles.import")} onclick={importProfile}><Icon name="upload" size={15} /></button>
      <button class="btn ghost sm icon" title={t("profiles.export")} disabled={!app.profile} onclick={exportProfile}><Icon name="download" size={15} /></button>
      <button class="btn ghost sm icon" title={t("profiles.new")} onclick={newProfile}><Icon name="plus" size={16} /></button>
    </div>
  </div>

  <div class="profiles">
    {#if app.profileName === ""}
      <div class="profile active unsaved">
        <Icon name="layers" size={15} />
        <span class="pname">{t("profiles.unsaved")}</span>
      </div>
    {/if}
    {#each app.profiles as name (name)}
      <div class="profile" class:active={app.profileName === name}>
        <button class="pbtn" onclick={() => openProfile(name)}>
          <Icon name="folder" size={15} />
          <span class="pname">{name}</span>
          {#if app.profileName === name && !isSaved()}<span class="dot" title={t("profiles.unsavedChanges")}></span>{/if}
        </button>
        <button class="del" title={t("profiles.delete")} onclick={() => deleteProfile(name)}><Icon name="trash" size={14} /></button>
      </div>
    {/each}
    {#if app.profiles.length === 0}
      <p class="faint small empty">{t("profiles.empty")}</p>
    {/if}
  </div>

  <div class="device card" class:ok={app.status.connected}>
    <div class="device-row">
      <span class="status-dot"></span>
      <div class="device-text">
        <div>{t(app.status.connected ? "device.connected" : app.status.needsPermission ? "device.noPermission" : "device.disconnected")}</div>
        <div class="faint small ellipsis">{tError(app.status.detail)}</div>
      </div>
      <Icon name="usb" size={16} class="faint" />
    </div>
    {#if app.status.needsPermission}
      <button class="btn sm primary full" onclick={installPermissions}>
        <Icon name="shield" size={15} /> {t("device.grant")}
      </button>
    {/if}
  </div>
</aside>

<style>
  .sidebar {
    width: 232px; flex: none; display: flex; flex-direction: column;
    background: var(--bg-elev); border-right: 1px solid var(--border);
  }
  .brand { display: flex; align-items: center; gap: 12px; padding: 16px; border-bottom: 1px solid var(--border); }
  .logo { width: 32px; height: 32px; display: grid; place-items: center; background: var(--accent); color: #0d0d0e; }
  .name { font-weight: 700; font-size: 14px; letter-spacing: .02em; text-transform: uppercase; }
  .small { font-size: 11px; }

  nav { display: flex; flex-direction: column; border-bottom: 1px solid var(--border); padding: 8px 0; }
  .nav-item {
    display: flex; align-items: center; gap: 11px; height: 36px; padding: 0 16px;
    border: 0; border-left: 2px solid transparent; background: transparent; color: var(--text-2); cursor: pointer; text-align: left;
  }
  .nav-item:hover { background: var(--surface); color: var(--text); }
  .nav-item.active { border-left-color: var(--accent); background: var(--surface); color: var(--text); }

  .section-head { display: flex; align-items: center; justify-content: space-between; padding: 12px 10px 6px 16px; }
  .head-actions { display: flex; gap: 2px; }
  .profiles { flex: 1; overflow-y: auto; display: flex; flex-direction: column; min-height: 0; }
  .profile { display: flex; align-items: center; color: var(--text-2); border-left: 2px solid transparent; }
  .profile.unsaved { gap: 10px; height: 32px; padding: 0 14px; font-style: italic; }
  .profile:hover { background: var(--surface); }
  .profile.active { background: var(--surface); color: var(--text); border-left-color: var(--text-2); }
  .pbtn {
    flex: 1; min-width: 0; display: flex; align-items: center; gap: 10px; height: 32px; padding: 0 14px;
    border: 0; background: transparent; cursor: pointer; text-align: left;
  }
  .pname { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .dot { width: 6px; height: 6px; background: var(--warning); flex: none; }
  .del {
    width: 26px; height: 26px; margin-right: 6px; border: 0; background: transparent;
    color: var(--text-3); cursor: pointer; display: grid; place-items: center; opacity: 0;
  }
  .profile:hover .del { opacity: 1; }
  .del:hover { color: var(--danger); }
  .empty { padding: 4px 16px; margin: 0; }

  .device { padding: 12px 16px; display: flex; flex-direction: column; gap: 10px; border: 0; border-top: 1px solid var(--border); background: transparent; }
  .device-row { display: flex; align-items: center; gap: 10px; }
  .device-text { flex: 1; min-width: 0; }
  .ellipsis { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-family: var(--mono); }
  .status-dot { width: 8px; height: 8px; flex: none; background: var(--danger); }
  .device.ok .status-dot { background: var(--success); }
  .full { width: 100%; }
</style>
