import { AlpineComponent } from "alpinejs";

interface Button {
    loading: boolean;
    showLoading: boolean;
    loadingTimer: number | undefined;
    setLoading(): void;
    clearLoading(): void;
}

export function button(): AlpineComponent<Button> {
    return {
        loading: false,
        showLoading: false,
        loadingTimer: undefined,

        // setLoading sets loading to true after 150ms to prevent instant glitches of loading state.
        setLoading() {
            this.loading = true;
            this.loadingTimer = setTimeout(() => {
                this.showLoading = true;
            }, 150);
        },

        clearLoading() {
            clearTimeout(this.loadingTimer);
            this.loading = false;
            this.showLoading = false;
        },
    };
}
