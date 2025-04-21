const AppState = new StateManager();

document.addEventListener("DOMContentLoaded", () => {
  AppState.set(SELECTED_NAV_OPTION, 0);
  AppState.set(NAVIGATION_STACK, []);
  if (AppState.get(MY_CART) === undefined) {
    if (localStorage.getItem("mycart"))
      AppState.set(MY_CART, [...localStorage.getItem("mycart").split(",").map(i => Number(i))]);
    else
      AppState.set(MY_CART, []);
  }
});

document.addEventListener("keydown", (e) => {
  if (e.key === "F12") {
    e.preventDefault();
  }
});
