<script lang="ts">
  const characters = ['ava', 'cubert', 'potat', 'azra', "eepy", "frog", "iconokeeb", "void", "zoe"];
  export let selected: string | null = null;
  export let onSelect: (character: string) => void;
</script>

<div class="grid grid-cols-3 gap-4">
  {#each characters as character}
    <button 
      type="button"
      aria-label="Select {character}"
      on:click={() => { selected = character; onSelect(character); }}
    >
      <img 
        src="/images/characters/{character}.png" 
        alt={character} 
        class="w-24 h-24 object-contain hover:scale-110 transition-transform {selected && selected !== character ? 'opacity-50' : ''}"
      />
    </button>
  {/each}
</div>
