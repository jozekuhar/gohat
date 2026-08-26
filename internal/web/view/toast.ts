import { AlpineComponent } from "alpinejs";

interface Toast {
    show: boolean;
    init(): void;
    destroy(): void;
}

export default function toast(): AlpineComponent<Toast> {
    return {
        show: false,
        init() {
            this.$nextTick(() => {
                this.show = true;
            });
            setTimeout(() => this.destroy(), 10000);
        },
        destroy() {
            this.show = false;
            setTimeout(() => this.$root.remove(), 500);
        },
    };
}
