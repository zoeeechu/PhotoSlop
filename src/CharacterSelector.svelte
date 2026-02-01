<script lang="ts">
  const characters = ['ava', 'cubert', 'potat', 'azra', "eepy", "frog", "iconokeeb", "void", "zoe"];
  export let selected: string | null = null;
  export let available: string[] = [];
  export let onSelect: (character: string) => void;
</script>

<div class="grid grid-cols-3 gap-4">
  {#each characters as character}
    {@const isSelected = selected === character}
    {@const isAvailable = available.includes(character) || isSelected}
    {@const isOtherSelected = selected && !isSelected}
    <button 
      type="button"
      aria-label="Select {character}"
      disabled={!isAvailable}
      on:click={() => { selected = character; onSelect(character); }}
    >
      <img 
        src="/images/characters/{character}.png" 
        alt={character} 
        class="w-24 h-24 object-contain transition-transform {isAvailable ? 'hover:scale-110' : 'opacity-25 pointer-events-none'} {isAvailable && isOtherSelected ? 'opacity-50' : ''}"
      />
    </button>
  {/each}
</div>
