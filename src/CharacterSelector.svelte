<script lang="ts">
  const characters = ['ava', 'cubert', 'potat'];
</script>

<div class="flex gap-4">
  {#each characters as character}
    <img src="/images/characters/{character}.png" alt={character} 
    class="w-24 h-24 object-contain  hover:scale-110 hover:outline-blue-500 transition-transform cursor-pointer" />
  {/each}
</div>
