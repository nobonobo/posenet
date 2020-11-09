build:
	yarn install
	esbuild index.js --format=esm --global-name=posenet --bundle --platform=node --outfile=dist/posenet.js
