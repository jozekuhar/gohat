import "./view.css";
import Alpine from "alpinejs";
import focus from "@alpinejs/focus";
import modal from "./modal";
import htmx from "htmx.org";
import toast from "./toast";

declare global {
    interface Window {
        htmx: typeof htmx;
        Alpine: typeof Alpine;
    }
}

Alpine.plugin(focus);

Alpine.data("modal", modal);
Alpine.data("toast", toast);

window.htmx = htmx;
window.Alpine = Alpine;

Alpine.start();
