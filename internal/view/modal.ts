import { AlpineComponent } from "alpinejs";

interface ModalEventDetail {
    title: string;
}

interface Modal {
    show: boolean;
    title: string;
    open(e: CustomEvent<ModalEventDetail>): void;
    close(): void;
}

export default function modal(): AlpineComponent<Modal> {
    return {
        show: false,
        title: "",
        open(e) {
            this.title = e.detail.title;
            this.show = true;
        },
        close() {
            this.show = false;
        },
    };
}
