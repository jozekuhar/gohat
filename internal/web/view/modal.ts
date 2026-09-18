import { AlpineComponent } from "alpinejs";

interface Modal {
    show: boolean;
    init(): void;
    destroy(): void;
}

export default function modal(): AlpineComponent<Modal> {
    return {
        show: false,
        init() {
            this.$nextTick(() => {
                this.show = true;
            });
        },
        destroy() {
            this.show = false;
            setTimeout(() => {
                this.$root.remove();
            }, 500);
        },
    };
}
