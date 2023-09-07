import {resolve} from "path"
import {defineConfig, loadEnv} from "vite"
import {svelte} from "@sveltejs/vite-plugin-svelte"

const root = resolve(__dirname, "src")

export default ({mode}) => {
    Object.assign(process.env, loadEnv(mode, process.cwd()));
    return defineConfig({
        root,
        plugins: [
            svelte()
        ],
        build: {
            outDir: "../dist",
            emptyOutDir: true,
        },
        envDir: "../"
    })
}