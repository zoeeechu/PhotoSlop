<script lang="ts">
  export let status = '';
  export let statusVisible = false;
  export let statusColor = 'green';

  let statusTimeout: ReturnType<typeof setTimeout> | undefined;
  $: {
    if (statusVisible) {
      if (statusTimeout !== undefined) {
        clearTimeout(statusTimeout);
      }
      statusTimeout = setTimeout(() => {
        statusVisible = false;
      }, 4000);
    }
  }
</script>

{#if statusVisible}
  <div
    class="w-full h-7 text-white flex items-center px-4"
    style="background: {statusColor};"
  >
    <div>{status}</div>
  </div>
{/if}