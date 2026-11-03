<script lang="ts">
  import Icon from "../lib/Icon.svelte";
  import { dialog, toasts } from "../lib/state.svelte";
  import { t } from "../lib/i18n.svelte";

  let value = $state("");
  let input: HTMLInputElement | undefined = $state();

  $effect(() => {
    if (dialog.current) {
      value = dialog.current.input ?? "";
      queueMicrotask(() => input?.focus());
    }
  });

  function close(result: string | boolean | null) {
    dialog.current?.resolve(result);
    dialog.current = null;
  }

  function confirm() {
    if (dialog.current?.input !== undefined) close(value);
    else close(true);
  }

  const icons = { ok: "check", error: "alert", info: "info" };
</script>

{#if dialog.current}
  {@const d = dialog.current}
  <div class="backdrop" role="presentation" onmousedown={(e) => e.target === e.currentTarget && close(null)}>
    <div class="dialog card" role="dialog" aria-modal="true">
      <h3>{d.title}</h3>
      {#if d.message}<p class="muted">{d.message}</p>{/if}
      {#if d.input !== undefined}
        <input class="input" bind:this={input} bind:value
          onkeydown={(e) => { if (e.key === "Enter") confirm(); if (e.key === "Escape") close(null); }} />
      {/if}
      <div class="actions">
        <button class="btn ghost" onclick={() => close(null)}>{t("dialog.cancel")}</button>
        <button class="btn" class:primary={!d.danger} class:danger={d.danger} onclick={confirm}>{d.confirm}</button>
      </div>
    </div>
  </div>
{/if}

<div class="toasts">
  {#each toasts as t (t.id)}
    <div class="toast {t.kind}">
      <Icon name={icons[t.kind]} size={16} />
      <span>{t.text}</span>
    </div>
  {/each}
</div>

<style>
  .backdrop { position: fixed; inset: 0; z-index: 50; display: grid; place-items: center; background: #000000b3; }
  .dialog { width: 420px; padding: 20px; display: flex; flex-direction: column; gap: 14px; border-color: var(--border-strong); }
  h3 { margin: 0; font-size: 15px; text-transform: uppercase; letter-spacing: .04em; }
  p { margin: 0; font-size: 13px; line-height: 1.5; }
  .actions { display: flex; justify-content: flex-end; gap: 6px; margin-top: 4px; }

  .toasts { position: fixed; bottom: 20px; right: 20px; z-index: 60; display: flex; flex-direction: column; gap: 6px; align-items: flex-end; }
  .toast {
    display: flex; align-items: center; gap: 10px; padding: 10px 14px; max-width: 440px;
    background: var(--surface-2); border: 1px solid var(--border-strong); border-left: 3px solid var(--success);
  }
  .toast.ok :global(svg) { color: var(--success); }
  .toast.info { border-left-color: var(--accent); }
  .toast.info :global(svg) { color: var(--accent-2); }
  .toast.error { border-left-color: var(--danger); }
  .toast.error :global(svg) { color: var(--danger); }
</style>
