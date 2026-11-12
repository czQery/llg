import browserslist from "browserslist"
import { browserslistToTargets } from "lightningcss"
import { resolve } from "path"
import { defineConfig, loadEnv } from "vite"
import {ViteMinifyPlugin} from "vite-plugin-minify"
import {svelte} from "@sveltejs/vite-plugin-svelte"

const root = resolve(import.meta.dirname, "src")

export default ({mode}) => {
    Object.assign(process.env, loadEnv(mode, process.cwd()));
    return defineConfig({
		root,
        css: { transformer: "lightningcss", lightningcss: { targets: browserslistToTargets(browserslist("defaults")) } },
        plugins: [
			svelte(),
			ViteMinifyPlugin(),
        ],
		build: {
			minify: "terser",
			cssMinify: "lightningcss",
			target: "esnext",
			outDir: "../dist",
			emptyOutDir: true,
			sourcemap: false,
			terserOptions: {
				mangle: { toplevel: true },
				compress: { passes: 2 },
				format: { comments: false },
			},
        },
        envDir: "../"
    })
}
