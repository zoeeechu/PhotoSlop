<script lang="ts">
  const characters = ['ava', 'cubert', 'potat'];
  export let selected: string | null = null;
  export let onSelect: (character: string) => void;
</script>

<div class="flex gap-4">
  {#each characters as character}
    <img 
      src="/images/characters/{character}.png" 
      alt={character} 
      class="w-24 h-24 object-contain hover:scale-110 transition-transform cursor-pointer {selected === character ? 'outline outline-4 outline-blue-500' : ''}"
      on:click={() => { selected = character; onSelect(character); }}
    />
  {/each}
</div>
