const AppState = new StateManager();

document.addEventListener("DOMContentLoaded", () => {
  AppState.set(SELECTED_NAV_OPTION, 0);
  AppState.set(NAVIGATION_STACK, []);
  AppState.set(MY_CART, [...localStorage.getItem("mycart").split(",").map(i => Number(i))]);
});

document.addEventListener("keydown", (e) => {
  if (e.key === "F12") {
    e.preventDefault();
  }
});
