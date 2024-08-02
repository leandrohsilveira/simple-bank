import { assert } from "./util";

/**
 *
 * @param {Event} event
 */
export function handleShowMessage(event) {
  assert(event instanceof CustomEvent, "event is not a CustomEvent");
  assertMessage(event.detail);
  addMessage(event.detail);
}

/**
 *
 * @param {import('./types.ts').Message} message
 */
export function addMessage(message) {
  const duration = 5000;

  const messageEl = document.createElement("li");
  messageEl.classList.add(message.type);

  const messageText = document.createElement("strong");
  messageText.textContent = message.message;
  messageEl.appendChild(messageText);

  if (message.description) {
    const messageDescription = document.createElement("p");
    messageDescription.textContent = message.description;
    messageEl.appendChild(messageDescription);
  }

  const messages = document.getElementById("messages");
  assert(messages !== null, "messages element not found");

  messages.appendChild(messageEl);

  setTimeout(() => {
    messageEl.classList.add("show");
  });

  setTimeout(() => {
    messageEl.classList.remove("show");
  }, duration);

  setTimeout(() => {
    messageEl.remove();
  }, duration + 500);
}

/**
 *
 * @param {unknown} message
 * @returns {asserts message is import('./types.ts').Message}
 */
export function assertMessage(message) {
  assert(typeof message === "object", "message must be an object");
  assert(message !== null, "message must not be null");
  assert("type" in message, "message must have a 'type' property");
  assert(typeof message.type === "string", "message.type must be a string");
  assert("message" in message, "message must have a 'message' property");
  assert(
    typeof message.message === "string",
    "message.message must be a string"
  );
}
