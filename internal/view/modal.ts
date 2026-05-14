interface ModalEventDetail {
    title: string;
}

interface ModalComponent {
    isOpen: boolean;
    title: string;
    open(e: CustomEvent<ModalEventDetail>): void;
    close(): void;
}

export default function modalComponent(): ModalComponent {
    return {
        isOpen: false,
        title: "",
        open(e) {
            this.title = e.detail.title;
            this.isOpen = true;
        },
        close() {
            console.log("close");
            this.isOpen = false;
        },
    };
}
