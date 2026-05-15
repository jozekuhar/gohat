import { AlpineComponent } from "alpinejs";

interface Toast {
    show: boolean;
    init(): void;
    remove(): void;
}

export default function toast(): AlpineComponent<Toast> {
    return {
        show: false,
        init() {
            this.$nextTick(() => {
                this.show = true;
            });
            setTimeout(() => this.remove(), 2000);
        },
        remove() {
            this.show = false;
            setTimeout(() => this.$root.remove(), 500);
        },
    };
}
