// some constants used as state values keys for consistency
// and for reflecting the application concept in code 
const SELECTED_NAV_OPTION = "selected-nav-option";
const NAVIGATION_STACK = "navigation-stack";

class StateManager {
  _stateMap = {};
  _listeners = {};

  get(key) {
    return this._stateMap[key];
  }

  set(key, value) {
    this._stateMap[key] = value;
    if (this._listeners[key]) {
      for (const func of this._listeners[key]) func();
    }
  }

  addListener(key, func) {
    if (!this._listeners[key]) {
      this._listeners[key] = [];
    }
    this._listeners[key].push(func);
  }
}
