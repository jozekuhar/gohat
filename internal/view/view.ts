import "./view.css";
import "htmx.org";
import Alpine from "alpinejs";
import CounterComponent from "./counter";

Alpine.data("CounterComponent", CounterComponent);

Alpine.start();
