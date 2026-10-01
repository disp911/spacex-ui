<template>
  <div class="waves" aria-hidden="true">
    <div class="waves__glow" />
    <div class="waves__band">
      <div class="waves__layer waves__layer--a">
        <svg viewBox="0 0 2880 340" preserveAspectRatio="none">
          <path d="M0 196 C 240 138 480 254 720 196 C 960 138 1200 254 1440 196 C 1680 138 1920 254 2160 196 C 2400 138 2640 254 2880 196 L2880 340 L0 340 Z" fill="var(--wave-1)" />
        </svg>
      </div>
      <div class="waves__layer waves__layer--b">
        <svg viewBox="0 0 2880 340" preserveAspectRatio="none">
          <path d="M0 232 C 120 190 360 274 480 232 C 600 190 840 274 960 232 C 1080 190 1320 274 1440 232 C 1560 190 1800 274 1920 232 C 2040 190 2280 274 2400 232 C 2520 190 2760 274 2880 232 L2880 340 L0 340 Z" fill="var(--wave-2)" />
        </svg>
      </div>
      <div class="waves__layer waves__layer--c">
        <svg viewBox="0 0 2880 340" preserveAspectRatio="none">
          <path d="M0 268 C 90 242 270 294 360 268 C 450 242 630 294 720 268 C 810 242 990 294 1080 268 C 1170 242 1350 294 1440 268 C 1530 242 1710 294 1800 268 C 1890 242 2070 294 2160 268 C 2250 242 2430 294 2520 268 C 2610 242 2790 294 2880 268 L2880 340 L0 340 Z" fill="var(--wave-3)" />
        </svg>
      </div>
    </div>
  </div>
</template>

<style scoped>
.waves {
  position: fixed;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}
.waves__glow {
  position: absolute;
  inset: 0;
  background: var(--glow);
}
.waves__band {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: min(440px, 52vh);
  overflow: hidden;
}
.waves__layer {
  position: absolute;
  left: 0;
  width: 200%;
  height: 100%;
}
.waves__layer svg {
  display: block;
  width: 100%;
  height: 100%;
}
.waves__layer--a { bottom: -6px; opacity: 0.3; animation: wave-l 26s linear infinite; }
.waves__layer--b { bottom: -14px; opacity: 0.2; animation: wave-r 17s linear infinite; }
.waves__layer--c { bottom: -20px; opacity: var(--wave-3-opacity); animation: wave-l 11s linear infinite; }
@keyframes wave-l {
  from { transform: translate3d(0, 0, 0); }
  to { transform: translate3d(-50%, 0, 0); }
}
@keyframes wave-r {
  from { transform: translate3d(-50%, 0, 0); }
  to { transform: translate3d(0, 0, 0); }
}
</style>
