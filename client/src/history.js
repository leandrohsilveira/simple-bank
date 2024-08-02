import { assert } from "./util";

/**
 *
 * @param {Event} event
 */
export function handleHistoryChange(event) {
  assert(event instanceof CustomEvent, "event is not a CustomEvent");
  assertHistoryEvent(event.detail);
  updateNavActive(event.detail.path);
}

/**
 *
 * @param {string} pathname
 */
export function updateNavActive(pathname) {
  const fragments = pathname.split("?");
  const path = fragments[0];
  const nav = document.getElementById("nav");
  assert(nav !== null, "nav element not found");
  nav.querySelectorAll(`a[href]`).forEach((el) => {
    assert(el instanceof HTMLAnchorElement, "el is not an HTMLAnchorElement");
    if (el.getAttribute("href") === path) el.classList.add("active");
    else el.classList.remove("active");
  });
}

/**
 *
 * @param {unknown} detail
 * @returns {asserts detail is { path: string }}
 */
function assertHistoryEvent(detail) {
  assert(typeof detail === "object", "detail is not an object");
  assert(detail !== null, "detail is null");
  assert("path" in detail, "detail is missing 'path' property");
  assert(typeof detail.path === "string", "detail.path is not a string");
}
