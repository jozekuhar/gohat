import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
    plugins: [tailwindcss()],
    build: {
        outDir: "./static/dist/",
        assetsDir: "",
        emptyOutDir: true,
        manifest: "manifest.json",
        rollupOptions: {
            input: {
                view: "./internal/view/view.ts",
            },
        },
    },
});
