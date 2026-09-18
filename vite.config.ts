import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
    plugins: [tailwindcss()],
    base: "/static/dist/",
    build: {
        outDir: "./static/dist/",
        assetsDir: "",
        emptyOutDir: true,
        manifest: "manifest.json",
        rollupOptions: {
            input: {
                view: "./internal/web/view/view.ts",
            },
        },
    },
});
