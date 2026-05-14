import "./view.css";
import Alpine from "alpinejs";
import focus from "@alpinejs/focus";
import modalComponent from "./modal";
import counterComponent from "./counter";
import htmx from "htmx.org";

declare global {
    interface Window {
        htmx: typeof htmx;
        Alpine: typeof Alpine;
    }
}

Alpine.plugin(focus);

Alpine.data("modal", modalComponent);
Alpine.data("counter", counterComponent);

window.htmx = htmx;
window.Alpine = Alpine;

Alpine.start();
