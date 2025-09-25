function initializeStickyNav() {
  const navs = document.querySelectorAll('.home-nav, .about-nav');

  // 만약 선택된 nav가 하나도 없다면 함수를 종료합니다.
  if (navs.length === 0) return;

  window.addEventListener('scroll', () => {
    // 선택된 모든 nav 요소에 대해 반복문을 실행합니다.
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
});