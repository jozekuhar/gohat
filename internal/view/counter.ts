interface CounterComponent {
    count: number;
    inc(): void;
    dec(): void;
}

export default function CounterComponent(): CounterComponent {
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
