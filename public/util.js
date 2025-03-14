function switchMode() {
  const html = document.getElementsByTagName("html")[0]
  const theme = html.getAttribute("theme")
  html.setAttribute("theme", theme !== "dark" ? "dark" : "light")
}   

function detectDevtools(event) {
  if (event?.detail.isOpen || devtools.isOpen) {
    window.location.replace("https://static.planetminecraft.com/files/resource_media/screenshot/1234/Hacker_Detected_3389798.jpg");
  }
}

// Get notified when it's opened/closed or orientation changes
window.addEventListener('devtoolschange', (event) => {
  console.log(event)
  detectDevtools(event);
});

detectDevtools();
