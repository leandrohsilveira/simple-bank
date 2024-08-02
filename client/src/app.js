import "htmx.org";
import "./message.js";
import { handleShowMessage } from "./message.js";
import { handleHistoryChange } from "./history.js";

document.body.addEventListener("showMessage", handleShowMessage);
document.body.addEventListener("htmx:pushedIntoHistory", handleHistoryChange);
document.body.addEventListener("htmx:historyRestore", handleHistoryChange);
document.body.addEventListener("onRedirect", handleHistoryChange);
