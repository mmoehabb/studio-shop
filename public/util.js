function switchMode() {
  const html = document.getElementsByTagName("html")[0]
  const theme = html.getAttribute("theme")
  html.setAttribute("theme", theme !== "dark" ? "dark" : "light")
}   

// Get notified when it's opened/closed or orientation changes
window.addEventListener('devtoolschange', event => {
    console.log('Is DevTools open:', event.detail.isOpen);
    console.log('DevTools orientation:', event.detail.orientation);
});
