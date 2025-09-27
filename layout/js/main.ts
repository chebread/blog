import mediumZoom from 'medium-zoom';

function initializeStickyNav() {
  const navs = document.querySelectorAll('.home-nav, .about-nav');

  if (navs.length === 0) return;

  window.addEventListener('scroll', () => {
    navs.forEach(nav => {
      if (window.scrollY > 10) {
        nav.classList.add('scrolled');
      } else {
        nav.classList.remove('scrolled');
      }
    });
  });
}


document.addEventListener('DOMContentLoaded', () => {
    initializeStickyNav();

    mediumZoom('.markdown-body img');
});