import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
    plugins: [tailwindcss()],
    build: {
        outDir: "./static",
        emptyOutDir: true,
        assetsDir: "dist",
        manifest: "dist/manifest.json",
        rollupOptions: {
            input: {
                view: "./internal/view/view.ts",
            },
        },
    },
});
