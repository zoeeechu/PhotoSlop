<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import p5 from "p5";

  export let mask: { x: number; y: number }[] = [];
  export let timeLimit: number = 5.0;
  export let direction: string = "";
  export let onComplete: (result: { x: number; y: number }[]) => void;

  let canvasContainer: HTMLDivElement;
  let p5Instance: p5;
  let drawnPoints: { x: number; y: number }[] = [];
  let isDrawing = false;
  let timeRemaining = timeLimit;
  let gameActive = true;
  let scaledMask: { x: number; y: number }[] = [];
  let scale = 1;
  let scaledDrawnPoints: { x: number; y: number }[] = [];

  onMount(() => {
    const sketch = (p: p5) => {
      p.setup = () => {
        const canvasWidth = window.innerWidth * 0.95;
        const canvasHeight = canvasWidth * 0.75;
        const canvas = p.createCanvas(canvasWidth, canvasHeight);
        canvas.style('display', 'block');
        
        // Scale based on 1920 width reference
        scale = canvasWidth / 1920;
        
        scaledMask = mask.map(point => ({
          x: point.x * scale,
          y: point.y * scale
        }));
        
        p.background(255);
      };

      p.draw = () => {
        if (!gameActive) return;

        if (isDrawing) {
          p.stroke(0, 0, 255);
          p.strokeWeight(8);
          p.line(p.pmouseX, p.pmouseY, p.mouseX, p.mouseY);
          drawnPoints.push({ x: p.mouseX, y: p.mouseY });
        }

        timeRemaining -= p.deltaTime / 1000;
        if (timeRemaining <= 0) {
          endGame(p);
        }
      };

      p.mousePressed = () => {
        if (gameActive && p.mouseX >= 0 && p.mouseX <= p.width && p.mouseY >= 0 && p.mouseY <= p.height) {
          isDrawing = true;
          drawnPoints = [{ x: p.mouseX, y: p.mouseY }];
        }
      };

      p.mouseReleased = () => {
        isDrawing = false;
      };

      p.touchStarted = () => {
        if (gameActive) {
          isDrawing = true;
          drawnPoints = [{ x: p.mouseX, y: p.mouseY }];
          return false;
        }
      };

      p.touchEnded = () => {
        isDrawing = false;
        return false;
      };
    };

    p5Instance = new p5(sketch, canvasContainer);

    const timer = setInterval(() => {
      if (gameActive && timeRemaining > 0) {
        timeRemaining -= 0.1;
      }
    }, 100);

    return () => clearInterval(timer);
  });

  function drawMaskOutline(p: p5) {
    p.stroke(255, 0, 0, 150);
    p.strokeWeight(3);
    p.noFill();
    p.beginShape();
    for (const point of scaledMask) {
      p.vertex(point.x, point.y);
    }
    p.endShape(p.CLOSE);
  }

  function calculateMatchPercentage(): number {
    if (drawnPoints.length === 0) return 0;

    let matchCount = 0;
    const threshold = 20;

    for (const drawnPoint of drawnPoints) {
      for (const maskPoint of scaledMask) {
        const distance = Math.sqrt(
          Math.pow(drawnPoint.x - maskPoint.x, 2) + 
          Math.pow(drawnPoint.y - maskPoint.y, 2)
        );
        if (distance < threshold) {
          matchCount++;
          break;
        }
      }
    }

    return Math.min(100, (matchCount / scaledMask.length) * 100);
  }

  

  function endGame(p: p5) {
    if (!gameActive) return;
    gameActive = false;
    
    p.textAlign(p.CENTER, p.CENTER);
    p.textSize(120);
    p.textFont('Darumadrop One');
    p.fill(245, 222, 179);
    p.stroke(255);
    p.strokeWeight(8);
    p.text("DONE!", p.width / 2, p.height / 2);
    
    const percentage = calculateMatchPercentage();

    scaledDrawnPoints = []

    for (const drawnPoint of drawnPoints) {
      let newPoint = drawnPoint;
      newPoint.x = newPoint.x / scale;
      newPoint.y = newPoint.y / scale;
      scaledDrawnPoints.push(newPoint)
    }

    onComplete(scaledDrawnPoints);
  }

  onDestroy(() => {
    if (p5Instance) {
      p5Instance.remove();
    }
  });
</script>

<div class="fixed inset-0 bg-white z-50 flex flex-col items-center justify-center">
  <div class="text-center mb-4">
    <h2 class="text-4xl font-bold mb-2">{direction}</h2>
    <div class="text-2xl">Time: {Math.max(0, timeRemaining).toFixed(1)}s</div>
  </div>
  <div bind:this={canvasContainer} class="border-4 border-black"></div>
</div>
