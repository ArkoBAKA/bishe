<template>
  <router-view />
</template>

<style>
:root {
  color-scheme: light;
  --bg-page:
    radial-gradient(
      900px 360px at 0% 0%,
      rgba(99, 102, 241, 0.08),
      transparent 55%
    ),
    radial-gradient(
      720px 280px at 100% 0%,
      rgba(56, 189, 248, 0.08),
      transparent 50%
    ),
    linear-gradient(180deg, #f8fafc 0%, #f5f7fb 48%, #f3f6fb 100%);
  --card-bg: rgba(255, 255, 255, 0.92);
  --card-border: rgba(226, 232, 240, 0.92);
  --card-shadow: 0 10px 30px rgba(15, 23, 42, 0.06);
  --card-shadow-hover: 0 18px 40px rgba(15, 23, 42, 0.1);
  --text-main: #0f172a;
  --text-sub: #64748b;
  --text-soft: #94a3b8;
  --brand: #4f46e5;
}

#app {
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue",
    Arial, "Noto Sans", "Apple Color Emoji", "Segoe UI Emoji",
    "Segoe UI Symbol", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--text-main);
  min-height: 100vh;
}

html,
body {
  height: 100%;
  margin: 0;
}

body {
  background: var(--bg-page);
}

button,
input,
select,
textarea {
  font: inherit;
}

button,
a,
input,
select,
textarea {
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease,
    box-shadow 0.2s ease,
    color 0.2s ease,
    transform 0.2s ease;
}

button:hover {
  transform: translateY(-1px);
}

button:active {
  transform: translateY(0);
}

button:focus-visible,
input:focus-visible,
select:focus-visible,
textarea:focus-visible,
[contenteditable="true"]:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.12);
  border-color: rgba(79, 70, 229, 0.42);
}

::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

::-webkit-scrollbar-thumb {
  background: rgba(148, 163, 184, 0.45);
  border-radius: 999px;
  border: 2px solid transparent;
  background-clip: padding-box;
}

::-webkit-scrollbar-track {
  background: transparent;
}

* {
  box-sizing: border-box;
}
</style>
