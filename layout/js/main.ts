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

function initializeBackLink() {
  const backLinks = document.querySelectorAll('#back');

  backLinks.forEach(link => {
    link.addEventListener('click', (event) => {
      event.preventDefault();
      history.back();
    });
  });
}

document.addEventListener('DOMContentLoaded', () => {
    initializeStickyNav();
    initializeBackLink();
    mediumZoom('.markdown-body img');
});