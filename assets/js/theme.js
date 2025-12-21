const storageKey = 'kafscale-theme';
const root = document.documentElement;
const toggle = document.querySelector('.theme-toggle');
const sunLabel = toggle ? toggle.querySelector('.sun') : null;
const moonLabel = toggle ? toggle.querySelector('.moon') : null;
const navToggle = document.querySelector('.nav-toggle');
const navLinks = document.getElementById('nav-links');
const banner = document.querySelector('.site-banner');
const bannerClose = document.querySelector('.banner-close');
const bannerKey = 'kafscale-banner-dismissed';

const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
const savedTheme = localStorage.getItem(storageKey);
const initialTheme = savedTheme || (prefersDark ? 'dark' : 'light');
root.setAttribute('data-theme', initialTheme);

function syncToggleLabels(theme) {
  if (!sunLabel || !moonLabel) return;
  if (theme === 'dark') {
    sunLabel.textContent = 'Light';
    moonLabel.textContent = '';
  } else {
    sunLabel.textContent = '';
    moonLabel.textContent = 'Dark';
  }
}

function setTheme(next) {
  root.setAttribute('data-theme', next);
  localStorage.setItem(storageKey, next);
  syncToggleLabels(next);
}

syncToggleLabels(initialTheme);

if (toggle) {
  toggle.addEventListener('click', () => {
    const current = root.getAttribute('data-theme');
    setTheme(current === 'dark' ? 'light' : 'dark');
  });
}

if (banner && bannerClose) {
  const hiddenUntil = Number(localStorage.getItem(bannerKey));
  if (hiddenUntil && Date.now() < hiddenUntil) {
    banner.remove();
  } else {
    bannerClose.addEventListener('click', () => {
      const next = Date.now() + 24 * 60 * 60 * 1000;
      localStorage.setItem(bannerKey, String(next));
      banner.remove();
    });
  }
}

if (navToggle && navLinks) {
  navToggle.addEventListener('click', () => {
    const expanded = navToggle.getAttribute('aria-expanded') === 'true';
    navToggle.setAttribute('aria-expanded', String(!expanded));
    navLinks.classList.toggle('open');
  });
}
