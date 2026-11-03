<script lang="ts" module>
  // Iconos de trazo estilo Lucide (24x24)
  const paths: Record<string, string> = {
    keyboard: "M4 5h16a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2Z M6 9h.01 M10 9h.01 M14 9h.01 M18 9h.01 M6 13h.01 M18 13h.01 M10 13h4 M7 16h10",
    sun: "M12 3v2 M12 19v2 M5.6 5.6l1.4 1.4 M17 17l1.4 1.4 M3 12h2 M19 12h2 M5.6 18.4 7 17 M17 7l1.4-1.4 M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8Z",
    settings: "M4 6h10 M18 6h2 M4 12h4 M12 12h8 M4 18h12 M20 18h0 M16 4v4 M10 10v4 M18 16v4",
    save: "M5 3h11l3 3v13a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V3Z M8 3v5h7V3 M8 21v-7h8v7",
    zap: "M13 2 4 14h7l-1 8 9-12h-7l1-8Z",
    target: "M12 12m-9 0a9 9 0 1 0 18 0 9 9 0 1 0-18 0 M12 12m-5 0a5 5 0 1 0 10 0 5 5 0 1 0-10 0 M12 12h.01",
    plus: "M12 5v14 M5 12h14",
    trash: "M3 6h18 M8 6V4h8v2 M6 6l1 14h10l1-14 M10 11v6 M14 11v6",
    copy: "M9 9h11v11H9z M5 15H4V4h11v1",
    rotate: "M3 12a9 9 0 1 0 3-6.7L3 8 M3 3v5h5",
    check: "M5 12l5 5L20 7",
    usb: "M12 2v14 M12 2l-3 3 M12 2l3 3 M8 9v3l4 3 4-3V8 M12 16a2 2 0 1 0 0 4 2 2 0 0 0 0-4Z",
    shield: "M12 3l8 3v6c0 5-3.5 8-8 9-4.5-1-8-4-8-9V6l8-3Z M9 12l2 2 4-4",
    x: "M6 6l12 12 M18 6 6 18",
    search: "M11 11m-7 0a7 7 0 1 0 14 0 7 7 0 1 0-14 0 M21 21l-5-5",
    brush: "M9.5 14.5 18 6a2.1 2.1 0 0 1 3 3l-8.5 8.5 M7 14c-2 0-3 1.5-3 3.5 0 1.5-1 2.5-2 2.5 1.5 1 3 1.5 4.5 1.5 2.5 0 4.5-2 4.5-4.5 0-1.5-1.5-3-4-3Z",
    pipette: "M14 4l6 6 M17 7l2-2a2.1 2.1 0 0 0-3-3l-2 2 M15 9 6 18l-3 3-1-1 3-3 9-9",
    eraser: "M7 21h12 M3.5 15.5l9-9a2 2 0 0 1 2.8 0l3.2 3.2a2 2 0 0 1 0 2.8L11 20H7.5l-4-4a2 2 0 0 1 0-.5Z M9 10l5 5",
    sparkles: "M12 3l1.8 4.7L18.5 9.5l-4.7 1.8L12 16l-1.8-4.7L5.5 9.5l4.7-1.8L12 3Z M19 15l.8 2.2L22 18l-2.2.8L19 21l-.8-2.2L16 18l2.2-.8L19 15Z",
    layers: "M12 3l9 5-9 5-9-5 9-5Z M3 13l9 5 9-5",
    alert: "M12 3l10 18H2L12 3Z M12 10v4 M12 17h.01",
    info: "M12 12m-9 0a9 9 0 1 0 18 0 9 9 0 1 0-18 0 M12 11v5 M12 8h.01",
    folder: "M3 6a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6Z",
    arrow: "M5 12h14 M13 6l6 6-6 6",
    loader: "M12 3a9 9 0 1 0 9 9",
    palette: "M12 3a9 9 0 1 0 0 18c1 0 1.5-.8 1.5-1.5 0-.4-.2-.8-.4-1.1-.3-.3-.4-.6-.4-1 0-.8.7-1.4 1.5-1.4H16a5 5 0 0 0 5-5C21 6.6 17 3 12 3Z M7.5 11.5h.01 M9.5 7.5h.01 M14.5 7.5h.01 M17 11h.01",
  };
</script>

<script lang="ts">
  let { name, size = 18, class: cls = "" }: { name: string; size?: number; class?: string } = $props();
</script>

<svg class={cls} width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor"
  stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
  <path d={paths[name] ?? ""} />
</svg>
