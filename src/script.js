document.addEventListener("mousemove", (event) => {
  const eyes = document.querySelector(".eyes");
  const pupil = document.getElementById("pupils");

  const mouseX = event.clientX;
  const mouseY = event.clientY;

  const eyesRect = eyes.getBoundingClientRect();
  const eyesX = eyesRect.left + eyesRect.width / 2;
  const eyesY = eyesRect.top + eyesRect.height / 2;
  const angle = Math.atan2(mouseY - eyesY, mouseX - eyesX);
  const distance = Math.min(eyesRect.width / 4, eyesRect.height / 4);

  const pupilX = Math.cos(angle) * distance;
  const pupilY = Math.sin(angle) * distance;

  pupil.style.transform = `translate(-50%, -50%) translate(${pupilX}px, ${pupilY}px)`;
});
