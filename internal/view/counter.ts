interface CounterComponent {
    count: number;
    inc(): void;
    dec(): void;
}

export default function counterComponent(): CounterComponent {
    return {
        count: 0,
        inc() {
            this.count++;
        },
        dec() {
            this.count--;
        },
    };
}
